package status

import "time"

type StatusID = int64

type Status struct {
	ID          StatusID
	Name        string
	Code        string
	Description string
	Sort        int
	IsFinal     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
