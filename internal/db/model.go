package db

import (
	"fmt"
	"time"
)

var ErrNotFound = fmt.Errorf("not found")

type Incident struct {
	ID          int64
	Title       string
	Description string

	PriorityID int64
	StatusID   int64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type DictEntity struct {
	ID          int64
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type (
	Priority = DictEntity
	Status   = DictEntity
	Team     = DictEntity
)

type User struct {
	ID        int64
	Name      string
	Username  string
	Email     string
	Password  string //nolint:gosec
	TeamID    int64
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Role = string

const (
	RoleUser    = "user"
	RoleManager = "manager"
	RoleAdmin   = "admin"
)
