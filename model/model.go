package model

import "time"

// Struct to define the columns for the User Details Table
type User struct {
	ID             int        `db:"id" json:"id"`
	PhoneNumber    string     `db:"phone_number" json:"phone_number"`
	Email          string     `db:"email" json:"email"`
	LinkedID       *int       `db:"linked_id" json:"linked_id"`
	LinkPrecedence string     `db:"link_precedence" json:"link_precedence"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`
}
