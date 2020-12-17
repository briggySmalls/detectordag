package models

import (
	"time"
)

type Device struct {
	// Name of the device
	// required: true
	Name string `json:"name"`
	// ID of the device
	// required: true
	DeviceId string `json:"deviceId"`
	// State of the device
	// required: true
	State *DeviceState `json:"state"`
	// When the state was last updated
	// required: true
	Updated time.Time `json:"updated"`
}

type DeviceState struct {
	// Power status of the device
	// required: true
	Power string `json:"power"`
}

type MutableDevice struct {
	// The name of the device
	Name string `json:"name"`
}

type DeviceRegisteredCertificate struct {
	Certificate string `json:"certificate"`
	PublicKey   string `json:"publicKey"`
	PrivateKey  string `json:"privateKey"`
}

type DeviceRegistered struct {
	Name        string                       `json:"name"`
	DeviceId    string                       `json:"deviceId"`
	Certificate *DeviceRegisteredCertificate `json:"certificate"`
}

// swagger:parameters updateDevice registerDevice
type DeviceParameter struct {
	// ID of device
	//
	// required: true
	// in: path
	DeviceID string `json:"deviceId"`
}

// swagger:parameters updateDevice registerDevice
type MutableDeviceParameter struct {
	// Properties to update about the device
	//
	// required: true
	// in: body
	Device MutableDevice
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

// Device has been successfully registered
// swagger:response deviceRegisteredResponse
type DeviceRegisteredResponse struct {
	// in:body
	Body DeviceRegistered
}
