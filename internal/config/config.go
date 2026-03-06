package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPServer HTTPServerConfig
}

type HTTPServerConfig struct {
	Addr              string        `envconfig:"FIXIK_HTTP_SERVER_ADDR" required:"true"`
	ReadHeaderTimeout time.Duration `envconfig:"FIXIK_HTTP_SERVER_READ_HEADER_TIMEOUT" required:"true"`
	ReadTimeout       time.Duration `envconfig:"FIXIK_HTTP_SERVER_READ_TIMEOUT" required:"true"`
	WriteTimeout      time.Duration `envconfig:"FIXIK_HTTP_SERVER_WRITE_TIMEOUT" required:"true"`
	IdleTimeout       time.Duration `envconfig:"FIXIK_HTTP_SERVER_IDLE_TIMEOUT" required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
