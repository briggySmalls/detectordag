package models

import (
	"time"
)

// swagger:model device
type Device struct {
	// Name of the device
	Name string `json:"name"`
	// ID of the device
	DeviceId string `json:"deviceId"`
	// State of the power of the device
	State *DeviceState `json:"state"`
	// When the power status was last updated
	Updated time.Time `json:"updated"`
}

// swagger:parameters updateDevice registerDevice
type DeviceParameter struct {
	// ID of device
	//
	// required: true
	// in: path
	DeviceID string `json:"deviceId"`
}

// Successful devices retrieval
// swagger:response getDevicesResponse
type GetDevicesResponse struct {
	// in: body
	Body []Device
}

// Successful device retrieval
// swagger:response getDeviceResponse
type GetDeviceResponse struct {
	// in: body
	Body Device
}

// Device with that ID not found
// swagger:response deviceNotFoundResponse
type DeviceNotFoundResponse struct {
	// in: body
	Body ModelError
}
