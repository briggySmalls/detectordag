package models

import (
	"time"
)

type Device struct {

	Name string `json:"name"`

	DeviceId string `json:"deviceId"`

	State *DeviceState `json:"state"`

	Updated time.Time `json:"updated"`
}
