package db

import (
	"database/sql"
	"fmt"
	"reflect"
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
