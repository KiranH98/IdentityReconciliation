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
	DB  *db.DataBase
	log *log.Logger
}

// CreateUser creates a new user.
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Accept json
// @Produce json
// @Param user body model.UserSwaggo true "User object"
// @Success 200 {object} model.UserSwaggo
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
// @Param user body model.UserSwaggo true "User object"
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

// Identity Return user data
// @Summary Return user data
// @Description This endpoint is used to return data related to email or phone number supplied
// @Accept json
// @Produce json
// @Param user body model.IdentityRequest true "enter email and phone number"
// @Success 200 {object} model.IdentityResponse
// @Router /users/create [post]
func (api *API) Identity(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request
	var request model.IdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		api.log.Fatal("Error occured while parsing input request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get users from db based on identity request in put

	users, err := api.DB.GetUsers(request)
	if err != nil {
		api.log.Fatal("error occured while fetching details drom DB", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(users) == 0 {
		api.log.Println("No entries in DB for the given details")
		api.log.Println("Creating entry in to the database")
		//todo addusers
	}

	contact := &model.Contact{
		PrimaryContactID:    getPrimaryID(users),
		Emails:              getEmails(users),
		PhoneNumbers:        getPhoneNumbers(users),
		SecondaryContactIDs: getSecondaryContactIDs(users),
	}

	identityResp := &model.IdentityResponse{
		Contact: *contact,
	}

	// Encode and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(identityResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getEmails(users []model.User) []string {
	var emails []string
	for _, user := range users {
		emails = append(emails, user.Email)
	}
	return emails
}

func getPhoneNumbers(users []model.User) []string {
	var phoneNumbers []string
	for _, user := range users {
		phoneNumbers = append(phoneNumbers, user.PhoneNumber)
	}
	return phoneNumbers
}

func getSecondaryContactIDs(users []model.User) []int {
	var secondaryContactIDs []int
	for _, user := range users {
		if user.LinkPrecedence == "secondary" {
			secondaryContactIDs = append(secondaryContactIDs, user.ID)
		}
	}
	return secondaryContactIDs
}

func getPrimaryID(users []model.User) int {
	var primaryID int
	for _, user := range users {
		if user.LinkPrecedence == "primary" {
			primaryID = user.ID
		}
	}
	return primaryID
}
