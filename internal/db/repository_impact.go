package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImpactRepository struct {
	db *pgxpool.Pool
}

func NewImpactRepository(db *pgxpool.Pool) *ImpactRepository {
	return &ImpactRepository{db: db}
}

func (r *ImpactRepository) Create(ctx context.Context, im *Impact) error {
	const query = `
		INSERT INTO incident_impacts (code, name)
		VALUES (@code, @name)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"code": im.Code,
		"name": im.Name,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&im.ID); err != nil {
		return fmt.Errorf("create impact: %w", err)
	}

	return nil
}

func (r *ImpactRepository) GetByID(ctx context.Context, id int64) (Impact, error) {
	const query = `
		SELECT id, code, name
		FROM incident_impacts
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	var impact Impact
	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&impact.ID, &impact.Code, &impact.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Impact{}, ErrNotFound
		}

		return Impact{}, fmt.Errorf("get impact: %w", err)
	}

	return impact, nil
}

func (r *ImpactRepository) List(ctx context.Context) ([]Impact, error) {
	const query = `
		SELECT id, code, name
		FROM incident_impacts
		ORDER BY id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list impacts: %w", err)
	}

	defer rows.Close()

	var list []Impact

	for rows.Next() {
		var im Impact
		if err := rows.Scan(&im.ID, &im.Code, &im.Name); err != nil {
			return nil, fmt.Errorf("scan impact: %w", err)
		}

		list = append(list, im)
	}

	return list, rows.Err()
}

func (r *ImpactRepository) Update(ctx context.Context, im *Impact) error {
	const query = `
		UPDATE incident_impacts
		SET code = @code, name = @name
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":   im.ID,
		"code": im.Code,
		"name": im.Name,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("update impact: %w", err)
	}

	return nil
}

func (r *ImpactRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM incident_impacts
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("delete impact: %w", err)
	}

	return nil
}
