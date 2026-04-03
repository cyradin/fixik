package incident

import (
	"time"

	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/user"
)

type IncidentID = int64

type Incident struct {
	ID          IncidentID
	Title       string
	Description string

	Status   status.Status
	Priority dict.Entity

	Team   *dict.Entity
	User   *user.User
	Author *user.User

	CreatedAt time.Time
	UpdatedAt time.Time

	CommentsCount int
}

type IncidentList struct {
	Items []Incident
	Total int
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

type CommentID = int64

type Comment struct {
	ID         CommentID
	IncidentID IncidentID
	Author     user.User
	Text       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CommentList struct {
	Items []Comment
	Total int
}
