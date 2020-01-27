package shared

import (
	"time"
)

type Message struct {
	Version string      `json:""`
	Payload interface{} `json:""`
}

type PowerStatusChangedV1 struct {
	Version   string    `json:""`
	Status    bool      `json:""`
	Timestamp time.Time `json:""`
	DeviceID  string    `json:""`
}
