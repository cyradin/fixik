package container

import (
	"log/slog"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	version string
	cfg     *config.Config
	logger  *slog.Logger
	pgPool  *pgxpool.Pool

	incidentRepo *db.IncidentRepository
	statusRepo   *db.StatusRepository
	impactRepo   *db.ImpactRepository
	priorityRepo *db.PriorityRepository
	teamRepo     *db.TeamRepository
	roleRepo     *db.RoleRepository

	statusManager   *incident.StatusManager
	priorityManager *incident.PriorityManager
	impactManager   *incident.ImpactManager
	incidentManager *incident.IncidentManager
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
