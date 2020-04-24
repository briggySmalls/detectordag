package swagger

import (
	"encoding/json"
	"fmt"
	models "github.com/briggysmalls/detectordag/api/swagger/go"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetDevicesSuccess(t *testing.T) {
	// Define some test constants
	const (
		accountID = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		token     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzU4MDcsImlzcyI6ImRldGVjdG9yZGFnIn0.CzyaCEIXlq1E0F89HR2Z9wbUn5gBDyQKTOCxTsX6iiQ"
	)
	// Create a client
	_, shdw, _, iotClient, tokens, router := createRealRouter(t)
	// Configure the tokens to expect a call to validate a token
	tokens.EXPECT().Validate(gomock.Eq(token)).Return(accountID, nil)
	// Configure the IoT client to expect a request for devices
	devices := []*iot.Device{
		{AccountId: "35581BF4-32C8-4908-8377-2E6A021D3D2B", DeviceId: "63eda5eb-7f56-417f-88ed-44a9eb9e5f67"},
		{AccountId: "35581BF4-32C8-4908-8377-2E6A021D3D2B", DeviceId: "4e9a7d26-d4de-4ea9-a0be-ec1b8264e35b"},
	}
	iotClient.EXPECT().GetThingsByAccount(gomock.Eq(accountID)).Return(devices, nil)
	// Configure the mock shadow client to expect calls for each device
	shadows := []shadow.Shadow{
		{
			Timestamp: shadow.Timestamp{createTime(t, "2020/03/22 00:27:00")},
			State:     shadow.State{Reported: map[string]interface{}{"status": true}},
			Metadata: shadow.Metadata{
				Reported: map[string]shadow.MetadataEntry{
					"status": {
						Timestamp: shadow.Timestamp{createTime(t, "2020/03/22 01:27:00")},
					},
				},
			},
		},
		{
			Timestamp: shadow.Timestamp{createTime(t, "2020/03/22 00:27:00")},
			State:     shadow.State{Reported: map[string]interface{}{"status": true}},
			Metadata: shadow.Metadata{
				Reported: map[string]shadow.MetadataEntry{
					"status": {
						Timestamp: shadow.Timestamp{createTime(t, "2020/03/22 01:27:00")},
					},
				},
			},
		},
	}
	shdw.EXPECT().Get(gomock.Any()).Return(&shadows[0], nil).Return(&shadows[1], nil).Times(2)
	// Create a request for devices
	req := createRequest(t, "GET", fmt.Sprintf("/v1/accounts/%s/devices", accountID), nil)
	// req = mux.SetURLVars(req, map[string]string{"accountId": accountID})
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// Execute the handler
	rr := runHandler(router, req)
	// Assert status ok
	assert.Equal(t, http.StatusOK, rr.Code)
	// Inspect the body of the response
	var resp []models.Device
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
}
