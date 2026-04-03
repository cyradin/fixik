package db

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}
func (r *CommentRepository) Create(ctx context.Context, c *Comment) error {
	const query = `
		INSERT INTO comments (author_id, incident_id, text, created_at, updated_at)
		VALUES (@author_id, @incident_id, @text, now(), now())
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"author_id":   c.AuthorID,
		"incident_id": c.IncidentID,
		"text":        c.Text,
	}

	tx := transaction.FromContext(ctx, r.db)

	if err := tx.QueryRow(ctx, query, args).Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return fmt.Errorf("insert comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) Update(ctx context.Context, c *Comment) error {
	const query = `
		UPDATE comments
		SET text = @text, updated_at = now()
		WHERE id = @id AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{
		"id":   c.ID,
		"text": c.Text,
	}

	tx := transaction.FromContext(ctx, r.db)

	tag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("update comment: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE comments
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tx := transaction.FromContext(ctx, r.db)

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete comment: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
func (r *CommentRepository) ListByIncident(ctx context.Context, incidentID int64, limit, offset int) (CommentListResult, error) {
	const listQuery = `
		SELECT
			id,
			incident_id,
			author_id,
			text,
			created_at,
			updated_at
		FROM comments
		WHERE incident_id = $1 AND deleted_at IS NULL
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`

	const countQuery = `
		SELECT COUNT(*)
		FROM comments
		WHERE incident_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(ctx, listQuery, incidentID, limit, offset)
	if err != nil {
		return CommentListResult{}, fmt.Errorf("db list query: %w", err)
	}
	defer rows.Close()

	var items []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(
			&c.ID,
			&c.IncidentID,
			&c.AuthorID,
			&c.Text,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return CommentListResult{}, fmt.Errorf("scan: %w", err)
		}
		items = append(items, c)
	}

	var total int
	if err := r.db.QueryRow(ctx, countQuery, incidentID).Scan(&total); err != nil {
		return CommentListResult{}, fmt.Errorf("db count query: %w", err)
	}

	return CommentListResult{
		Items: items,
		Total: total,
	}, nil
}
