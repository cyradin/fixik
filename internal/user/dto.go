package user

import (
	"time"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/role"
)

type ID = int64

type CreateUser struct {
	Name     string
	Username string
	Email    string
	Password string //nolint:gosec
	TeamID   *int64
	Role     role.Type
}

// UpdateUser
type UpdateUser struct {
	ID       int64
	Name     *string
	Username *string
	Email    *string
	Password *string //nolint:gosec
	TeamID   *int64
	Role     *role.Type
}

type User struct {
	ID        ID
	Name      string
	Username  string
	Email     string
	TeamID    *int64
	Role      role.Type
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFromDB(u db.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		TeamID:    u.TeamID,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
