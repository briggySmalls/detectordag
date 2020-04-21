package server

import (
	"encoding/json"
	"errors"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

type server struct {
	db     database.Client
	shadow shadow.Client
	email  email.Client
	config Config
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
		config: params.Config,
	}
}

func (s *server) validateAccount(w http.ResponseWriter, r *http.Request) *string {
	// Ensure that there is a token sent
	token, err := s.getToken(&r.Header)
	if err != nil {
		setError(w, err, http.StatusUnauthorized)
		return nil
	}
	// Pull out the account ID
	vars := mux.Vars(r)
	accountId, ok := vars["accountId"]
	if !ok {
		setError(w, errors.New("Account ID not supplied in path"), http.StatusBadRequest)
		return nil
	}
	// Check the user is authorised
	err = s.checkAuthorized(token, accountId)
	if err != nil {
		setError(w, err, http.StatusForbidden)
		return nil
	}
	return &accountId
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
