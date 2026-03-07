package incident

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PriorityRepository struct {
	db *pgxpool.Pool
}

func NewPriorityRepository(db *pgxpool.Pool) *PriorityRepository {
	return &PriorityRepository{db: db}
}

func (r *PriorityRepository) Create(ctx context.Context, p *Priority) error {
	const query = `
		INSERT INTO incident_priorities (code, name)
		VALUES (@code, @name)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"code": p.Code,
		"name": p.Name,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&p.ID); err != nil {
		return fmt.Errorf("create priority: %w", err)
	}

	return nil
}

func (r *PriorityRepository) GetByID(ctx context.Context, id int64) (Priority, error) {
	const query = `
		SELECT id, code, name
		FROM incident_priorities
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	var p Priority

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&p.ID, &p.Code, &p.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Priority{}, nil
		}

		return Priority{}, fmt.Errorf("get priority: %w", err)
	}

	return p, nil
}

func (r *PriorityRepository) List(ctx context.Context) ([]Priority, error) {
	const query = `
		SELECT id, code, name
		FROM incident_priorities
		ORDER BY id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list priorities: %w", err)
	}

	defer rows.Close()

	var list []Priority

	for rows.Next() {
		var p Priority
		if err := rows.Scan(&p.ID, &p.Code, &p.Name); err != nil {
			return nil, fmt.Errorf("scan priority: %w", err)
		}

		list = append(list, p)
	}

	return list, rows.Err()
}

func (r *PriorityRepository) Update(ctx context.Context, p *Priority) error {
	const query = `
		UPDATE incident_priorities
		SET code = @code, name = @name
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":   p.ID,
		"code": p.Code,
		"name": p.Name,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("update priority: %w", err)
	}

	return nil
}

func (r *PriorityRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM incident_priorities
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("delete priority: %w", err)
	}

	return nil
}
