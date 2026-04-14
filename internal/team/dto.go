package team

import "time"

type ID = int64

type Team struct {
	ID          ID
	Name        string
	Code        string
	Description string
	Sort        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
