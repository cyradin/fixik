package db

import (
	"fmt"
	"time"
)

var ErrNotFound = fmt.Errorf("not found")

type IncidentFilter struct {
	AuthorIDs   []int64
	UserIDs     []int64
	TeamIDs     []int64
	PriorityIDs []int64
	StatusIDs   []int64
}

type Incident struct {
	ID          int64
	Title       string
	Description string

	PriorityID int64
	StatusID   int64

	TeamID   *int64
	UserID   *int64
	AuthorID *int64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	CommentsCount int
}

type IncidentListResult struct {
	Items []Incident
	Total int
}

type Comment struct {
	ID         int64
	AuthorID   int64
	IncidentID int64
	Text       string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CommentListResult struct {
	Items []Comment
	Total int
}

type DictEntity struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Sort        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type (
	Priority = DictEntity
	Team     = DictEntity
)

type User struct {
	ID        int64
	Name      string
	Username  string
	Email     string
	Password  string //nolint:gosec
	TeamID    *int64
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

type Status struct {
	ID          int64
	Code        string
	Name        string
	Description string
	Sort        int
	IsFinal     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
