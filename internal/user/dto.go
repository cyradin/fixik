package user

import (
	"time"

	"github.com/cyradin/fixik/internal/db"
)

type Role struct {
	Code        RoleType
	Name        string
	Description string
}

type RoleType = db.Role

const (
	RoleUser    = db.RoleUser
	RoleManager = db.RoleManager
	RoleAdmin   = db.RoleAdmin
)

func Roles() []Role {
	return []Role{
		{
			Name:        "Пользователь",
			Code:        RoleUser,
			Description: "Может работать с инцидентами",
		},
		{
			Name:        "Менеджер",
			Code:        RoleManager,
			Description: "Может выполнять все операции, кроме изменения пользователей и команд",
		},
		{
			Name:        "Администратор",
			Code:        RoleAdmin,
			Description: "Может выполнять все операции",
		},
	}
}

func RoleTypes() []RoleType {
	return []RoleType{
		RoleUser,
		RoleManager,
		RoleAdmin,
	}
}

// CreateUser
type CreateUser struct {
	Name     string
	Username string
	Email    string
	Password string //nolint:gosec
	TeamID   int64
	Role     RoleType
}

// UpdateUser
type UpdateUser struct {
	ID       int64
	Name     *string
	Username *string
	Email    *string
	Password *string //nolint:gosec
	TeamID   *int64
	Role     *RoleType
}

type User struct {
	ID        int64
	Name      string
	Username  string
	Email     string
	TeamID    int64
	Role      RoleType
	CreatedAt time.Time
	UpdatedAt time.Time
}
