package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (r *IncidentRepository) List(ctx context.Context, filter IncidentFilter, limit, offset int) (IncidentListResult, error) {
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
		LEFT JOIN statuses s ON i.status_id = s.id
		WHERE %s
		ORDER BY i.updated_at DESC
		LIMIT @limit
		OFFSET @offset
	`

	const countQuery = `
		SELECT COUNT(*)
		FROM incidents i
		LEFT JOIN statuses s ON i.status_id = s.id
		WHERE %s
	`

	conditions := []string{"i.deleted_at IS NULL"}
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	conditions = r.addFilter(conditions, args, "i.author_id", filter.AuthorIDs, "author_ids")
	conditions = r.addFilter(conditions, args, "i.user_id", filter.UserIDs, "user_ids")
	conditions = r.addFilter(conditions, args, "i.team_id", filter.TeamIDs, "team_ids")
	conditions = r.addFilter(conditions, args, "i.priority_id", filter.PriorityIDs, "priority_ids")
	conditions = r.addFilter(conditions, args, "i.status_id", filter.StatusIDs, "status_ids")

	if filter.ActiveOnly {
		conditions = append(conditions, "(s.is_final IS NOT TRUE OR i.updated_at > NOW() - INTERVAL '2 weeks')")
	}

	where := strings.Join(conditions, " AND ")

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, fmt.Sprintf(listQuery, where), args)
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
		QueryRow(ctx, fmt.Sprintf(countQuery, where), args).
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

func (r *IncidentRepository) addFilter(
	conditions []string,
	args pgx.NamedArgs,
	field string, values []int64, argName string) []string {
	if len(values) == 0 {
		return conditions
	}

	// special case: [0] => IS NULL
	if len(values) == 1 && values[0] == 0 {
		conditions = append(conditions, fmt.Sprintf("%s IS NULL", field))

		return conditions
	}

	conditions = append(conditions, fmt.Sprintf("%s = ANY(@%s)", field, argName))
	args[argName] = values

	return conditions
}
