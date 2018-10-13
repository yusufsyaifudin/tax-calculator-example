package model

import "time"

// User is represent data structure in database.
type User struct {
	ID        int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
