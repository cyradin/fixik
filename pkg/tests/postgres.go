package tests

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresSuite struct {
	suite.Suite

	migrationsDir string
	postgres      *postgres.PostgresContainer
	postgresPool  *pgxpool.Pool
}

func (s *PostgresSuite) SetupSuite() {
	_, filename, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(filename)
	s.migrationsDir = filepath.Join(pkgDir, "../../migrations")
}

func (s *PostgresSuite) Postgres() *pgxpool.Pool {
	if s.postgres == nil {
		postgresContainer, pool, err := StartPostgresContainer(context.Background(), "file://"+s.migrationsDir)
		s.Require().NoError(err)

		s.postgres = postgresContainer
		s.postgresPool = pool
	}

	return s.postgresPool
}

func (s *PostgresSuite) TearDownSuite() {
	if s.postgres != nil {
		s.postgresPool.Close()

		err := s.postgres.Terminate(context.Background())
		require.NoError(s.T(), err)
	}
}

func StartPostgresContainer(ctx context.Context, migrationsPath string) (*postgres.PostgresContainer, *pgxpool.Pool, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:18.1-alpine3.23",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).                  // nolint:mnd
				WithStartupTimeout(30*time.Second), // nolint:mnd
		),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("get connection string: %w", err)
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, err
	}

	m, err := migrate.New(
		migrationsPath,
		connStr,
	)
	if err != nil {
		pool.Close()

		_ = pgContainer.Terminate(ctx)

		return nil, nil, fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		pool.Close()

		_ = pgContainer.Terminate(ctx)

		return nil, nil, fmt.Errorf("apply migrations: %w", err)
	}

	return pgContainer, pool, nil
}
