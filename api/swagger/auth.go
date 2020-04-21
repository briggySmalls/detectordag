package swagger

import (
	"context"
	"github.com/briggysmalls/detectordag/api/swagger/tokens"
	"net/http"
)

const (
	AuthenticationHeaderPrefix = "Bearer "
)

var (
	ErrNoAuthHeader           = errors.New("Authorization header not set")
	ErrMalformattedAuthHeader = errors.New("Authorization header badly formed")
)

type auth struct {
	tokens tokens.Tokens
}

// Middleware for authorizing requests
func (a *auth) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure there is a token
		token, err := getToken(r.Header)
		if err != nil {
			setError()
		}
		// Check the user is authorised
		accountID, err = a.tokens.Validate(token)
		if err != nil {
			setError(w, err, http.StatusForbidden)
			return nil
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
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
