package model

import "time"

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusCanceled   = "canceled"
) //Other statuses for future improvements

var allowedStatuses = map[string]struct{}{
	StatusPending:    {},
	StatusInProgress: {},
	StatusDone:       {},
	StatusCanceled:   {},
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func IsValidStatus(s string) bool {
	_, ok := allowedStatuses[s]

	return ok
}
