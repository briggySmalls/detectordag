package models

type DeviceRegisteredCertificate struct {

	Certificate string `json:"certificate"`

	PublicKey string `json:"publicKey"`

	PrivateKey string `json:"privateKey"`
}
