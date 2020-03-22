package swagger

import (
	"encoding/json"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

func TestGetDevicesSuccess(t *testing.T) {
	// Define some test constants
	const (
		accountId = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		token     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzU4MDcsImlzcyI6ImRldGVjdG9yZGFnIn0.CzyaCEIXlq1E0F89HR2Z9wbUn5gBDyQKTOCxTsX6iiQ"
	)
	// Create a client
	db, shadow := createMocks(t)
	// Create unit under test
	s := server{db: db, config: Config{JwtSecret: jwtSecret}, shadow: shadow}
	// Configure the mock db client to expect a call to query for devices in an account
	devices := []database.Device{
		{AccountId: "35581BF4-32C8-4908-8377-2E6A021D3D2B", DeviceId: "63eda5eb-7f56-417f-88ed-44a9eb9e5f67"},
		{AccountId: "35581BF4-32C8-4908-8377-2E6A021D3D2B", DeviceId: "4e9a7d26-d4de-4ea9-a0be-ec1b8264e35b"},
	}
	db.EXPECT().GetDevicesByAccount(gomock.Eq(accountId)).Return(devices, nil)
	// Create a request for devices
	req := createRequest(t, "GET", fmt.Sprintf("/v1/accounts/%s/devices", accountId), nil)
	req = mux.SetURLVars(req, map[string]string{
		"accountId": accountId,
	})
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// Execute the handler
	rr := runHandler(s.GetDevices, req)
	// Assert status ok
	assertStatus(t, rr, http.StatusOK)
	// Inspect the body of the response
	var resp []models.Device
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Body could not be unmarshalled as device array: %v", rr.Body.String())
	}
}
