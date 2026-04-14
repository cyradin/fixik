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
			i.id,
			i.title,
			i.description,
			i.priority_id,
			i.status_id,
			i.team_id,
			i.user_id,
			i.author_id,
			i.created_at,
			i.updated_at,
			i.deleted_at,
			COALESCE(c.comments_count, 0) AS comments_count
		FROM incidents i
		LEFT JOIN (
			SELECT incident_id, COUNT(*) AS comments_count
			FROM comments
			WHERE deleted_at IS NULL
			GROUP BY incident_id
		) c ON c.incident_id = i.id
		WHERE i.id = @id
		  AND i.deleted_at IS NULL
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
		&i.CommentsCount,
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
func (r *IncidentRepository) List(ctx context.Context, limit, offset int) (IncidentListResult, error) {
	const listQuery = `
		SELECT
			i.id,
			i.title,
			i.description,
			i.priority_id,
			i.status_id,
			i.team_id,
			i.user_id,
			i.author_id,
			i.created_at,
			i.updated_at,
			i.deleted_at,
			COALESCE(c.comments_count, 0) AS comments_count
		FROM incidents i
		LEFT JOIN (
			SELECT incident_id, COUNT(*) AS comments_count
			FROM comments
			WHERE deleted_at IS NULL
			GROUP BY incident_id
		) c ON c.incident_id = i.id
		WHERE i.deleted_at IS NULL
		ORDER BY i.id DESC
		LIMIT @limit
		OFFSET @offset
	`

	const countQuery = `
		SELECT COUNT(*)
		FROM incidents
		WHERE deleted_at IS NULL
	`

	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, listQuery, args)
	if err != nil {
		return IncidentListResult{}, fmt.Errorf("db list query: %w", err)
	}
	defer rows.Close()

	var items []Incident

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
			&i.CommentsCount,
		); err != nil {
			return IncidentListResult{}, fmt.Errorf("scan: %w", err)
		}

		items = append(items, i)
	}

	var total int
	if err := transaction.FromContext(ctx, r.db).
		QueryRow(ctx, countQuery).
		Scan(&total); err != nil {
		return IncidentListResult{}, fmt.Errorf("db count query: %w", err)
	}

	return IncidentListResult{
		Items: items,
		Total: total,
	}, nil
}

func (r *IncidentRepository) CountByStatus(ctx context.Context, statusID int64) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM incidents
		WHERE status_id = @id
		  AND deleted_at IS NULL
	`

	return r.count(ctx, query, statusID, "status")
}

func (r *IncidentRepository) CountByPriority(ctx context.Context, priorityID int64) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM incidents
		WHERE priority_id = @id
		  AND deleted_at IS NULL
	`

	return r.count(ctx, query, priorityID, "priority")
}

func (r *IncidentRepository) CountByTeam(ctx context.Context, teamID int64) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM incidents
		WHERE team_id = @id
		  AND deleted_at IS NULL
	`

	return r.count(ctx, query, teamID, "team")
}

func (r *IncidentRepository) CountByUser(ctx context.Context, userID int64) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM incidents
		WHERE (user_id = @id OR author_id = @id)
		  AND deleted_at IS NULL
	`

	return r.count(ctx, query, userID, "user")
}

func (r *IncidentRepository) count(ctx context.Context, query string, id int64, entity string) (int, error) {
	args := pgx.NamedArgs{
		"id": id,
	}

	var count int

	if err := transaction.FromContext(ctx, r.db).
		QueryRow(ctx, query, args).
		Scan(&count); err != nil {
		return 0, fmt.Errorf("db count by %s: %w", entity, err)
	}

	return count, nil
}
