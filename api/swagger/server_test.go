package swagger

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	timeLeeway, err := time.ParseDuration("10ns")
	if err != nil {
		t.Fatalf(err.Error())
	}
	// Create some test inputs
	testParams := []struct {
		secret    string
		duration  string
		accountID string
		error     error
	}{
		{secret: "mysecret", duration: "2h", accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B"},
		{secret: "anothersecret", duration: "1m", accountID: "22222222-32C8-4908-8377-2E6A021D3D2B"},
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
			// Confirm the accounts match
			if claims.AccountId != params.accountID {
				t.Fatalf("Token did not save correct account ID")
			}
			// Confirm the expiry time
			dur, err := time.ParseDuration(params.duration)
			if err != nil {
				t.Fatalf(err.Error())
			}
			exp := time.Unix(claims.ExpiresAt, 0)
			if math.Abs(float64(time.Until(exp)-dur)) < float64(timeLeeway) {
				t.Fatalf("Token did not save correct duration")
			}
		} else {
			t.Fatalf(err.Error())
		}
	}
}
