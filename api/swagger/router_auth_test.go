package swagger

import (
	"encoding/json"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		accountID      = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		jwtDuration    = "2h"
		password       = "mypassword"
		hashedPassword = "$2y$12$Nt3ajpggM4ViynWVGLOpW.JSbnVVVKRjNuw/ZYI71cj1WNG3Fty0K"
	)
	// Create a mock client
	db, _, _, _, tokens, router := createRealRouter(t)
	// Configure the mock db client to expect a call to fetch the account
	account := database.Account{AccountId: accountID, Username: username, Password: hashedPassword}
	db.EXPECT().GetAccountByUsername(gomock.Eq(username)).Return(&account, nil)
	// Configure the mock tokens to create a token
	expectedToken := "dummy-token"
	tokens.EXPECT().Create(gomock.Eq(accountID)).Return(expectedToken, nil)
	// Create a request to authenticate
	req := createRequest(t, "POST", "/v1/auth", []byte(fmt.Sprintf(`{"username": "email@example.com", "password": "%s"}`, password)))
	// Execute the handler
	rr := runHandler(router, req)
	// Assert the HTTP status
	assert.Equal(t, http.StatusOK, rr.Code)
	// Check the response body is what we expect.
	var resp models.Token
	var err error
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, accountID, resp.AccountId)
	// Parse the token contents
	assert.Equal(t, expectedToken, resp.Token)
}
