package incident

import (
	"time"

	"github.com/cyradin/fixik/internal/dict"
)

type IncidentID = int64

type Incident struct {
	ID          IncidentID
	Title       string
	Description string

	Priority dict.Entity
	Status   dict.Entity

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateIncident struct {
	Title       string
	Description string

	PriorityID dict.EntityID
	StatusID   dict.EntityID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateIncident struct {
	ID IncidentID

	Title       *string
	Description *string

	PriorityID *dict.EntityID
	StatusID   *dict.EntityID
}
