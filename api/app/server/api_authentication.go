package server

import (
	"encoding/json"
	"github.com/briggysmalls/detectordag/api/app/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s *server) Auth(w http.ResponseWriter, r *http.Request) {
	// Try to parse the body
	var creds models.Credentials
	var err error
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		SetError(w, err, http.StatusBadRequest)
		return
	}
	// Query for an account with the given username
	account, err := s.db.GetAccountByUsername(creds.Username)
	if err != nil {
		SetError(w, err, http.StatusForbidden)
		return
	}
	// Check that the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(creds.Password))
	if err != nil {
		SetError(w, err, http.StatusForbidden)
		return
	}
	// Create a token for the authenticated user
	token, err := s.tokens.Create(account.AccountId)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
	}
	// Build response content
	content := models.Token{
		AccountId: account.AccountId,
		Token:     token,
	}
	body, err := json.Marshal(content)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
	}
	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
