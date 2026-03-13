package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	docs "github.com/cyradin/fixik/docs"
	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/web"
	"github.com/cyradin/fixik/pkg/logger"
)

var GitCommit string = "dev"

// @title Fixik API
// @version 1.0
// @description Incident management system API
// @BasePath /api
// @host localhost:8080
// @schemes http
func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	container := container.New(GitCommit, cfg)

	if err := run(cfg, container); err != nil {
		container.Logger().Error("application error", logger.Error(err))
	}
}

func run(cfg *config.Config, container *container.Container) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		server := initHTTPServer(ctx, container, cfg)

		container.Logger().Info("started http server", logger.Address(cfg.HTTPServer.Addr))

		if err := server.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	go func() {
		server := initHTTPDebugServer(ctx, container, cfg)

		container.Logger().Info("started http debug server", logger.Address(cfg.HTTPDebugServer.Addr))

		if err := server.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("debug server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func initHTTPServer(ctx context.Context, container *container.Container, cfg *config.Config) *http.Server {
	r := web.NewRouter(container, cfg.HTTPServer.AllowedOrigins)

	return &http.Server{
		Addr:              cfg.HTTPServer.Addr,
		ReadTimeout:       cfg.HTTPServer.ReadTimeout,
		ReadHeaderTimeout: cfg.HTTPServer.ReadHeaderTimeout,
		WriteTimeout:      cfg.HTTPServer.WriteTimeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
		Handler:           r,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}

func initHTTPDebugServer(ctx context.Context, container *container.Container, cfg *config.Config) *http.Server {
	docs.SwaggerInfo.Host = cfg.HTTPDebugServer.SwaggerAddr
	r := web.NewDebugRouter(container)

	return &http.Server{
		Addr:              cfg.HTTPDebugServer.Addr,
		ReadTimeout:       cfg.HTTPDebugServer.ReadTimeout,
		ReadHeaderTimeout: cfg.HTTPDebugServer.ReadHeaderTimeout,
		WriteTimeout:      cfg.HTTPDebugServer.WriteTimeout,
		IdleTimeout:       cfg.HTTPDebugServer.IdleTimeout,
		Handler:           r,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}
