package container

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	version string
	cfg     *config.Config
	logger  *slog.Logger
	pgPool  *pgxpool.Pool

	incidentRepo *incident.Repository
}

func New(
	version string,
	cfg *config.Config,
) *Container {
	return &Container{
		version: version,
		cfg:     cfg,
	}
}

func (c *Container) Logger() *slog.Logger {
	if c.logger == nil {
		c.logger = logger.New(c.version, c.cfg.Log.Level)
	}

	return c.logger
}

func (c *Container) PgPool() *pgxpool.Pool {
	const dbInitTimeout = 3 * time.Second

	if c.pgPool == nil {
		ctx, cancel := context.WithTimeout(context.Background(), dbInitTimeout)
		defer cancel()

		pool, err := pgxpool.New(ctx, c.cfg.Postgres.URL)
		if err != nil {
			panic(fmt.Errorf("init database: %w", err))
		}

		c.pgPool = pool
	}

	return c.pgPool
}

func (c *Container) IncidentRepository() *incident.Repository {
	if c.incidentRepo == nil {
		c.incidentRepo = incident.NewRepository(c.PgPool())
	}

	return c.incidentRepo
}
