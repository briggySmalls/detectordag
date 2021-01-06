package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/briggysmalls/detectordag/api/app/models"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetDevicesSuccess(t *testing.T) {
	// Define some test constants
	const (
		accountID = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		token     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzU4MDcsImlzcyI6ImRldGVjdG9yZGFnIn0.CzyaCEIXlq1E0F89HR2Z9wbUn5gBDyQKTOCxTsX6iiQ"
	)
	// Create the results we expect
	devices := []models.Device{
		{
			Name:     "one",
			DeviceId: "63eda5eb-7f56-417f-88ed-44a9eb9e5f67",
			State: &models.DeviceState{
				Power:   "on",
				Updated: createTime(t, "2020/03/22 01:27:00"),
			},
			Connection: &models.DeviceConnection{
				Status:  "connected",
				Updated: createTime(t, "2020/03/22 01:27:01"),
			},
		},
		{
			Name:     "two",
			DeviceId: "4e9a7d26-d4de-4ea9-a0be-ec1b8264e35b",
			State: &models.DeviceState{
				Power:   "off",
				Updated: createTime(t, "2020/03/22 01:20:00"),
			},
			Connection: &models.DeviceConnection{
				Status:  "connected",
				Updated: createTime(t, "2020/03/22 01:20:01"),
			},
		},
	}
	// Create a client
	_, shdw, _, iotClient, tokens, router := createRealRouter(t)
	gomock.InOrder(
		// Configure the tokens to expect a call to validate a token
		tokens.EXPECT().Validate(token).Return(accountID, nil),
		// Configure the IoT client to expect a request for devices
		iotClient.EXPECT().GetThingsByAccount(accountID).Return([]*iot.Device{
			{AccountId: accountID, DeviceId: devices[0].DeviceId},
			{AccountId: accountID, DeviceId: devices[1].DeviceId},
		}, nil),
		// Configure the mock shadow client to expect calls for each device
		shdw.EXPECT().Get(devices[0].DeviceId).Return(&shadow.Shadow{
			Name: devices[0].Name,
			Time: createTime(t, "2020/03/22 00:27:00"),
			Power: shadow.PowerShadow{
				Value:   devices[0].State.Power,
				Updated: devices[0].State.Updated,
			},
			Connection: shadow.ConnectionShadow{
				Status:      devices[0].Connection.Status,
				TransientID: "0a72fb5e-6489-49d5-aadb-782d54d1ed1f",
				Updated:     devices[0].Connection.Updated,
			},
		}, nil),
		shdw.EXPECT().Get(devices[1].DeviceId).Return(&shadow.Shadow{
			Name: devices[1].Name,
			Time: createTime(t, "2020/03/22 00:27:00"),
			Power: shadow.PowerShadow{
				Value:   devices[1].State.Power,
				Updated: devices[1].State.Updated,
			},
			Connection: shadow.ConnectionShadow{
				Status:      devices[1].Connection.Status,
				TransientID: "63c7d830-724a-4715-aa1e-dc8934cd32fb",
				Updated:     devices[1].Connection.Updated,
			},
		}, nil),
	)
	// Create a request for devices
	req := createRequest(t, "GET", fmt.Sprintf("/v1/accounts/%s/devices", accountID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// Execute the handler
	rr := runHandler(router, req)
	// Assert status ok
	assert.Equal(t, http.StatusOK, rr.Code)
	// Inspect the body of the response
	var resp []models.Device
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, devices, resp)
}
