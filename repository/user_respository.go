package repository

import (
	"identityreconciliation/model"
	"time"
)

// GetUsers retrieves users from the database based on the provided IdentityRequest.
func (repository *Repository) GetUsers(request model.IdentityRequest) ([]model.User, error) {
	var result []model.User

	// creating the sql query to be executed
	query := "SELECT * FROM users WHERE ("
	args := []interface{}{}

	if request.Email != "" {
		query += " email = ?"
		args = append(args, request.Email)
	}

	if request.PhoneNumber != "" {
		if len(args) > 0 && request.Email != "" {
			query += " AND"
		}
		query += " phoneNumber = ?"
		args = append(args, request.PhoneNumber)
	}

	query += ")"

	// Execute the SQL query
	rows, err := repository.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and scan into User structs
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Email, &user.PhoneNumber); err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	// Check if no rows were found, and return an empty slice
	if len(result) == 0 {
		return []model.User{}, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// InsertUser inserts a new user into the database and returns the ID of the created user.
func (repository *Repository) InsertUser(user model.User) (int, error) {
	stmt := `INSERT INTO users (phone_number, email, linked_id, link_precedence, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now().UTC()
	result, err := repository.db.Exec(stmt, user.PhoneNumber, user.Email, user.LinkedID, user.LinkPrecedence, now, now)
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
func (repository *Repository) UpdateUser(user model.User) error {
	stmt := `UPDATE users
			SET phone_number = ?, email = ?, linked_id = ?, link_precedence = ?, updated_at = ?
			WHERE id = ?`

	now := time.Now().UTC()
	_, err := repository.db.Exec(stmt, user.PhoneNumber, user.Email, user.LinkedID, user.LinkPrecedence, now, user.ID)
	if err != nil {
		return err
	}

	return nil
}
