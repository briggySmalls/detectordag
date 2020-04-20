package server

import (
	"encoding/json"
	"errors"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

const (
	AuthenticationHeaderPrefix = "Bearer "
)

var (
	ErrNoAuthHeader           = errors.New("Authorization header not set")
	ErrMalformattedAuthHeader = errors.New("Authorization header badly formed")
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

type CustomAuthClaims struct {
	AccountId string `json:"accountId"`
	jwt.StandardClaims
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

func (s *server) createToken(accountId string) (string, error) {
	// Get the duration tokens should be alive for
	dur, err := s.config.ParseDuration()
	if err != nil {
		return "", err
	}
	// Create the Claims
	claims := CustomAuthClaims{
		accountId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(dur).Unix(),
			Issuer:    issuer,
		},
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JwtSecret))
}

// checkAuthorized checks that the token authorises access to the specified account
func (s *server) checkAuthorized(tokenString, accountId string) error {
	// Parse takes the token string and a function for looking up the key.
	token, err := jwt.ParseWithClaims(tokenString, &CustomAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return our secret
		return []byte(s.config.JwtSecret), nil
	})
	if err != nil {
		return err
	}
	// Check the token contents
	if claims, ok := token.Claims.(*CustomAuthClaims); ok && token.Valid {
		// Confirm the account IDs match
		if claims.AccountId != accountId {
			return fmt.Errorf("Not authorized to access account: %s", accountId)
		}
		// The token authorises access
		return nil
	}
	return err
}

func (s *server) getToken(header *http.Header) (string, error) {
	// Check the auth header is set
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}
	// Ensure we've been given a JWT how we expect
	if !strings.HasPrefix(authHeader, AuthenticationHeaderPrefix) {
		return "", ErrMalformattedAuthHeader
	}
	// Return the token
	return strings.TrimPrefix(authHeader, AuthenticationHeaderPrefix), nil
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
