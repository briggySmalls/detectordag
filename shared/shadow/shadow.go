package shadow

import (
	"encoding/json"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

type MetadataEntry struct {
	Timestamp Timestamp `json:""`
}

// UnmarshalJSON defines a custom method of deserialising a timestamp
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Parse data to int
	epoch, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	// Conver to a time
	t.Time = time.Unix(int64(epoch), 0)
	return nil
}

type Shadow struct {
	Time       time.Time
	Version    int
	Connection struct {
		Value   bool
		Updated time.Time
	}
	Power struct {
		Value   bool
		Updated time.Time
	}
}

type DeviceShadowSchema struct {
	Timestamp Timestamp
	Version   int
	State     struct {
		Reported struct {
			Connection bool `json:""`
			Status     bool
		} `json:""`
	} `json:""`
	Metadata struct {
		Reported struct {
			Connection MetadataEntry
			Status     MetadataEntry
		}
	}
}

// Extract converts the information into a more user-friendly form
func (c *DeviceShadowSchema) Extract(payload []byte) (*Shadow, error) {
	// Load the json into this struct
	if err := json.Unmarshal(payload, c); err != nil {
		return nil, err
	}
	// Create a shadow
	s := Shadow{}
	s.Time = c.Timestamp.Time
	s.Version = c.Version
	s.Connection.Value = c.State.Reported.Connection
	s.Connection.Updated = c.Metadata.Reported.Connection.Timestamp.Time
	s.Power.Value = c.State.Reported.Status
	s.Power.Updated = c.Metadata.Reported.Status.Timestamp.Time
	// Extract the fields we care about
	return &s, nil
}
