package container

import (
	"log/slog"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/internal/priority"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/team"
	"github.com/cyradin/fixik/internal/user"
	"github.com/cyradin/fixik/pkg/logger"
	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	version    string
	cfg        *config.Config
	logger     *slog.Logger
	pgPool     *pgxpool.Pool
	txExecutor *transaction.Executor

	statusRepo   *db.StatusRepository
	priorityRepo *db.DictRepository
	teamRepo     *db.DictRepository
	incidentRepo *db.IncidentRepository
	userRepo     *db.UserRepository
	commentRepo  *db.CommentRepository

	statusManager *status.StatusManager

	priorityManager *priority.PriorityManager
	teamManager     *team.TeamManager

	incidentManager *incident.IncidentManager
	commentManager  *incident.CommentManager
	userManager     *user.UserManager
	jwtManager      *user.JWTManager
	authService     *user.AuthService
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

func (c *Container) Cfg() config.Config {
	return *c.cfg
}
