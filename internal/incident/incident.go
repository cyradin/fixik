package incident

import "time"

type Incident struct {
	ID          int64
	Title       string
	Description string
	Impact      string
	Priority    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
