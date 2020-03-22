package shadow

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface"
	"time"
)

type client struct {
	dp iotdataplaneiface.IoTDataPlaneAPI
}

type MetadataEntry struct {
	Timestamp time.Time
}

type Metadata struct {
	Reported map[string]MetadataEntry `json:""`
}

type State struct {
	Reported map[string]interface{} `json:""`
}

type Shadow struct {
	Timestamp time.Time `json:""`
	Metadata  Metadata  `json:""`
	State     State     `json:""`
}

// New creates a new shadow client
func New(sess *session.Session) Client {
	// Create the service client
	svc := iotdataplane.New(sess)
	// Return our client wrapper
	return &client{
		dp: svc,
	}
}

// Client represents a client to the device shadow service
type Client interface {
	Get(deviceId string) (*Shadow, error)
}

func (c *client) Get(deviceId string) (*Shadow, error) {
	// Request the shadow
	resp, err := c.dp.GetThingShadow(&iotdataplane.GetThingShadowInput{
		ThingName: aws.String(deviceId),
	})
	if err != nil {
		return nil, err
	}
	// Unpack
	var shadow Shadow
	err = json.Unmarshal(resp.Payload, &shadow)
	if err != nil {
		return nil, err
	}
	// Return
	return &shadow, nil
}
