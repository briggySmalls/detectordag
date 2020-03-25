package swagger

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math"
	"net/http"
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

func TestCheckAuthorized(t *testing.T) {
	// Create some test inputs
	testParams := []struct {
		secret    string
		token     string
		now       time.Time
		accountID string
		errors    uint32
	}{
		{
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjE1ODQ3OTk0NzQsImlzcyI6ImRldGVjdG9yZGFnIn0.qqMDypPk5BT1dz_8KT6S9eNLABWcYIfnaRr_BroisKo",
			now:       createTime(t, "2020/03/21 12:06:00"),
			errors:    0,
		},
		{
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjE1ODQ3OTk0NzQsImlzcyI6ImRldGVjdG9yZGFnIn0.qqMDypPk5BT1dz_8KT6S9eNLABWcYIfnaRr_BroisKo",
			now:       createTime(t, "2020/03/22 12:06:00"),
			errors:    jwt.ValidationErrorExpired,
		},
		{
			secret:    "mysecret",
			accountID: "35581BF4-32C8-4908-8377-2E6A021D3D2B",
			token:     "",
			now:       createTime(t, "2020/03/22 12:06:00"),
			errors:    jwt.ValidationErrorMalformed,
		},
	}
	for _, params := range testParams {
		// Create a server
		srv := server{
			config: Config{JwtSecret: params.secret},
		}
		// Check if the token authorises the supplied account
		at(params.now, func() {
			err := srv.checkAuthorized(params.token, params.accountID)
			if err == nil && params.errors == 0 {
				// We weren't expecting an error
				return
			}
			vErr, ok := err.(*jwt.ValidationError)
			if !ok {
				t.Errorf("Unexpected error format: %v", err)
			} else if vErr.Errors&params.errors == 0 {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	// Create some test inputs
	testParams := []struct {
		authHeader http.Header
		token      string
		error      error
	}{
		{authHeader: http.Header{"Authorization": {"Bearer mytoken"}}, token: "mytoken", error: nil},
		{authHeader: http.Header{"Authorization": {"badHeader"}}, token: "", error: ErrMalformattedAuthHeader},
		{authHeader: http.Header{}, token: "", error: ErrNoAuthHeader},
	}
	// Run the test
	for _, params := range testParams {
		// Create a server
		srv := server{}
		// Get the token
		token, err := srv.getToken(&params.authHeader)
		// Assert results
		if err != params.error {
			t.Errorf("Unexpected error: %v", err)
		} else if token != params.token {
			t.Errorf("Unexpected token: %s", token)
		}
	}
}
