package incident

import (
	"context"
	"errors"
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
			impact_id,
			priority_id,
			status_id
		)
		VALUES (
			@title,
			@description,
			@impact_id,
			@priority_id,
			@status_id
		)
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"title":       i.Title,
		"description": i.Description,
		"impact_id":   i.ImpactID,
		"priority_id": i.PriorityID,
		"status_id":   i.StatusID,
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

func (r *Repository) GetByID(ctx context.Context, id int64) (Incident, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			impact_id,
			priority_id,
			status_id,
			created_at,
			updated_at,
			deleted_at
		FROM incidents
		WHERE id = @id
		  AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var incident Incident

	err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Description,
		&incident.ImpactID,
		&incident.PriorityID,
		&incident.StatusID,
		&incident.CreatedAt,
		&incident.UpdatedAt,
		&incident.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Incident{}, nil
		}

		return Incident{}, fmt.Errorf("db query: %w", err)
	}

	return incident, nil
}

func (r *Repository) Update(ctx context.Context, i *Incident) error {
	const query = `
		UPDATE incidents
		SET
			title = @title,
			description = @description,
			impact_id = @impact_id,
			priority_id = @priority_id,
			status_id = @status_id,
			updated_at = now()
		WHERE id = @id
		  AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := pgx.NamedArgs{
		"id":          i.ID,
		"title":       i.Title,
		"description": i.Description,
		"impact_id":   i.ImpactID,
		"priority_id": i.PriorityID,
		"status_id":   i.StatusID,
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
		  AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}
