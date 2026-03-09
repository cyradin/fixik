package incident

import "time"

type StatusID = int64

type Status struct {
	ID   StatusID
	Name string
	Code string
}

type PriorityID = int64

type Priority struct {
	ID   PriorityID
	Name string
	Code string
}

type ImpactID = int64

type Impact struct {
	ID   ImpactID
	Name string
	Code string
}

type IncidentID = int64

type Incident struct {
	ID          IncidentID
	Title       string
	Description string

	Impact   Impact
	Priority Priority
	Status   Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateIncident struct {
	Title       string
	Description string

	ImpactID   ImpactID
	PriorityID PriorityID
	StatusID   StatusID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateIncident struct {
	ID IncidentID

	Title       string
	Description string

	ImpactID   ImpactID
	PriorityID PriorityID
	StatusID   StatusID
}
