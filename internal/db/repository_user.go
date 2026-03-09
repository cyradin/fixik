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
	const insertUser = `
		INSERT INTO users (
			username,
			email,
			password,
			team_id
		)
		VALUES (
			@username,
			@email,
			@password,
			@team_id
		)
		RETURNING id
	`

	args := pgx.NamedArgs{
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
		"team_id":  u.TeamID,
	}

	tx := transaction.FromContext(ctx, r.db)

	if err := tx.QueryRow(ctx, insertUser, args).Scan(&u.ID); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	if err := r.insertRoles(ctx, tx, u.ID, u.RoleIDs); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (User, error) {
	const query = `
		SELECT
			u.id,
			u.username,
			u.email,
			u.password,
			u.team_id,
			COALESCE(array_agg(ur.role_id) FILTER (WHERE ur.role_id IS NOT NULL), '{}')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		WHERE u.id = $1
		GROUP BY u.id
	`

	var u User

	err := transaction.FromContext(ctx, r.db).QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.TeamID,
		&u.RoleIDs,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("db query: %w", err)
	}

	return u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]User, error) {
	const query = `
		SELECT
			u.id,
			u.username,
			u.email,
			u.password,
			u.team_id,
			COALESCE(array_agg(ur.role_id) FILTER (WHERE ur.role_id IS NOT NULL), '{}')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		GROUP BY u.id
		ORDER BY u.id
	`

	rows, err := transaction.FromContext(ctx, r.db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.TeamID,
			&u.RoleIDs,
		); err != nil {
			return nil, fmt.Errorf("scan value: %w", err)
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {
	const query = `
		UPDATE users
		SET
			username = @username,
			email = @email,
			password = @password,
			team_id = @team_id
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
		"team_id":  u.TeamID,
	}

	tx := transaction.FromContext(ctx, r.db)

	if _, err := tx.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, u.ID); err != nil {
		return fmt.Errorf("delete roles: %w", err)
	}

	if err := r.insertRoles(ctx, tx, u.ID, u.RoleIDs); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM users
		WHERE id = @id
	`

	args := pgx.NamedArgs{"id": id}

	if _, err := transaction.FromContext(ctx, r.db).Exec(ctx, query, args); err != nil {
		return fmt.Errorf("db query: %w", err)
	}

	return nil
}

func (r *UserRepository) insertRoles(ctx context.Context, tx transaction.Postgres, userID int64, roleIDs []int64) error {
	const query = `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
	`

	for _, roleID := range roleIDs {
		if _, err := tx.Exec(ctx, query, userID, roleID); err != nil {
			return fmt.Errorf("insert role: %w", err)
		}
	}

	return nil
}
