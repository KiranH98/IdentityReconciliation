package core

import (
	"identityreconciliation/api"
	db "identityreconciliation/database"
	"identityreconciliation/model"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Core struct {
	DB  *db.DataBase
	API *api.API
	log *log.Logger
}

func NewCoreService(db *db.DataBase, apiHandler *api.API) *Core {
	return &Core{
		DB:  db,
		API: apiHandler,
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
	}
}

func (c *Core) CallRun() {
	c.run()
}

func (c *Core) run() {

	c.log.Println("starting core")
	r := mux.NewRouter()

	// Configure database connection
	if err := c.DB.InitDB("identity.db"); err != nil {
		c.log.Fatalf("failed to initialize database: %v", err)
	}
	defer c.DB.CloseDB()

	// Create tables based on struct definitions
	if err := c.DB.CreateTable(model.User{}); err != nil {
		c.log.Fatalf("failed to create table: %v", err)
	}

	// Register API handlers using gorilla/mux
	//r.HandleFunc("/users", c.API.ReadUsers).Methods("GET")
	r.HandleFunc("/users/create", c.API.CreateUser).Methods("POST")
	r.HandleFunc("/users/update", c.API.UpdateUser).Methods("PUT")
	//r.HandleFunc("/users/delete", c.API.DeleteUser).Methods("DELETE")
	r.HandleFunc("/identity", c.API.Identity).Methods("POST")
	// Serve Swagger UI at /swagger/
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	// Use the gorilla/mux router
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
