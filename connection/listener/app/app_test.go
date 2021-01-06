package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_connection_updater.go -package app github.com/briggysmalls/detectordag/connection ConnectionUpdater
//go:generate go run github.com/golang/mock/mockgen -destination mock_sqs.go -package app -mock_names Client=MockSQSClient github.com/briggysmalls/detectordag/shared/sqs Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient github.com/briggysmalls/detectordag/shared/iot Client

import (
	"testing"
	"time"

	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/shared/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		app, _, _, _ := getStubbedApp(t)
		// Prepare an event
		event := DeviceLifecycleEvent{EventType: params.event}
		// Run the test
		assert.NotNil(t, app.RunJob(nil, event))
	}
}

func TestConnectedEvent(t *testing.T) {
	const (
		deviceID   = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		accountID  = "0ba69d11-28a2-433a-8403-8a269b94e61f"
		eventType  = "connected"
		status     = shadow.CONNECTION_STATUS_CONNECTED
		timestamp  = 0
		timeString = "1970/01/01 00:00:00"
	)
	// Create app under test
	app, mockConnectionUpdater, mockIoTClient, _ := getStubbedApp(t)
	// Configure call to fetching device
	device := iot.Device{
		DeviceId: deviceID, AccountId: accountID,
	}
	mockIoTClient.EXPECT().GetThing(deviceID).Return(&device, nil)
	// Configure call to update the connection status
	mockConnectionUpdater.EXPECT().UpdateConnectionStatus(&device, createTime(t, timeString), status)
	// Prepare an event
	event := DeviceLifecycleEvent{DeviceID: deviceID, EventType: eventType, Timestamp: timestamp}
	// Run the test
	assert.Nil(t, app.RunJob(nil, event))
}

func TestDisconnectedEvent(t *testing.T) {
	const (
		deviceID    = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		eventType   = "disconnected"
		status      = false
		timestamp   = 0
		timeString  = "1970/01/01 00:00:00"
		TransientID = "f5e0f6e2-e8b4-4233-bba2-8b2f0725d483"
	)
	// Create app under test
	app, _, _, mockSQS := getStubbedApp(t)
	// Configure call to queue the disconnected event
	mockSQS.EXPECT().QueueConnectionEvent(sqs.ConnectionEventPayload{
		DeviceID: deviceID,
		Status:   shadow.CONNECTION_STATUS_DISCONNECTED,
		Time:     createTime(t, timeString),
		ID:       TransientID,
	})
	// Prepare an event
	event := DeviceLifecycleEvent{DeviceID: deviceID, EventType: eventType, Timestamp: timestamp}
	// Run the test
	assert.Nil(t, app.RunJob(nil, event))
}

func getStubbedApp(t *testing.T) (*app, *MockConnectionUpdater, *MockIoTClient, *MockSQSClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create connection status updater
	connUpdater := NewMockConnectionUpdater(ctrl)
	// Create mock sqs
	sqs := NewMockSQSClient(ctrl)
	// Create mock iot client
	iot := NewMockIoTClient(ctrl)
	// Bundle up into an app
	return &app{updater: connUpdater, iot: iot, sqs: sqs}, connUpdater, iot, sqs
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
