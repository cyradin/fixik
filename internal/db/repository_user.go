package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/pkg/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	const query = `
		INSERT INTO users (name, username, email, password, team_id, role, created_at, updated_at)
		VALUES (@name, @username, @email, @password, @team_id, @role, now(), now())
		RETURNING id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"name":     u.Name,
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
		"team_id":  u.TeamID,
		"role":     u.Role,
	}

	tx := transaction.FromContext(ctx, r.db)
	if err := tx.QueryRow(ctx, query, args).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (User, error) {
	const query = `
		SELECT
			u.id, u.name, u.username, u.email, u.password, u.team_id,
			u.role, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		WHERE u.id = $1 AND u.deleted_at IS NULL
		GROUP BY u.id
	`

	var u User

	err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.TeamID,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("db query: %w", err)
	}

	return u, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
	const query = `
		SELECT
			u.id, u.name, u.username, u.email, u.password, u.team_id,
			u.role, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		WHERE u.deleted_at IS NULL
		GROUP BY u.id
		ORDER BY u.id
		LIMIT $1 OFFSET $2
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0, limit)

	for rows.Next() {
		var u User

		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.TeamID,
			&u.Role,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {
	const query = `
		UPDATE users
		SET name = @name, username = @username, email = @email, password = @password,
		    team_id = @team_id, role = @role, updated_at = now()
		WHERE id = @id AND deleted_at IS NULL
	`

	args := pgx.NamedArgs{
		"id":       u.ID,
		"name":     u.Name,
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
		"team_id":  u.TeamID,
		"role":     u.Role,
	}

	tx := transaction.FromContext(ctx, r.db)

	tag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE users
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tx := transaction.FromContext(ctx, r.db)

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}

	return nil
}
