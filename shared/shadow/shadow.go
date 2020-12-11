package shadow

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface"
	"log"
	"strconv"
	"time"
)

// Client represents a client to the device shadow service
type Client interface {
	Get(deviceId string) (*Shadow, error)
	UpdateConnectionStatus(deviceID string, status bool) error
	GetConnectionStatus(deviceID string) (ConnectionState, error)
}

type Timestamp struct {
	time.Time
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

type client struct {
	dp iotdataplaneiface.IoTDataPlaneAPI
}

type MetadataEntry struct {
	Timestamp Timestamp `json:""`
}

type Shadow struct {
	Timestamp Timestamp `json:""`
	Metadata  struct {
		Reported map[string]MetadataEntry `json:""`
	} `json:""`
	State struct {
		Reported map[string]interface{} `json:"reported"`
	} `json:""`
}

type ConnectionUpdatePayload struct {
	State struct {
		Reported struct {
			ConnectionState bool
		} `json:"reported"`
	} `json:"state"`
}

// New creates a new shadow client
func New(sess *session.Session) (Client, error) {
	// We need to use an IoT control plane client to get an endpoint address
	ctrlSvc := iot.New(sess)
	descResp, err := ctrlSvc.DescribeEndpoint(&iot.DescribeEndpointInput{})
	if err != nil {
		return nil, err
	}
	// Create a IoT data plane client using the endpoint address we retrieved
	svc := iotdataplane.New(sess, &aws.Config{
		Endpoint: descResp.EndpointAddress,
	})
	// Return our client wrapper
	return &client{
		dp: svc,
	}, nil
}

func (c *client) Get(deviceId string) (*Shadow, error) {
	// Request the shadow
	payload, err := c.getShadow(deviceId)
	// Unpack
	var shadow Shadow
	err = json.Unmarshal(payload, &shadow)
	if err != nil {
		return nil, err
	}
	// Return
	return &shadow, nil
}

func (c *client) UpdateConnectionStatus(deviceID string, status bool) error {
	// Create new reported state
	newState := ConnectionUpdatePayload{State: {Reported: {ConnectionState: status}}}
	// Bundle up the request
	payload, err := json.Marshal(newState)
	if err != nil {
		return nil
	}
	// Form the request
	log.Print(string(payload))
	c.dp.UpdateThingShadow(&iotdataplane.UpdateThingShadowInput{
		ThingName: aws.String(deviceID),
		Payload:   payload,
	})
	return nil
}

func (c *client) GetConnectionStatus(deviceID string) (*ConnectionState, error) {
	// Request the shadow
	payload, err := c.getShadow(deviceID)
	if err != nil {
		return nil, err
	}
	// Unpack the payload
	var connState ConnectionStateSchema
	if err := connState.Load(payload); err != nil {
		return nil, err
	}
	// Repackage nicely
	flat := connState.Flatten()
	return &flat, nil
}

func (c *client) getShadow(deviceID string) ([]byte, error) {
	// Request the shadow
	resp, err := c.dp.GetThingShadow(&iotdataplane.GetThingShadowInput{
		ThingName: aws.String(deviceID),
	})
	// Bail on error
	if err != nil {
		return nil, fmt.Errorf("Get shadow failure for '%s': %w", deviceID, err)
	}
	// Just return the payload
	return resp.Payload, nil
}
