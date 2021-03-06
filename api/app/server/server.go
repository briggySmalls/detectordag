package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/briggysmalls/detectordag/api/app/models"
	"github.com/briggysmalls/detectordag/api/app/tokens"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
)

var (
	ErrAccountIDMissing = errors.New("AccountID missing from context")
)

type AccountIdKey struct {
}

type server struct {
	db     database.Client
	shadow shadow.Client
	email  email.Verifier
	iot    iot.Client
	tokens tokens.Tokens
}

type Server interface {
	Auth(w http.ResponseWriter, r *http.Request)
	GetAccount(w http.ResponseWriter, r *http.Request)
	GetDevices(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	UpdateDevice(w http.ResponseWriter, r *http.Request)
}

func New(db database.Client, shadow shadow.Client, email email.Verifier, iot iot.Client, tokens tokens.Tokens) Server {
	return &server{
		db:     db,
		shadow: shadow,
		email:  email,
		iot:    iot,
		tokens: tokens,
	}
}

func SetError(w http.ResponseWriter, err error, status int) {
	// TODO: If 5xx error then hide message unless in debug
	log.Printf("%+v", err)
	// Create the error struct
	m := models.ModelError{
		Error_: err.Error(),
	}
	// Marshal into string
	content, err := json.Marshal(m)
	if err != nil {
		// What do ew
		http.Error(w, "{\"error\": \"Failed to format error message\"}", http.StatusInternalServerError)
		return
	}
	// Write the output
	http.Error(w, string(content), status)
}

func getAccountId(context context.Context) (string, error) {
	// Ensure the auth middleware provided us with the account ID
	accountID := context.Value(AccountIdKey{})
	if accountID == nil {
		return "", ErrAccountIDMissing
	}
	// Cast the value to a string
	accountIDString, ok := accountID.(string)
	if !ok {
		return "", ErrAccountIDMissing
	}
	return accountIDString, nil
}
