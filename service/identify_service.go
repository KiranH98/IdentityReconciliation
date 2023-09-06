package service

import (
	"encoding/json"
	"identityreconciliation/model"
	"log"
	"net/http"
)

const (
	Primary   = "primary"
	Secondary = "secondary"
)

// Identity Return user data
// @Summary Return user data
// @Description This endpoint is used to return data related to email or phone number supplied
// @Accept json
// @Produce json
// @Param user body model.IdentifyRequest true "enter email and phone number"
// @Success 200 {object} model.IdentifyResponse
// @Router /identify [post]
func (s *Service) Identify(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming request")
	// Parse the JSON request
	var request model.IdentifyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.log.Fatal("Error occured while parsing input request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get users from db based on identity request
	contacts, err := s.storage.GetContacts(request)
	if err != nil {
		s.log.Fatal("error occured while fetching details drom DB", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contact := &model.ContactResponse{
		PrimaryContactID:    getPrimaryID(contacts),
		Emails:              getEmails(contacts),
		PhoneNumbers:        getPhoneNumbers(contacts),
		SecondaryContactIDs: getSecondaryContactIDs(contacts),
	}

	identityResp := &model.IdentifyResponse{
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

func getEmails(contacts []model.Contact) []string {
	var emails []string
	for _, contact := range contacts {
		emails = append(emails, contact.Email)
	}
	uniqueEmails, ok := RemoveDuplicates(emails).([]string)
	if !ok {
		return nil
	}
	return uniqueEmails
}

func getPhoneNumbers(contacts []model.Contact) []string {
	var phoneNumbers []string
	for _, contact := range contacts {
		phoneNumbers = append(phoneNumbers, contact.PhoneNumber)
	}
	uniquePhoneNumbers, ok := RemoveDuplicates(phoneNumbers).([]string)
	if !ok {
		return nil
	}
	return uniquePhoneNumbers
}

func getSecondaryContactIDs(Contacts []model.Contact) []int {
	var secondaryContactIDs []int
	for _, contact := range Contacts {
		if contact.LinkPrecedence == "secondary" {
			secondaryContactIDs = append(secondaryContactIDs, contact.ID)
		}
	}
	uniqueSecondaryIDs, ok := RemoveDuplicates(secondaryContactIDs).([]int)
	if !ok {
		return nil
	}
	return uniqueSecondaryIDs
}

func getPrimaryID(Contacts []model.Contact) int {
	var primaryID int
	for _, contact := range Contacts {
		if contact.LinkPrecedence == "primary" {
			primaryID = contact.ID
		}
	}
	return primaryID
}

func RemoveDuplicates(input interface{}) interface{} {
	switch input := input.(type) {
	case []string:
		uniqueMap := make(map[string]struct{})
		var result []string

		for _, str := range input {
			if _, found := uniqueMap[str]; !found {
				uniqueMap[str] = struct{}{}
				result = append(result, str)
			}
		}

		return result

	case []int:
		uniqueMap := make(map[int]struct{})
		var result []int

		for _, num := range input {
			if _, found := uniqueMap[num]; !found {
				uniqueMap[num] = struct{}{}
				result = append(result, num)
			}
		}

		return result

	default:
		return nil
	}
}
