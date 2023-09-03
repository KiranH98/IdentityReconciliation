package repository

import (
	"identityreconciliation/model"
	"time"
)

// GetUsers retrieves users from the database based on the provided IdentityRequest.
func (repository *Repository) GetContacts(request model.IdentifyRequest) ([]model.Contact, error) {
	var result []model.Contact

	// creating the sql query to be executed
	queries := []string{}
	args := [][]interface{}{}

	if request.Email != "" && request.Email != "null" {
		// Query for email only
		queries = append(queries, "SELECT * FROM contact WHERE email = ?")
		args = append(args, []interface{}{request.Email})
	}

	if request.PhoneNumber != "" && request.PhoneNumber != "null" {
		// Query for phone number only
		queries = append(queries, "SELECT * FROM contact WHERE phone_number = ?")
		args = append(args, []interface{}{request.PhoneNumber})
	}

	if request.Email != "" && request.Email != "null" && request.PhoneNumber != "" && request.PhoneNumber != "null" {
		// Query for both email and phone number
		queries = append(queries, "SELECT * FROM contact WHERE email = ? AND phone_number = ?")
		args = append(args, []interface{}{request.Email, request.PhoneNumber})
	}

	repository.log.Printf("Query: %v, Args: %v\n", queries, args)

	//iterate over the differnt queries
	for i, query := range queries {
		// Execute the SQL query with the args slice
		rows, err := repository.db.Query(query, args[i]...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// Iterate through the rows and scan into User structs
		for rows.Next() {
			var contact model.Contact
			if err := rows.Scan(&contact.ID, &contact.PhoneNumber, &contact.Email, &contact.LinkedID, &contact.LinkPrecedence, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt); err != nil {
				return nil, err
			}
			result = append(result, contact)
		}

		// Check if no rows were found, and return an empty slice
		if len(result) == 0 {
			repository.log.Println("no records found")
			return []model.Contact{}, nil
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// InsertUser inserts a new user into the database and returns the ID of the created user.
func (repository *Repository) InsertContact(contact model.Contact) (int, error) {
	stmt := `INSERT INTO contact (phone_number, email, linked_id, link_precedence, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now().UTC()
	result, err := repository.db.Exec(stmt, contact.PhoneNumber, contact.Email, contact.LinkedID, contact.LinkPrecedence, now, now)
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
func (repository *Repository) UpdateContact(contact model.Contact) error {
	stmt := `UPDATE contact
			SET phone_number = ?, email = ?, linked_id = ?, link_precedence = ?, updated_at = ?
			WHERE id = ?`

	now := time.Now().UTC()
	_, err := repository.db.Exec(stmt, contact.PhoneNumber, contact.Email, contact.LinkedID, contact.LinkPrecedence, now, contact.ID)
	if err != nil {
		return err
	}

	return nil
}
