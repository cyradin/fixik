package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DictRepository struct {
	db        *pgxpool.Pool
	tableName string
}

func newDictRepository(db *pgxpool.Pool, tableName string) *DictRepository {
	return &DictRepository{db: db, tableName: tableName}
}

func NewImpactRepository(db *pgxpool.Pool) *DictRepository {
	return newDictRepository(db, "incident_impacts")
}

func NewPriorityRepository(db *pgxpool.Pool) *DictRepository {
	return newDictRepository(db, "incident_priorities")
}

func NewStatusRepository(db *pgxpool.Pool) *DictRepository {
	return newDictRepository(db, "incident_statuses")
}

func NewTeamRepository(db *pgxpool.Pool) *DictRepository {
	return newDictRepository(db, "teams")
}

func NewRoleRepository(db *pgxpool.Pool) *DictRepository {
	return newDictRepository(db, "roles")
}

func (r *DictRepository) Create(ctx context.Context, e *DictEntity) error {
	const queryTemplate = `
		INSERT INTO %s (code, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING id, created_at, updated_at
	`

	query := fmt.Sprintf(queryTemplate, r.tableName)

	err := transaction.FromContext(ctx, r.db).QueryRow(
		ctx,
		query,
		e.Code,
		e.Name,
		e.Description,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *DictRepository) GetByID(ctx context.Context, id int64) (DictEntity, error) {
	const queryTemplate = `
		SELECT id, code, name, description, created_at, updated_at, deleted_at
		FROM %s
		WHERE id = $1 AND deleted_at IS NULL
	`

	query := fmt.Sprintf(queryTemplate, r.tableName)

	var e DictEntity
	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, id).Scan(
		&e.ID, &e.Code, &e.Name, &e.Description, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DictEntity{}, ErrNotFound
		}

		return DictEntity{}, fmt.Errorf("db query: %w", err)
	}

	return e, nil
}

func (r *DictRepository) List(ctx context.Context) ([]DictEntity, error) {
	const queryTemplate = `
		SELECT id, code, name, description, created_at, updated_at, deleted_at
		FROM %s
		WHERE deleted_at IS NULL
		ORDER BY id
	`

	query := fmt.Sprintf(queryTemplate, r.tableName)

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	var list []DictEntity

	for rows.Next() {
		var e DictEntity
		if err := rows.Scan(&e.ID, &e.Code, &e.Name, &e.Description, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		list = append(list, e)
	}

	return list, rows.Err()
}

func (r *DictRepository) Update(ctx context.Context, e *DictEntity) error {
	const queryTemplate = `
		UPDATE %s
		SET code = $1, name = $2, description = $3, updated_at = now()
		WHERE id = $4 AND deleted_at IS NULL
	`

	query := fmt.Sprintf(queryTemplate, r.tableName)

	tag, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, e.Code, e.Name, e.Description, e.ID)
	if err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *DictRepository) Delete(ctx context.Context, id int64) error {
	const queryTemplate = `
		UPDATE %s
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`

	query := fmt.Sprintf(queryTemplate, r.tableName)

	_, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}
