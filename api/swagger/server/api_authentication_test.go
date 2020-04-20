package server

import (
	"encoding/json"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

const (
	jwtSecret = "mysecret"
)

func TestAuthSuccess(t *testing.T) {
	// Create some test constants
	const (
		username       = "email@example.com"
		accountId      = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		jwtDuration    = "2h"
		password       = "mypassword"
		hashedPassword = "$2y$12$Nt3ajpggM4ViynWVGLOpW.JSbnVVVKRjNuw/ZYI71cj1WNG3Fty0K"
	)
	// Create a mock client
	db, shadow := createMocks(t)
	// Create unit under test
	s := server{
		db:     db,
		config: Config{JwtSecret: jwtSecret, JwtDuration: jwtDuration},
		shadow: shadow,
	}
	// Configure the mock db client to expect a call to fetch the account
	account := database.Account{AccountId: accountId, Username: username, Password: hashedPassword}
	db.EXPECT().GetAccountByUsername(gomock.Eq(username)).Return(&account, nil)
	// Create a request to authenticate
	req := createRequest(t, "POST", "/v1/auth", []byte(fmt.Sprintf(`{"username": "email@example.com", "password": "%s"}`, password)))
	// Execute the handler
	rr := runHandler(s.Auth, req)
	// Assert the HTTP status
	assertStatus(t, rr, http.StatusOK)
	// Check the response body is what we expect.
	var resp models.Token
	var err error
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Body could not be unmarshalled as a token: %v", rr.Body.String())
	}
	if resp.AccountId != accountId {
		t.Fatalf("handler returned unexpected account ID: got %s want %s", resp.AccountId, accountId)
	}
	// Parse the token contents
	token, err := jwt.ParseWithClaims(resp.Token, &CustomAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	// Check the claims
	if claims, ok := token.Claims.(*CustomAuthClaims); ok && token.Valid {
		if claims.AccountId != resp.AccountId {
			t.Fatalf("Token did not correspond with authenticated user")
		}
	} else {
		t.Fatalf(err.Error())
	}
}
