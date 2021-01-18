package shadow

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface"
)

const (
	TopicPatternRequestStatusUpdate = "dags/%s/status/request"
)

// Client represents a client to the device shadow service
type Client interface {
	Get(deviceId string) (*Shadow, error)
	UpdateConnectionStatus(deviceID string, status string, updated time.Time) (*Shadow, error)
	UpdateConnectionTransientID(deviceID string, ID string) error
	UpdateName(deviceId, name string) (*Shadow, error)
	RequestStatusUpdate(deviceID string) error
}

type client struct {
	dp iotdataplaneiface.IoTDataPlaneAPI
}

type ConnectionUpdatePayload struct {
	State struct {
		Reported struct {
			Connection struct {
				Status  string    `json:"current"`
				Updated Timestamp `json:"updated"`
			} `json:"connection"`
		} `json:"reported"`
	} `json:"state"`
}

func (p *ConnectionUpdatePayload) Dump() ([]byte, error) {
	return json.Marshal(p)
}

type TransientConnectionUpdatePayload struct {
	State struct {
		Reported struct {
			Connection struct {
				TransientID string `json:"transientId"`
			} `json:"connection"`
		} `json:"reported"`
	} `json:"state"`
}

type NameUpdatePayload struct {
	State struct {
		Reported struct {
			Name string `json:"name"`
		} `json:"reported"`
	} `json:"state"`
}

func (p *NameUpdatePayload) Dump() ([]byte, error) {
	return json.Marshal(p)
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
	if err != nil {
		return nil, err
	}
	// Unpack
	var shadowSchema DeviceShadowSchema
	return shadowSchema.Extract(payload)
}

func (c *client) UpdateConnectionStatus(deviceID, status string, updated time.Time) (*Shadow, error) {
	// Create new reported state
	updatePayload := ConnectionUpdatePayload{}
	updatePayload.State.Reported.Connection.Status = status
	updatePayload.State.Reported.Connection.Updated.Time = updated
	// Bundle up the request
	payload, err := updatePayload.Dump()
	if err != nil {
		return nil, err
	}
	// Update
	return c.updateShadow(deviceID, payload)
}

func (c *client) UpdateName(deviceID, name string) (*Shadow, error) {
	// Create new reported state
	updatePayload := NameUpdatePayload{}
	updatePayload.State.Reported.Name = name
	// Bundle up the request
	payload, err := updatePayload.Dump()
	if err != nil {
		return nil, err
	}
	// Make the request
	return c.updateShadow(deviceID, payload)
}

func (c *client) updateShadow(deviceID string, payload []byte) (*Shadow, error) {
	// Make the request
	log.Print(string(payload))
	_, err := c.dp.UpdateThingShadow(&iotdataplane.UpdateThingShadowInput{
		ThingName: aws.String(deviceID),
		Payload:   payload,
	})
	if err != nil {
		return nil, err
	}
	// Request the shadow
	shdw, err := c.getShadow(deviceID)
	if err != nil {
		return nil, err
	}
	// Parse the response
	var shadowSchema DeviceShadowSchema
	return shadowSchema.Extract([]byte(shdw))
}

func (c *client) UpdateConnectionTransientID(deviceID, ID string) error {
	// Create new reported state
	updatePayload := TransientConnectionUpdatePayload{}
	updatePayload.State.Reported.Connection.TransientID = ID
	// Bundle up the request
	payload, err := json.Marshal(updatePayload)
	if err != nil {
		return err
	}
	// Make the request
	log.Print(string(payload))
	_, err = c.dp.UpdateThingShadow(&iotdataplane.UpdateThingShadowInput{
		ThingName: aws.String(deviceID),
		Payload:   payload,
	})
	return err
}

func (c *client) RequestStatusUpdate(deviceID string) error {
	_, err := c.dp.Publish(&iotdataplane.PublishInput{
		Qos:     aws.Int64(1),
		Topic:   aws.String(fmt.Sprintf(TopicPatternRequestStatusUpdate, deviceID)),
		Payload: []byte("{}"),
	})
	return err
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
