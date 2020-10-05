package models

type DeviceRegistered struct {

	Name string `json:"name"`

	DeviceId string `json:"deviceId"`

	Certificate *DeviceRegisteredCertificate `json:"certificate"`
}
