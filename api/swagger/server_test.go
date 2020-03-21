package swagger

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestCreateToken(t *testing.T) {
	// Create some test inputs
	testParams := []struct {
		secret    string
		duration  string
		accountID string
		error     error
	}{
		{secret: "mysecret", duration: "2h", accountID: ""},
	}
	for _, params := range testParams {
		// Create a server
		srv := server{
			config: Config{JwtSecret: params.secret, JwtDuration: params.duration},
		}
		// Create a token
		ss, err := srv.createToken(params.accountID)
		// Check the error
		if err != params.error {
			t.Fatalf(err.Error())
		}
		// Parse the token contents
		token, err := jwt.ParseWithClaims(ss, &CustomAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret
			return []byte(params.secret), nil
		})
		// Check the claims
		if claims, ok := token.Claims.(*CustomAuthClaims); ok && token.Valid {
			if claims.AccountId != params.accountID {
				t.Fatalf("Token did not save correct account ID")
			}
		} else {
			t.Fatalf(err.Error())
		}
	}
}
