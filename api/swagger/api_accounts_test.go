package swagger

import (
	"encoding/json"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"

	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

func TestGetDevicesSuccess(t *testing.T) {
	// Define some test constants
	const (
		accountId = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		token     = ""
	)
	// Create a client
	c := createMockClient(t)
	// Create unit under test
	s := server{db: c}
	// Configure the mock db client to expect a call to fetch the account
	account := database.Account{AccountId: accountId, Emails: []string{"email@email@example.com"}}
	c.EXPECT().GetAccountById(gomock.Eq(accountId)).Return(&account, nil)
	// Create a request for devices
	req := createRequest(t, "GET", fmt.Sprintf("/v1/accounts/%s/devices", accountId), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// Execute the handler
	rr := runHandler(s.GetDevices, req)
	// Assert status ok
	assertStatus(t, rr, http.StatusOK)
	// Inspect the body of the response
	var resp []models.Device
	var err error
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Body could not be unmarshalled as device array: %v", rr.Body.String())
	}
}
