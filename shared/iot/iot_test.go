package iot

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetThing(t *testing.T) {
	// Create a session
	sesh := createSession(t)
	// Create a client under test
	c, err := New(sesh)
	assert.NoError(t, err)
	// Query for a known device
	device, err := c.GetThing("92f59eeb298c4f8c8773e4704d9afe65")
	assert.NoError(t, err)
	// Assert it has expected fields
	assert.Equal(t, "aac45d02-c97d-442c-8431-336d578fdcf7", device.AccountId)
}

func TestGetThingsByAccount(t *testing.T) {
	// Create a session
	sesh := createSession(t)
	// Create a client under test
	c, err := New(sesh)
	assert.NoError(t, err)
	// Query for devices associated with an account
	devices, err := c.GetThingsByAccount("aac45d02-c97d-442c-8431-336d578fdcf7")
	assert.NoError(t, err)
	// Assert it has expected fields
	assert.Len(t, devices, 1)
	assert.Equal(t, "aac45d02-c97d-442c-8431-336d578fdcf7", devices[0].AccountId)
}

func createSession(t *testing.T) *session.Session {
	// Create a session
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	assert.NoError(t, err)
	return sesh
}
