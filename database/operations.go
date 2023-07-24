package db

import (
	"identityreconciliation/model"
	"time"
)

// InsertUser inserts a new user into the database and returns the ID of the created user.
func (DB *DataBase) InsertUser(user model.User) (int, error) {
	stmt := `INSERT INTO users (phone_number, email, linked_id, link_precedence, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now().UTC()
	result, err := DB.db.Exec(stmt, user.PhoneNumber, user.Email, user.LinkedID, user.LinkPrecedence, now, now)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

// UpdateUser updates an existing user in the database.
func (DB *DataBase) UpdateUser(user model.User) error {
	stmt := `UPDATE users
			SET phone_number = ?, email = ?, linked_id = ?, link_precedence = ?, updated_at = ?
			WHERE id = ?`

	now := time.Now().UTC()
	_, err := DB.db.Exec(stmt, user.PhoneNumber, user.Email, user.LinkedID, user.LinkPrecedence, now, user.ID)
	if err != nil {
		return err
	}

	return nil
}
