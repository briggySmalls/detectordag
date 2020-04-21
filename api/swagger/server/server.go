package server

import (
	"encoding/json"
	"errors"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/api/swagger/tokens"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"net/http"
)

var (
	ErrAccountIDMissing = errors.New("AccountID missing from context")
)

type server struct {
	db     database.Client
	shadow shadow.Client
	email  email.Client
	tokens tokens.Tokens
}

type Server interface {
	Auth(w http.ResponseWriter, r *http.Request)
	GetAccount(w http.ResponseWriter, r *http.Request)
	GetDevices(w http.ResponseWriter, r *http.Request)
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	UpdateDevice(w http.ResponseWriter, r *http.Request)
}

func New(params Params) Server {
	return &server{
		db:     params.Db,
		shadow: params.Shadow,
		email:  params.Email,
		tokens: params.Tokens,
	}
}

func setError(w http.ResponseWriter, err error, status int) {
	// TODO: If 5xx error then hide message unless in debug
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
