package core

import (
	"identityreconciliation/api"
	db "identityreconciliation/database"
	"identityreconciliation/model"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Core struct {
	DB  *db.DataBase
	API *api.API
}

func (c *Core) run() {
	// Configure database connection
	if err := c.DB.InitDB("mydb.sqlite"); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer c.DB.CloseDB()

	// Create tables based on struct definitions
	if err := c.DB.CreateTable(model.User{}); err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	// Register API handlers
	//http.HandleFunc("/users", c.API.ReadUsers)
	http.HandleFunc("/users/create", c.API.CreateUser)
	http.HandleFunc("/users/update", c.API.UpdateUser)
	//http.HandleFunc("/users/delete", c.API.DeleteUser)

	// Serve the Swagger UI at /swagger/index.html
	http.Handle("/swagger/", httpSwagger.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
