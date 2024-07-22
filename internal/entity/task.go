package entity

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          string       `json:"id" db:"id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Priority    string       `json:"priority" db:"priority"`
	Status      string       `json:"status" db:"status"`
	Assignee    string       `json:"assignee" db:"assignee"`
	Project     string       `json:"project" db:"project"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	FinishedAt  sql.NullTime `json:"finished_at" db:"finished_at"`
}
