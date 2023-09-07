package repository

import (
	"database/sql"
	"identityreconciliation/model"
	"sort"
	"time"
)

const (
	Primary   = "primary"
	Secondary = "secondary"
)

// GetUsers retrieves users from the database based on the provided IdentityRequest.
func (repository *Repository) GetContacts(request model.IdentifyRequest) ([]model.Contact, error) {
	var result []model.Contact
	var existing []model.Contact
	var emailRowsexist, phoneNumberRowsExist bool
	// Check if entry exists for email in the incoming request
	if request.Email != "" {
		emailRows, err := repository.GetContactbyEmail(request.Email)
		if err != nil {
			return nil, err
		}
		// Iterate through the rows and scan into Contact model
		existing = append(existing, emailRows...)
		if len(emailRows) != 0 {
			emailRowsexist = true
			repository.log.Println("emailRowsexist ", emailRowsexist)
		}
	}

	// Query for contacts where phoneNumber matches and no entries for email
	if request.PhoneNumber != "" {
		phoneNumberRows, err := repository.GetContactbyPhoneNumber(request.PhoneNumber)
		if err != nil {
			return nil, err
		}

		// Iterate through the rows and scan into Contact structs
		existing = append(existing, phoneNumberRows...)
		if len(phoneNumberRows) != 0 {
			phoneNumberRowsExist = true
			repository.log.Println("phoneNumberRowsExist ", phoneNumberRowsExist)
		}
	}
	// If no rows found, create a new entry in the database
	if len(existing) == 0 {
		repository.log.Println("No entries in DB for the given details")
		repository.log.Println("Creating primary entry in the database")
		// Add a new primary contact entry
		newContact := &model.Contact{
			PhoneNumber:    request.PhoneNumber,
			Email:          request.Email,
			LinkPrecedence: Primary,
		}
		repository.InsertContact(*newContact)
	} else if emailRowsexist && phoneNumberRowsExist {
		// Sort existing contacts by created_at, oldest first
		sortedContacts := make([]model.Contact, len(existing))
		copy(sortedContacts, existing)
		sort.Slice(sortedContacts, func(i, j int) bool {
			if (sortedContacts[i].LinkPrecedence == Primary && (sortedContacts[i].Email == request.Email || sortedContacts[i].PhoneNumber == request.PhoneNumber)) && (sortedContacts[j].LinkPrecedence == Primary && (sortedContacts[j].Email == request.Email || sortedContacts[j].PhoneNumber == request.PhoneNumber)) {
				return sortedContacts[i].CreatedAt.After(sortedContacts[j].CreatedAt)
			}
			return false
		})

		var linkedID sql.NullInt64
		linkedID.Int64 = int64(sortedContacts[1].ID)
		linkedID.Valid = true

		updateContact := sortedContacts[0]
		updateContact.LinkPrecedence = Secondary
		updateContact.LinkedID = linkedID

		repository.log.Println("unliked Primary contact exists")
		repository.log.Println("updating latest entry to secondary")
		repository.UpdateContact(updateContact)
	} else {
		for _, contact := range existing {
			// Check if the current contact matches the request
			if (request.Email != "" && contact.Email == request.Email) || (request.PhoneNumber != "" && contact.PhoneNumber == request.PhoneNumber) {
				// If rows are found, create a "secondary" Contact entry with the new information
				repository.log.Println("Matching entry found in the database")
				repository.log.Println("Creating secondary entry in the database")

				// Create a new Contact struct with the updated information
				var primaryID sql.NullInt64
				for _, contact := range existing {
					if contact.LinkPrecedence == Primary {
						primaryID.Int64 = int64(contact.ID)
						primaryID.Valid = true
						break
					}
				}
				newSecondaryContact := model.Contact{
					PhoneNumber:    request.PhoneNumber,
					Email:          request.Email,
					LinkedID:       primaryID,
					LinkPrecedence: Secondary,
				}
				repository.InsertContact(newSecondaryContact)
			}
		}
	}

	// Create queries based on email and phoneNumber
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

	// Iterate over the different queries
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
			SET linked_id = ?, link_precedence = ?, updated_at = ?
			WHERE id = ?`

	now := time.Now().UTC()
	_, err := repository.db.Exec(stmt, contact.LinkedID, contact.LinkPrecedence, now, contact.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetContactbyEmail(email string) ([]model.Contact, error) {
	var result []model.Contact
	// Query for both email and phoneNumber
	query := "SELECT * FROM contact WHERE email = ?"
	rows, err := repository.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var contact model.Contact
		if err := rows.Scan(&contact.ID, &contact.PhoneNumber, &contact.Email, &contact.LinkedID, &contact.LinkPrecedence, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt); err != nil {
			return nil, err
		}
		result = append(result, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (repository *Repository) GetContactbyPhoneNumber(phoneNumber string) ([]model.Contact, error) {
	var result []model.Contact
	// Query for both email and phoneNumber
	query := "SELECT * FROM contact WHERE phone_number = ?"
	rows, err := repository.db.Query(query, phoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var contact model.Contact
		if err := rows.Scan(&contact.ID, &contact.PhoneNumber, &contact.Email, &contact.LinkedID, &contact.LinkPrecedence, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt); err != nil {
			return nil, err
		}
		result = append(result, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
