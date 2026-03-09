package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusRepository struct {
	db *pgxpool.Pool
}

func NewStatusRepository(db *pgxpool.Pool) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) Create(ctx context.Context, s *Status) error {
	const query = `
		INSERT INTO incident_statuses (
			code,
			name
		)
		VALUES (
			@code,
			@name
		)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"code": s.Code,
		"name": s.Name,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&s.ID); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *StatusRepository) GetByID(ctx context.Context, id int64) (Status, error) {
	const query = `
		SELECT id, code, name
		FROM incident_statuses
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var s Status

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&s.ID,
		&s.Code,
		&s.Name,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Status{}, ErrNotFound
		}

		return Status{}, fmt.Errorf("db query: %w", err)
	}

	return s, nil
}

func (r *StatusRepository) List(ctx context.Context) ([]Status, error) {
	const query = `
		SELECT id, code, name
		FROM incident_statuses
		ORDER BY id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer rows.Close()

	var statuses []Status

	for rows.Next() {
		var s Status

		if err := rows.Scan(&s.ID, &s.Code, &s.Name); err != nil {
			return nil, fmt.Errorf("scan value: %w", err)
		}

		statuses = append(statuses, s)
	}

	return statuses, rows.Err()
}

func (r *StatusRepository) Update(ctx context.Context, s *Status) error {
	const query = `
		UPDATE incident_statuses
		SET code = @code, name = @name
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":   s.ID,
		"code": s.Code,
		"name": s.Name,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *StatusRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM incident_statuses
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}
