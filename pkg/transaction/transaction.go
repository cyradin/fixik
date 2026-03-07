package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgres interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ctxKey struct{}

var txKey ctxKey = struct{}{}

func FromContext(ctx context.Context, db postgres) postgres {
	if tx, ok := ctx.Value(txKey).(postgres); ok {
		return tx
	}

	return db
}

type Executor struct {
	db postgres
}

func NewExecutor(db postgres) *Executor {
	return &Executor{db: db}
}

func (e *Executor) Exec(ctx context.Context, callback func(ctx context.Context) error) (txError error) {
	if tx := ctx.Value(txKey); tx != nil {
		return callback(ctx)
	}

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			txError = err
		}
	}()

	ctx = context.WithValue(ctx, txKey, tx)

	if err := callback(ctx); err != nil {
		return fmt.Errorf("exec tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
