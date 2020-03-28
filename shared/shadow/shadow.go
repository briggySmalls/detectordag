package shadow

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface"
	"log"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Parse data to int
	epoch, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	log.Print(epoch)
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

// New creates a new shadow client
func New(sess *session.Session) Client {
	// We need to use an IoT control plane client to get an endpoint address
	ctrlSvc := iot.New(sess)
	descResp, err := ctrlSvc.DescribeEndpoint(&iot.DescribeEndpointInput{})
	if err != nil {
		log.Fatalf("failed to get dataplane endpoint: %v", err)
	}
	// Create a IoT data plane client using the endpoint address we retrieved
	svc := iotdataplane.New(sess, &aws.Config{
		Endpoint: descResp.EndpointAddress,
	})
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
