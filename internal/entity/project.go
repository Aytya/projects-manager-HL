package entity

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`

	// type: string
	// format: date-time
	FinishedAt sql.NullTime `db:"finished_at" json:"finished_at"`
	Manager    string       `db:"manager" json:"manager"`
}
