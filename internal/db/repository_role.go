package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *Role) error {
	const query = `
		INSERT INTO roles (
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
		"code": role.Code,
		"name": role.Name,
	}

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(&role.ID); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id int64) (Role, error) {
	const query = `
		SELECT id, code, name
		FROM roles
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var role Role

	if err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, args).Scan(
		&role.ID,
		&role.Code,
		&role.Name,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Role{}, ErrNotFound
		}

		return Role{}, fmt.Errorf("db query: %w", err)
	}

	return role, nil
}

func (r *RoleRepository) List(ctx context.Context) ([]Role, error) {
	const query = `
		SELECT id, code, name
		FROM roles
		ORDER BY id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer rows.Close()

	var roles []Role

	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.Code, &role.Name); err != nil {
			return nil, fmt.Errorf("scan value: %w", err)
		}

		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *RoleRepository) Update(ctx context.Context, role *Role) error {
	const query = `
		UPDATE roles
		SET code = @code, name = @name
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":   role.ID,
		"code": role.Code,
		"name": role.Name,
	}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM roles
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
