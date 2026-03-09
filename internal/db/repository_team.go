package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, t *Team) error {
	const query = `
		INSERT INTO teams (
			code,
			name
		)
		VALUES (
			@code,
			@name
		)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"code": t.Code,
		"name": t.Name,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&t.ID); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id int64) (Team, error) {
	const query = `
		SELECT id, code, name
		FROM teams
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var t Team

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&t.ID,
		&t.Code,
		&t.Name,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Team{}, ErrNotFound
		}

		return Team{}, fmt.Errorf("db query: %w", err)
	}

	return t, nil
}

func (r *TeamRepository) List(ctx context.Context) ([]Team, error) {
	const query = `
		SELECT id, code, name
		FROM teams
		ORDER BY id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer rows.Close()

	var teams []Team

	for rows.Next() {
		var t Team

		if err := rows.Scan(&t.ID, &t.Code, &t.Name); err != nil {
			return nil, fmt.Errorf("scan value: %w", err)
		}

		teams = append(teams, t)
	}

	return teams, rows.Err()
}

func (r *TeamRepository) Update(ctx context.Context, t *Team) error {
	const query = `
		UPDATE teams
		SET code = @code, name = @name
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":   t.ID,
		"code": t.Code,
		"name": t.Name,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *TeamRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM teams
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
