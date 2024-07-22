package entity

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Role         string    `json:"role" db:"role"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
}
