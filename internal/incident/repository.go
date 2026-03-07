package incident

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, i *Incident) error {
	const query = `
		INSERT INTO incidents (
			title,
			description,
			impact,
			urgency,
			priority,
			status
		)
		VALUES (
			@title,
			@description,
			@impact,
			@urgency,
			@priority,
			@status
		)
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"title":       i.Title,
		"description": i.Description,
		"impact":      i.Impact,
		"urgency":     i.Urgency,
		"priority":    i.Priority,
		"status":      i.Status,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Incident, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			impact,
			urgency,
			priority,
			status,
			created_at,
			updated_at
		FROM incidents
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var i Incident

	err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Impact,
		&i.Urgency,
		&i.Priority,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	return &i, nil
}

func (r *Repository) Update(ctx context.Context, i *Incident) error {
	const query = `
		UPDATE incidents
		SET
			title = @title,
			description = @description,
			impact = @impact,
			urgency = @urgency,
			priority = @priority,
			status = @status,
			updated_at = now()
		WHERE id = @id
		RETURNING updated_at
	`

	args := pgx.NamedArgs{
		"id":          i.ID,
		"title":       i.Title,
		"description": i.Description,
		"impact":      i.Impact,
		"urgency":     i.Urgency,
		"priority":    i.Priority,
		"status":      i.Status,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&i.UpdatedAt); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE incidents
		SET
			deleted_at = now(),
			updated_at = now()
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
