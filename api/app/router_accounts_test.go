package app

import (
	"encoding/json"
	"fmt"
	"github.com/briggysmalls/detectordag/api/app/models"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetDevicesSuccess(t *testing.T) {
	// Define some test constants
	const (
		accountID = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		token     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiIzNTU4MUJGNC0zMkM4LTQ5MDgtODM3Ny0yRTZBMDIxRDNEMkIiLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzU4MDcsImlzcyI6ImRldGVjdG9yZGFnIn0.CzyaCEIXlq1E0F89HR2Z9wbUn5gBDyQKTOCxTsX6iiQ"
	)
	deviceInputs := []struct {
		name    string
		ID      string
		state   string
		updated time.Time
	}{
		{name: "one", ID: "63eda5eb-7f56-417f-88ed-44a9eb9e5f67", state: "on", updated: createTime(t, "2020/03/22 01:27:00")},
		{name: "two", ID: "4e9a7d26-d4de-4ea9-a0be-ec1b8264e35b", state: "off", updated: createTime(t, "2020/03/22 01:20:00")},
	}
	// Create a client
	_, shdw, _, iotClient, tokens, router := createRealRouter(t)
	// Configure the tokens to expect a call to validate a token
	tokens.EXPECT().Validate(token).Return(accountID, nil)
	// Configure the IoT client to expect a request for devices
	devices := []*iot.Device{
		{Name: deviceInputs[0].name, AccountId: accountID, DeviceId: deviceInputs[0].ID},
		{Name: deviceInputs[1].name, AccountId: accountID, DeviceId: deviceInputs[1].ID},
	}
	iotClient.EXPECT().GetThingsByAccount(accountID).Return(devices, nil)
	// Configure the mock shadow client to expect calls for each device
	shdw.EXPECT().Get(deviceInputs[0].ID).Return(&shadow.Shadow{
		Time: createTime(t, "2020/03/22 00:27:00"),
		Power: shadow.StringShadowField{
			Value:   deviceInputs[0].state,
			Updated: deviceInputs[0].updated,
		},
	}, nil)
	shdw.EXPECT().Get(deviceInputs[1].ID).Return(&shadow.Shadow{
		Time: createTime(t, "2020/03/22 00:27:00"),
		Power: shadow.StringShadowField{
			Value:   deviceInputs[1].state,
			Updated: deviceInputs[1].updated,
		},
	}, nil)
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
	assert.Equal(t, []models.Device{
		{
			Name:     deviceInputs[0].name,
			DeviceId: deviceInputs[0].ID,
			State:    &models.DeviceState{Power: deviceInputs[0].state},
			Updated:  deviceInputs[0].updated,
		},
		{
			Name:     deviceInputs[1].name,
			DeviceId: deviceInputs[1].ID,
			State:    &models.DeviceState{Power: deviceInputs[1].state},
			Updated:  deviceInputs[1].updated,
		},
	}, resp)
}

func TestRegisterDevice(t *testing.T) {
	// Define some test constants
	const (
		accountID   = "35581BF4-32C8-4908-8377-2E6A021D3D2B"
		deviceID    = "63eda5eb-7f56-417f-88ed-44a9eb9e5f67"
		token       = "my-crazy-token"
		desiredName = "device-name"
		publicCert  = "impublic"
		privateCert = "imprivate"
		cert        = "imcert"
	)
	// Create a client
	_, _, _, iotClient, tokens, router := createRealRouter(t)
	// Configure the tokens to expect a call to validate a token
	tokens.EXPECT().Validate(token).Return(accountID, nil)
	// Configure the IoT client to expect a request to register a new device
	device := iot.Device{Name: desiredName, AccountId: accountID, DeviceId: deviceID}
	certs := iot.Certificates{Public: publicCert, Private: privateCert, Certificate: cert}
	iotClient.EXPECT().RegisterThing(accountID, deviceID, desiredName).Return(&device, &certs, nil)
	// Create a request for devices
	req := createRequest(t, "PUT",
		fmt.Sprintf("/v1/accounts/%s/devices/%s", accountID, deviceID),
		[]byte(fmt.Sprintf(`{"name": "%s"}`, desiredName)),
	)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// Execute the handler
	rr := runHandler(router, req)
	// Assert status ok
	assert.Equal(t, http.StatusOK, rr.Code)
	// Inspect the body of the response
	const expectedBody = `{"name":"%s","deviceId":"%s","certificate":{"certificate":"%s","publicKey":"%s","privateKey":"%s"}}`
	assert.Equal(t, fmt.Sprintf(expectedBody, desiredName, deviceID, cert, publicCert, privateCert), string(rr.Body.Bytes()))
}
