package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/briggysmalls/detectordag/api/mocks"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	jwtSecret = "mysecret"
)

func TestAuthSuccess(t *testing.T) {
	// Create some test constants
	const (
		username    = "email@example.com"
		accountId   = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		jwtDuration = "2h"
	)
	// Run the test
	rr := runTest(t, username, accountId, jwtDuration, "mypassword", "$2y$12$Nt3ajpggM4ViynWVGLOpW.JSbnVVVKRjNuw/ZYI71cj1WNG3Fty0K")
	// Assert the HTTP status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	var resp Token
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

func runTest(t *testing.T, username, accountId, jwtDuration, password, hashedPassword string) *httptest.ResponseRecorder {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database client
	c := mocks.NewMockClient(ctrl)
	// Create unit under test
	s := server{
		db:     c,
		config: Config{JwtSecret: jwtSecret, JwtDuration: jwtDuration},
	}
	// Configure the mock db client to expect a call to fetch the account
	account := database.Account{
		AccountId: accountId,
		Username:  username,
		Password:  hashedPassword,
	}
	c.EXPECT().GetAccountByUsername(gomock.Eq(username)).Return(&account, nil)
	// Create a request to authenticate
	body := fmt.Sprintf(`{"username": "email@example.com", "password": "%s"}`, password)
	req, err := http.NewRequest("POST", "/v1/auth", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// Run the handler using test code
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Auth)
	handler.ServeHTTP(rr, req)
	return rr
}
