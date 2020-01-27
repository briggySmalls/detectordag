package shared

type Message struct {
	Version string      `json:""`
	Payload interface{} `json:""`
}

type StatusMessageV1 struct {
	Version   string    `json:""`
	Status    boolean   `json:""`
	Timestamp time.Time `json:""`
	DeviceID  string    `json:""`
}
