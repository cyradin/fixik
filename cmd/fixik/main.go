package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		server := initHTTPServer(ctx, cfg)

		if err := server.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func initHTTPServer(ctx context.Context, cfg *config.Config) *http.Server {
	r := router.New(cfg)

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
