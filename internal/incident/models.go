package incident

import "time"

type Incident struct {
	ID          int64
	Title       string
	Description string

	ImpactID   int64
	PriorityID int64
	StatusID   int64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Priority struct {
	ID   int64
	Code string
	Name string
}

type Impact struct {
	ID   int64
	Code string
	Name string
}

type Status struct {
	ID   int64
	Code string
	Name string
}
