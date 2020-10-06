package models

type DeviceRegisteredCertificate struct {

	Certificate string `json:"certificate"`

	PublicKey string `json:"publicKey"`

	PrivateKey string `json:"privateKey"`
}

// swagger:model deviceRegistered
type DeviceRegistered struct {

	Name string `json:"name"`

	DeviceId string `json:"deviceId"`

	Certificate *DeviceRegisteredCertificate `json:"certificate"`
}
