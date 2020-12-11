package shadow

import (
	"encoding/json"
)

type ConnectionState struct {
	State     bool      `json:""`
	Timestamp Timestamp `json:""`
}

type ConnectionStateSchema struct {
	State struct {
		Reported struct {
			Connection bool `json:""`
		} `json:""`
	} `json:""`
	Metadata struct {
		Reported struct {
			Connection MetadataEntry
		}
	}
}

// Flatten converts the information into a more user-friendly form
func (c *ConnectionStateSchema) Flatten() ConnectionState {
	return ConnectionState{
		State:     c.State.Reported.Connection,
		Timestamp: c.Metadata.Reported.Connection.Timestamp,
	}
}

// Load populates a connection state struct with values
func (c *ConnectionStateSchema) Load(payload []byte) error {
	return json.Unmarshal(payload, c)
}
