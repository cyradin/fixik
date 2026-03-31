package dict

import "time"

type EntityID = int64

type Entity struct {
	ID          EntityID
	Name        string
	Code        string
	Description string
	Sort        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
