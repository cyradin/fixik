package container

import (
	"context"
	"fmt"
	"time"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/pkg/transaction"
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

		if err := pool.Ping(ctx); err != nil {
			panic(err)
		}

		c.pgPool = pool
	}

	return c.pgPool
}

func (c *Container) TxExecutor() *transaction.Executor {
	if c.txExecutor == nil {
		c.txExecutor = transaction.NewExecutor(c.PgPool())
	}

	return c.txExecutor
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

func (c *Container) PriorityRepository() *db.DictRepository {
	if c.priorityRepo == nil {
		c.priorityRepo = db.NewPriorityRepository(c.PgPool())
	}

	return c.priorityRepo
}

func (c *Container) TeamRepository() *db.DictRepository {
	if c.teamRepo == nil {
		c.teamRepo = db.NewTeamRepository(c.PgPool())
	}

	return c.teamRepo
}

func (c *Container) UserRepository() *db.UserRepository {
	if c.userRepo == nil {
		c.userRepo = db.NewUserRepository(c.PgPool())
	}

	return c.userRepo
}

func (c *Container) CommentRepository() *db.CommentRepository {
	if c.commentRepo == nil {
		c.commentRepo = db.NewCommentRepository(c.PgPool())
	}

	return c.commentRepo
}
