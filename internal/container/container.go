package container

import (
	"log/slog"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/internal/user"
	"github.com/cyradin/fixik/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	version string
	cfg     *config.Config
	logger  *slog.Logger
	pgPool  *pgxpool.Pool

	statusRepo   *db.DictRepository
	impactRepo   *db.DictRepository
	priorityRepo *db.DictRepository
	teamRepo     *db.DictRepository
	incidentRepo *db.IncidentRepository
	userRepo     *db.UserRepository

	statusManager   *dict.EntityManager
	priorityManager *dict.EntityManager
	impactManager   *dict.EntityManager
	teamManager     *dict.EntityManager

	incidentManager *incident.IncidentManager
	userManager     *user.UserManager
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
