package incident

import (
	"time"

	"github.com/cyradin/fixik/internal/priority"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/team"
	"github.com/cyradin/fixik/internal/user"
)

type ID = int64

type Incident struct {
	ID          ID
	Title       string
	Description string

	Status   status.Status
	Priority priority.Priority

	Team   *team.Team
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

	PriorityID priority.ID
	StatusID   status.ID

	TeamID   *team.ID
	UserID   *user.ID
	AuthorID *user.ID
}

type UpdateIncident struct {
	ID ID

	Title       *string
	Description *string

	PriorityID *priority.ID
	StatusID   *status.ID

	TeamID   *team.ID
	UserID   *user.ID
	AuthorID *user.ID
}

type CommentID = int64

type Comment struct {
	ID         CommentID
	IncidentID ID
	Author     user.User
	Text       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CommentList struct {
	Items []Comment
	Total int
}

type Filter struct {
	AuthorIDs   []user.ID
	UserIDs     []user.ID
	TeamIDs     []team.ID
	PriorityIDs []priority.ID
	StatusIDs   []status.ID
	ActiveOnly  bool
}
