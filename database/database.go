package db

import (
	"database/sql"
	"fmt"
	"identityreconciliation/model"
	"reflect"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	db *sql.DB
}

// InitDB initializes the database connection.
func (DB *DataBase) InitDB(dbName string) error {

	connString := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbName)
	database, err := sql.Open("sqlite3", connString)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	DB.db = database

	return nil
}

// CloseDB closes the database connection.
func (DB *DataBase) CloseDB() {
	if DB.db != nil {
		DB.db.Close()
	}
}

// CreateTable creates a table in the database based on the struct definition.
func (DB *DataBase) CreateTable(data interface{}) error {
	query := generateCreateTableQuery(data)
	if query == "" {
		return fmt.Errorf("failed to generate create table query")
	}
	_, err := DB.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

// generateCreateTableQuery generates the SQL query to create a table based on the struct definition.
func generateCreateTableQuery(data interface{}) string {
	// Extract the table name from the struct
	tableName := getTableName(data)

	// Get the fields and their types from the struct
	fields := getStructFields(data)

	// Generate the CREATE TABLE query
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, fields)

	return query
}

// getTableName gets the table name from the struct definition.
func getTableName(data interface{}) string {
	structName := fmt.Sprintf("%T", data)
	lastDotIndex := strings.LastIndex(structName, ".")
	if lastDotIndex != -1 {
		structName = structName[lastDotIndex+1:]
	}
	return structName
}

// getStructFields gets the fields and their types from the struct definition.
func getStructFields(data interface{}) string {
	fields := ""
	structType := reflect.TypeOf(data)
	numFields := structType.NumField()

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		if fields != "" {
			fields += ", "
		}
		fields += field.Name + " " + getFieldType(field.Type)
	}

	return fields
}

// getFieldType returns the corresponding SQL type based on the Go type.
func getFieldType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int64:
		return "INTEGER"
	case reflect.Float64:
		return "REAL"
	case reflect.String:
		return "TEXT"
	case reflect.Bool:
		return "INTEGER" // SQLite does not have a BOOLEAN type, so we use INTEGER (0 or 1)
	case reflect.Struct:
		// Handle time.Time as a special case
		if t == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP"
		}
	}
	return ""
}

// GetUsers retrieves users from the database based on the provided IdentityRequest.
func (DB *DataBase) GetUsers(request model.IdentityRequest) ([]model.User, error) {
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
	rows, err := DB.db.Query(query, args...)
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
