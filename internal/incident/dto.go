package incident

import (
	"time"

	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/user"
)

type IncidentID = int64

type Incident struct {
	ID          IncidentID
	Title       string
	Description string

	Priority dict.Entity
	Status   dict.Entity

	Team   *dict.Entity
	User   *user.User
	Author *user.User

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateIncident struct {
	Title       string
	Description string

	PriorityID dict.EntityID
	StatusID   dict.EntityID

	TeamID   *int64
	UserID   *int64
	AuthorID *int64
}

type UpdateIncident struct {
	ID IncidentID

	Title       *string
	Description *string

	TeamID   *int64
	UserID   *int64
	AuthorID *int64

	PriorityID *dict.EntityID
	StatusID   *dict.EntityID
}
