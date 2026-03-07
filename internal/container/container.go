package container

import (
	"log/slog"

	"github.com/cyradin/fixik/internal/config"
	"github.com/cyradin/fixik/pkg/logger"
)

type Container struct {
	version string
	cfg     *config.Config
	logger  *slog.Logger
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
