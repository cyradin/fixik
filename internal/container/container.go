package container

import "github.com/cyradin/fixik/internal/config"

type Container struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Container {
	return &Container{
		cfg: cfg,
	}
}
