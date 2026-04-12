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
	Auth            AuthConfig
}

//nolint:gosec
type AuthConfig struct {
	Secret          string        `envconfig:"FIXIK_AUTH_SECRET" required:"true" default:"secret"`
	AccessTokenTTL  time.Duration `envconfig:"FIXIK_AUTH_ACCESS_TOKEN_TTL" required:"true" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"FIXIK_AUTH_REFRESH_TOKEN_TTL" required:"true" default:"672h"` // default - 28 days
	SecureCookies   bool          `envconfig:"FIXIK_AUTH_SECURE_COOKIES" required:"true" default:"false"`
}

type HTTPServerConfig struct {
	Addr              string        `envconfig:"FIXIK_HTTP_SERVER_ADDR" required:"true" default:":8080"`
	ReadHeaderTimeout time.Duration `envconfig:"FIXIK_HTTP_SERVER_READ_HEADER_TIMEOUT" required:"true" default:"5s"`
	ReadTimeout       time.Duration `envconfig:"FIXIK_HTTP_SERVER_READ_TIMEOUT" required:"true" default:"10s"`
	WriteTimeout      time.Duration `envconfig:"FIXIK_HTTP_SERVER_WRITE_TIMEOUT" required:"true" default:"10s"`
	IdleTimeout       time.Duration `envconfig:"FIXIK_HTTP_SERVER_IDLE_TIMEOUT" required:"true" default:"60s"`
	AllowedOrigins    []string      `envconfig:"FIXIK_HTTP_SERVER_ALLOWED_ORIGINS" required:"true" default:"http://localhost:8081"`
}

type HTTPDebugServerConfig struct {
	Addr              string        `envconfig:"FIXIK_HTTP_DEBUG_SERVER_ADDR" required:"true" default:":8081"`
	ReadHeaderTimeout time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_READ_HEADER_TIMEOUT" required:"true" default:"5s"`
	ReadTimeout       time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_READ_TIMEOUT" required:"true" default:"10s"`
	WriteTimeout      time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_WRITE_TIMEOUT" required:"true" default:"10s"`
	IdleTimeout       time.Duration `envconfig:"FIXIK_HTTP_DEBUG_SERVER_IDLE_TIMEOUT" required:"true" default:"60s"`
	// Куда ходит веб-интерфейс сваггера в Try it out
	SwaggerAddr string `envconfig:"FIXIK_HTTP_DEBUG_SWAGGER_ADDR" required:"true" default:"localhost:8080"`
}

type LogConfig struct {
	Level string `envconfig:"FIXIK_LOG_LEVEL" required:"true" default:"info"`
}

type PostgresConfig struct {
	URL string `envconfig:"FIXIK_POSTGRES_URL" required:"true" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
