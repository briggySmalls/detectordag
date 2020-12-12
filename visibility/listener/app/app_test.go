package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package app -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_sqs.go -package app -mock_names Client=MockSQSClient github.com/briggysmalls/detectordag/shared/sqs Client

import (
	"github.com/briggysmalls/detectordag/shared/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInvalidEventType(t *testing.T) {
	testParams := []struct {
		event string
	}{
		{event: "other"},
	}
	// Run the test
	for _, params := range testParams {
		// Create app under test
		app, _, _ := getStubbedApp(t)
		// Prepare an event
		event := DeviceLifecycleEvent{EventType: params.event}
		// Run the test
		assert.NotNil(t, app.RunJob(nil, event))
	}
}

func TestConnectedEvent(t *testing.T) {
	const (
		deviceID  = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		eventType = "connected"
		status    = true
		timestamp = 0
	)
	// Create app under test
	app, mockShadow, _ := getStubbedApp(t)
	// Configure call to update the connection status
	mockShadow.EXPECT().UpdateConnectionStatus(deviceID, status)
	// Prepare an event
	event := DeviceLifecycleEvent{DeviceID: deviceID, EventType: eventType, Timestamp: timestamp}
	// Run the test
	assert.Nil(t, app.RunJob(nil, event))
}

func TestDisconnectedEvent(t *testing.T) {
	const (
		deviceID   = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		eventType  = "disconnected"
		status     = false
		timestamp  = 0
		timeString = "1970/01/01 00:00:00"
	)
	// Create app under test
	app, _, mockSQS := getStubbedApp(t)
	// Configure call to queue the disconnected event
	mockSQS.EXPECT().QueueDisconnectedEvent(sqs.DisconnectedPayload{
		DeviceID: deviceID,
		Time:     createTime(t, timeString),
	})
	// Prepare an event
	event := DeviceLifecycleEvent{DeviceID: deviceID, EventType: eventType, Timestamp: timestamp}
	// Run the test
	assert.Nil(t, app.RunJob(nil, event))
}

func getStubbedApp(t *testing.T) (*app, *MockShadowClient, *MockSQSClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock shadow
	shadow := NewMockShadowClient(ctrl)
	// Create mock sqs
	sqs := NewMockSQSClient(ctrl)
	// Bundle up into an app
	return &app{shadow: shadow, sqs: sqs}, shadow, sqs
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
