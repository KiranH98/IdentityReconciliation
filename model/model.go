package model

import "time"

// Struct to define the columns for the User Details Table
type User struct {
	ID             int        `db:"id"`
	PhoneNumber    string     `db:"phone_number"`
	Email          string     `db:"email"`
	LinkedID       *int       `db:"linked_id"`
	LinkPrecedence string     `db:"link_precedence"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}
