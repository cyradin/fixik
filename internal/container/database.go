package container

import (
	"context"
	"fmt"
	"time"

	"github.com/cyradin/fixik/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (c *Container) IncidentRepository() *db.IncidentRepository {
	if c.incidentRepo == nil {
		c.incidentRepo = db.NewIncidentRepository(c.PgPool())
	}

	return c.incidentRepo
}

func (c *Container) StatusRepository() *db.StatusRepository {
	if c.statusRepo == nil {
		c.statusRepo = db.NewStatusRepository(c.PgPool())
	}

	return c.statusRepo
}

func (c *Container) ImpactRepository() *db.ImpactRepository {
	if c.impactRepo == nil {
		c.impactRepo = db.NewImpactRepository(c.PgPool())
	}

	return c.impactRepo
}

func (c *Container) PriorityRepository() *db.PriorityRepository {
	if c.priorityRepo == nil {
		c.priorityRepo = db.NewPriorityRepository(c.PgPool())
	}

	return c.priorityRepo
}

func (c *Container) TeamRepository() *db.TeamRepository {
	if c.teamRepo == nil {
		c.teamRepo = db.NewTeamRepository(c.PgPool())
	}

	return c.teamRepo
}

func (c *Container) RoleRepository() *db.RoleRepository {
	if c.roleRepo == nil {
		c.roleRepo = db.NewRoleRepository(c.PgPool())
	}

	return c.roleRepo
}
