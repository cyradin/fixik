package user

import (
	"time"

	"github.com/cyradin/fixik/internal/dict"
)

// CreateUser
type CreateUser struct {
	Username string
	Email    string
	Password string //nolint:gosec
	TeamID   int64
	RoleIDs  []int64
}

// UpdateUser
type UpdateUser struct {
	ID       int64
	Username *string
	Email    *string
	Password *string //nolint:gosec
	TeamID   *int64
	RoleIDs  *[]int64
}

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string //nolint:gosec
	TeamID    int64
	Roles     []dict.Entity
	CreatedAt time.Time
	UpdatedAt time.Time
}
