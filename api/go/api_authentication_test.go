package swagger

import (
	"encoding/json"
	"github.com/briggysmalls/detectordag/api/mocks"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthSuccess(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database client
	c := mocks.NewMockClient(ctrl)
	// Create unit under test
	h := handlerer{db: c}
	// Configure the mock db client to expect a call to fetch the account
	const (
		username  = "email@example.com"
		accountId = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
	)
	account := database.Account{
		AccountId: accountId,
		Username:  username,
		Password:  "$2y$12$Nt3ajpggM4ViynWVGLOpW.JSbnVVVKRjNuw/ZYI71cj1WNG3Fty0K",
	}
	c.EXPECT().GetAccountByUsername(gomock.Eq(username)).Return(&account, nil)
	// Create a request to authenticate
	const body = `{"username": "email@example.com", "password": "mypassword"}`
	req, err := http.NewRequest("POST", "/v1/auth", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Auth)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	var token Token
	err = json.Unmarshal(rr.Body.Bytes(), &token)
	if err != nil {
		t.Fatalf("Body could not be unmarshalled as a token: %v", rr.Body.String())
	}
	if token.AccountId != accountId {
		t.Fatalf("handler returned unexpected account ID: got %s want %s", token.AccountId, accountId)
	}
}
