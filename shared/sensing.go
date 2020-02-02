package shared

import (
	"encoding/json"
	"github.com/denisbrodbeck/machineid"
	"github.com/streadway/amqp"
	"time"
)

var powerStatusConfig = queueConfig{
	Name: "power_status",
}

type SensingMessenger interface {
	clientInterface
	PowerStatusChanged(status bool) error
}

type SensingReceiver interface {
	clientInterface
	PowerStatusConsumer() (<-chan amqp.Delivery, error)
}

func NewSensingMessenger() SensingMessenger {
	return newClient([]queueConfig{powerStatusConfig})
}

func NewSensingReceiver() SensingReceiver {
	return newClient([]queueConfig{powerStatusConfig})
}

// PowerStatusChanged sends a message saying the power status has changed
func (c *client) PowerStatusChanged(status bool) error {
	// Create the message data
	data, err := newStatusMessage(status)
	if err != nil {
		return WrapError(err, "Failed to create message")
	}
	// Serialise the message
	b, err := json.Marshal(data)
	if err != nil {
		return WrapError(err, "Failed to serialise message")
	}
	// Send the message
	return c.send(powerStatusConfig.Name, b)
}

func (c *client) PowerStatusConsumer() (<-chan amqp.Delivery, error) {
	return c.getConsumer(powerStatusConfig.Name)
}

func newStatusMessage(status bool) (*PowerStatusChangedV1, error) {
	// Get machine ID
	id, err := machineid.ID()
	if err != nil {
		return nil, WrapError(err, "Failed to get machine ID")
	}
	// Create the message data
	return &PowerStatusChangedV1{
		Version:   "0.1.0",
		Status:    status,
		Timestamp: time.Now(),
		DeviceID:  id,
	}, nil
}
