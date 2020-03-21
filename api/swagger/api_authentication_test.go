package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/briggysmalls/detectordag/api/mocks"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
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
	c := createMockClient(t)
	// Create unit under test
	s := server{
		db:     c,
		config: Config{JwtSecret: jwtSecret, JwtDuration: jwtDuration},
	}
	// Configure the mock db client to expect a call to fetch the account
	account := database.Account{AccountId: accountId, Username: username, Password: hashedPassword}
	c.EXPECT().GetAccountByUsername(gomock.Eq(username)).Return(&account, nil)
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

func createMockClient(t *testing.T) *mocks.MockClient {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database client
	return mocks.NewMockClient(ctrl)
}

func runHandler(h func(http.ResponseWriter, *http.Request), req *http.Request) *httptest.ResponseRecorder {
	// Run the handler using test code
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)
	handler.ServeHTTP(rr, req)
	return rr
}

func createRequest(t *testing.T, method, route string, body []byte) *http.Request {
	req, err := http.NewRequest(method, route, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expected)
	}
}
