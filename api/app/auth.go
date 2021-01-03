package app

import (
	"context"
	"errors"
	"github.com/briggysmalls/detectordag/api/app/server"
	"github.com/briggysmalls/detectordag/api/app/tokens"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const (
	AuthenticationHeaderPrefix = "Bearer "
)

var (
	ErrNoAuthHeader           = errors.New("Authorization header not set")
	ErrMalformattedAuthHeader = errors.New("Authorization header badly formed")
	// Internal because gorilla should catch this
	errPathParameterMissing = errors.New("Path parameter missing")
)

type auth struct {
	tokens tokens.Tokens
	iot    iot.Client
}

type accountFetcher func(r *http.Request) (string, error)

// Middleware for authorizing requests
func (a *auth) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure there is a token
		token, err := getToken(&r.Header)
		if err != nil {
			server.SetError(w, err, http.StatusInternalServerError)
		}
		// Check that the token is valid
		tokenAccountID, err := a.tokens.Validate(token)
		switch err {
		case tokens.ErrBadToken:
			server.SetError(w, err, http.StatusForbidden)
			return
		case tokens.ErrInternalError:
			server.SetError(w, err, http.StatusInternalServerError)
			return
		default:
			break
		}
		// Fetch the account associated with the resource request
		accountID, err := a.getAccount(r)
		// Ensure we were able to get the account
		switch err {
		case nil:
			break
		case errPathParameterMissing:
			// We shouldn't ever get this, gorilla should handle it
			server.SetError(w, err, http.StatusInternalServerError)
			return
		default:
			// Something else went wrong, e.g. accessing database
			server.SetError(w, err, http.StatusInternalServerError)
			return
		}
		// Ensure we are authorised to access the account's resources
		if accountID != tokenAccountID {
			// Don't reveal that the account exists by returning 403
			server.SetError(w, errors.New("Not found"), http.StatusNotFound)
			return
		}
		// Record the account ID in the context, in case people want it
		ctx := context.WithValue(r.Context(), server.AccountIdKey{}, accountID)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function for fetching the accountID requested
func (a *auth) getAccount(r *http.Request) (string, error) {
	// Create a map of prefixes to match
	matchers := map[string]accountFetcher{
		"/v1/accounts": getAccountFromVars,
		"/v1/devices":  a.getAccountFromDevice,
	}
	// Try to find appropriate fetcher for the route
	for prefix, fetcher := range matchers {
		if strings.HasPrefix(r.URL.Path, prefix) {
			return fetcher(r)
		}
	}
	return "", errors.New("No match found")
}

func getAccountFromVars(r *http.Request) (string, error) {
	// Pull out the account ID
	vars := mux.Vars(r)
	accountID, ok := vars["accountId"]
	if !ok {
		return "", errPathParameterMissing
	}
	return accountID, nil
}

func (a *auth) getAccountFromDevice(r *http.Request) (string, error) {
	// Pull out the device ID
	vars := mux.Vars(r)
	deviceID, ok := vars["deviceId"]
	if !ok {
		return "", errPathParameterMissing
	}
	// Lookup the device in the database
	d, err := a.iot.GetThing(deviceID)
	if err != nil {
		return "", err
	}
	return d.AccountId, nil
}

func getToken(header *http.Header) (string, error) {
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
