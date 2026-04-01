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

func (r *StatusRepository) Create(ctx context.Context, e *Status) error {
	const query = `
		INSERT INTO statuses (code, name, description, sort, is_final, created_at, updated_at)
		VALUES (@code, @name, @description, @sort, @is_final, now(), now())
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"code":        e.Code,
		"name":        e.Name,
		"description": e.Description,
		"sort":        e.Sort,
		"is_final":    e.IsFinal,
	}

	err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *StatusRepository) GetByID(ctx context.Context, id int64) (Status, error) {
	const query = `
		SELECT id, code, name, description, sort, is_final, created_at, updated_at, deleted_at
		FROM statuses
		WHERE id = @id AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{"id": id}
	var e Status
	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&e.ID, &e.Code, &e.Name, &e.Description, &e.Sort, &e.IsFinal, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Status{}, ErrNotFound
		}
		return Status{}, fmt.Errorf("db query: %w", err)
	}

	return e, nil
}

func (r *StatusRepository) List(ctx context.Context) ([]Status, error) {
	const query = `
		SELECT id, code, name, description, sort, is_final, created_at, updated_at, deleted_at
		FROM statuses
		WHERE deleted_at IS NULL
		ORDER BY sort, id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	var list []Status
	for rows.Next() {
		var e Status
		if err := rows.Scan(
			&e.ID, &e.Code, &e.Name, &e.Description, &e.Sort, &e.IsFinal,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		list = append(list, e)
	}

	return list, rows.Err()
}

func (r *StatusRepository) Update(ctx context.Context, e *Status) error {
	const query = `
		UPDATE statuses
		SET code = @code, name = @name, description = @description, sort = @sort, is_final = @is_final, updated_at = now()
		WHERE id = @id AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{
		"id":          e.ID,
		"code":        e.Code,
		"name":        e.Name,
		"description": e.Description,
		"sort":        e.Sort,
		"is_final":    e.IsFinal,
	}

	tag, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *StatusRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE statuses
		SET deleted_at = now()
		WHERE id = @id AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{"id": id}
	_, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}
