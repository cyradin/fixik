package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPServer      HTTPServerConfig
	HTTPDebugServer HTTPDebugServerConfig
	Log             LogConfig
	Postgres        PostgresConfig
}

type HTTPServerConfig struct {
	Addr              string        `envconfig:"FIXIK_HTTP_SERVER_ADDR" required:"true"`
	ReadHeaderTimeout time.Duration `envconfig:"FIXIK_HTTP_SERVER_READ_HEADER_TIMEOUT" required:"true"`
	ReadTimeout       time.Duration `envconfig:"FIXIK_HTTP_SERVER_READ_TIMEOUT" required:"true"`
	WriteTimeout      time.Duration `envconfig:"FIXIK_HTTP_SERVER_WRITE_TIMEOUT" required:"true"`
	IdleTimeout       time.Duration `envconfig:"FIXIK_HTTP_SERVER_IDLE_TIMEOUT" required:"true"`
	AllowedOrigins    []string      `envconfig:"FIXIK_HTTP_SERVER_ALLOWED_ORIGINS" required:"true"`
}

type HTTPDebugServerConfig struct {
	Addr              string        `envconfig:"FIXIK_HTTP_DEBUG_SERVER_ADDR" required:"true"`
	ReadHeaderTimeout time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_READ_HEADER_TIMEOUT" required:"true"`
	ReadTimeout       time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_READ_TIMEOUT" required:"true"`
	WriteTimeout      time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_WRITE_TIMEOUT" required:"true"`
	IdleTimeout       time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_IDLE_TIMEOUT" required:"true"`
	// Куда ходит веб-интерфейс сваггера в Try it out
	SwaggerAddr string `envconfig:"FIXIK_HTTP_DEBUG_SWAGGER_ADDR" required:"true"`
}

type LogConfig struct {
	Level string `envconfig:"FIXIK_LOG_LEVEL" required:"true"`
}

type PostgresConfig struct {
	URL string `envconfig:"FIXIK_POSTGRES_URL" required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
