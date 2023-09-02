package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db  *sql.DB
	log *log.Logger
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:  db,
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
	}
}

// InitDB initializes the database connection.
func (repository *Repository) InitDB(dbName string) error {

	connString := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbName)
	database, err := sql.Open("sqlite3", connString)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	repository.db = database

	return nil
}

// CloseDB closes the database connection.
func (repository *Repository) CloseDB() {
	if repository.db != nil {
		repository.db.Close()
	}
}

// CreateTable creates a table in the database based on the model definition.
func (repository *Repository) CreateTable(data interface{}) error {
	query := generateCreateTableQuery(data)
	if query == "" {
		return fmt.Errorf("failed to generate create table query")
	}
	_, err := repository.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

// generateCreateTableQuery generates the SQL query to create a table based on the model definition.
func generateCreateTableQuery(data interface{}) string {
	// Extract the table name from the struct
	tableName := getTableName(data)

	// Get the fields and their types from the struct
	//fields := getStructFields(data)

	// Generate SQL CREATE TABLE statement
	createTableSQL := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	fields := reflect.TypeOf(data)
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		dbTag := field.Tag.Get("db")
		fieldType := getFieldType(field.Type)
		createTableSQL += dbTag + " " + fieldType + ","
	}
	createTableSQL = strings.TrimSuffix(createTableSQL, ",") + ");"

	return createTableSQL
}

// getTableName gets the table name from the model definition.
func getTableName(data interface{}) string {
	structName := fmt.Sprintf("%T", data)
	lastDotIndex := strings.LastIndex(structName, ".")
	if lastDotIndex != -1 {
		structName = structName[lastDotIndex+1:]
	}
	return strings.ToLower(structName)
}

// getStructFields gets the fields and their types from the model definition.
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
		return "INTEGER"
	case reflect.Struct:
		// Handle time.Time as a special case
		if t == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP"
		}
	}
	return ""
}
