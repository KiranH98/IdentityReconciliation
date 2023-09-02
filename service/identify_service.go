package service

import (
	"encoding/json"
	"identityreconciliation/model"
	"log"
	"net/http"
)

// Identity Return user data
// @Summary Return user data
// @Description This endpoint is used to return data related to email or phone number supplied
// @Accept json
// @Produce json
// @Param user body model.IdentityRequest true "enter email and phone number"
// @Success 200 {object} model.IdentityResponse
// @Router /identify [post]
func (s *Service) Identify(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming request")
	// Parse the JSON request
	var request model.IdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.log.Fatal("Error occured while parsing input request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get users from db based on identity request in put

	users, err := s.storage.GetUsers(request)
	if err != nil {
		s.log.Fatal("error occured while fetching details drom DB", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(users) == 0 {
		s.log.Println("No entries in DB for the given details")
		s.log.Println("Creating entry in to the database")
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