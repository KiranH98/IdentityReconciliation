package api

import (
	"encoding/json"
	db "identityreconciliation/database"
	"identityreconciliation/model"
	"log"
	"net/http"
	"time"
)

type API struct {
	DB *db.DataBase
}

type UserSwaggo struct {
	ID             int        `db:"id" json:"id"`
	PhoneNumber    string     `db:"phone_number" json:"phone_number"`
	Email          string     `db:"email" json:"email"`
	LinkedID       *int       `db:"linked_id" json:"linked_id"`
	LinkPrecedence string     `db:"link_precedence" json:"link_precedence"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`
}

// CreateUser creates a new user.
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Accept json
// @Produce json
// @Param user body UserSwaggo true "User object"
// @Success 200 {object} UserSwaggo
// @Router /users/create [post]
func (api *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the CreatedAt and UpdatedAt timestamps
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	userID, err := api.DB.InsertUser(user)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.ID = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates an existing user.
// @Summary Update an existing user
// @Description Update an existing user with the provided data
// @Accept json
// @Produce json
// @Param user body UserSwaggo true "User object"
// @Success 200
// @Router /users/update [put]
func (api *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the UpdatedAt timestamp
	now := time.Now().UTC()
	user.UpdatedAt = now

	err = api.DB.UpdateUser(user)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *API) Identity(w http.ResponseWriter, r *http.Request) {

}

/* func (api *API) ReadUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.DB.GetUsers()
	if err != nil {
		log.Printf("failed to get users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
} */

/* func (api *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Parse the ID from the URL path
	id := r.URL.Path[len("/users/delete/"):]

	err := api.DB.DeleteUser(id)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} */
