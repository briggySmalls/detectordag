package shadow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/briggysmalls/detectordag/shared"
)

const (
	CONNECTION_STATUS_CONNECTED    = "connected"
	CONNECTION_STATUS_DISCONNECTED = "disconnected"
	POWER_STATUS_ON                = "on"
	POWER_STATUS_OFF               = "off"
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

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	// Convert the timestamp to a number
	return []byte(fmt.Sprint(t.Unix())), nil
}

type PowerShadow struct {
	Value   string
	Updated time.Time
}

type ConnectionShadow struct {
	Status      string
	Updated     time.Time
	TransientID string
}

type Shadow struct {
	Time       time.Time
	Version    int
	Name       string
	Connection ConnectionShadow
	Power      PowerShadow
}

type DeviceShadowSchema struct {
	Timestamp Timestamp
	Version   int
	State     struct {
		Reported struct {
			Name       string
			Connection struct {
				Current     string    `validate:"required,eq=connected|eq=disconnected"`
				Updated     Timestamp `validate:"required"`
				TransientID string    `validate:"required,uuid"`
			}
			Status string `validate:"required,eq=on|eq=off"`
		}
	}
	Metadata struct {
		Reported struct {
			Status MetadataEntry `validate:"required"`
		}
	}
}

// Extract converts the information into a more user-friendly form
func (c *DeviceShadowSchema) Extract(payload []byte) (*Shadow, error) {
	// Load the json into this struct
	if err := json.Unmarshal(payload, c); err != nil {
		return nil, err
	}
	// Validate the struct
	if err := shared.Validate.Struct(c); err != nil {
		return nil, err
	}
	// Create a shadow
	s := Shadow{
		Time:    c.Timestamp.Time,
		Version: c.Version,
		Name:    c.State.Reported.Name,
		Connection: ConnectionShadow{
			Status:      c.State.Reported.Connection.Current,
			Updated:     c.State.Reported.Connection.Updated.Time,
			TransientID: c.State.Reported.Connection.TransientID,
		},
		Power: PowerShadow{
			Value:   c.State.Reported.Status,
			Updated: c.Metadata.Reported.Status.Timestamp.Time,
		},
	}
	// Extract the fields we care about
	return &s, nil
}
