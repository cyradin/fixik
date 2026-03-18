package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository struct {
	db *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) Create(ctx context.Context, i *Incident) error {
	const query = `
		INSERT INTO incidents (
			title,
			description,
			priority_id,
			status_id,
			team_id,
			user_id,
			author_id
		)
		VALUES (
			@title,
			@description,
			@priority_id,
			@status_id,
			@team_id,
			@user_id,
			@author_id
		)
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"title":       i.Title,
		"description": i.Description,
		"priority_id": i.PriorityID,
		"status_id":   i.StatusID,
		"team_id":     i.TeamID,
		"user_id":     i.UserID,
		"author_id":   i.AuthorID,
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

func (r *IncidentRepository) GetByID(ctx context.Context, id int64) (Incident, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			priority_id,
			status_id,
			team_id,
			user_id,
			author_id,
			created_at,
			updated_at,
			deleted_at
		FROM incidents
		WHERE id = @id
		  AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{"id": id}

	var i Incident

	err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.PriorityID,
		&i.StatusID,
		&i.TeamID,
		&i.UserID,
		&i.AuthorID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Incident{}, ErrNotFound
		}

		return Incident{}, fmt.Errorf("db query: %w", err)
	}

	return i, nil
}

func (r *IncidentRepository) Update(ctx context.Context, i *Incident) error {
	const query = `
		UPDATE incidents
		SET
			title = @title,
			description = @description,
			priority_id = @priority_id,
			status_id = @status_id,
			team_id = @team_id,
			user_id = @user_id,
			author_id = @author_id,
			updated_at = now()
		WHERE id = @id
		  AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := pgx.NamedArgs{
		"id":          i.ID,
		"title":       i.Title,
		"description": i.Description,
		"priority_id": i.PriorityID,
		"status_id":   i.StatusID,
		"team_id":     i.TeamID,
		"user_id":     i.UserID,
		"author_id":   i.AuthorID,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&i.UpdatedAt); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *IncidentRepository) Delete(ctx context.Context, id int64) error {
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

func (r *IncidentRepository) List(ctx context.Context, limit, offset int) ([]Incident, error) {
	const query = `
		SELECT
			id,
			title,
			description,
			priority_id,
			status_id,
			team_id,
			user_id,
			author_id,
			created_at,
			updated_at,
			deleted_at
		FROM incidents
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT @limit
		OFFSET @offset
	`

	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	var result []Incident

	for rows.Next() {
		var i Incident
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.PriorityID,
			&i.StatusID,
			&i.TeamID,
			&i.UserID,
			&i.AuthorID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, i)
	}

	return result, nil
}
