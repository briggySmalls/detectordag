package shadow

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface"
	"strconv"
	"time"
)

// Client represents a client to the device shadow service
type Client interface {
	Get(deviceId string) (*Shadow, error)
}

type Timestamp struct {
	time.Time
}

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

type Metadata struct {
	Reported map[string]MetadataEntry `json:""`
}

type State struct {
	Reported map[string]interface{} `json:""`
}

type Shadow struct {
	Timestamp Timestamp `json:""`
	Metadata  Metadata  `json:""`
	State     State     `json:""`
}

type ConnectionStatusUpdate struct {
	Status bool
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
	resp, err := c.dp.GetThingShadow(&iotdataplane.GetThingShadowInput{
		ThingName: aws.String(deviceId),
	})
	if err != nil {
		return nil, fmt.Errorf("Get shadow failure for '%s': %w", deviceId, err)
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

func (c *client) UpdateConnectionStatus(deviceId string, update ConnectionStatusUpdate) error {
	return nil
}
