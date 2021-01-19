package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_sqs.go -package app -mock_names Client=MockSQSClient github.com/briggysmalls/detectordag/shared/sqs Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package app -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_connection.go -package app github.com/briggysmalls/detectordag/connection ConnectionUpdater

import (
	"testing"
	"time"

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

func TestDisconnectedEvents(t *testing.T) {
	const (
		deviceID   = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		timestamp  = 0
		timeString = "1970/01/01 00:00:00"
	)
	testParams := []struct {
		currentStatus string
		newStatus     string
	}{
		{ // Transition from connected to disconnected
			currentStatus: shadow.CONNECTION_STATUS_CONNECTED,
			newStatus:     shadow.CONNECTION_STATUS_DISCONNECTED,
		},
	}
	for _, params := range testParams {
		// Create app under test
		app, _, mockShadowClient, mockSQSClient := getStubbedApp(t)
		// Prepare a shadow to return
		gomock.InOrder(
			// We expect to always update the transient status
			mockShadowClient.EXPECT().UpdateConnectionTransientID(deviceID, gomock.Any()),
			// We also always get the shadow
			mockShadowClient.EXPECT().Get(deviceID).Return(&shadow.Shadow{Connection: shadow.ConnectionShadow{
				Status: params.currentStatus,
			}}, nil),
			// Expect a
			mockShadowClient.EXPECT().RequestStatusUpdate(deviceID),
			// This test checks that events are enqueued
			mockSQSClient.EXPECT().QueueConnectionEvent(gomock.Any()).Do(func(payload sqs.ConnectionEventPayload) {
				assert.Equal(t, deviceID, payload.DeviceID)
				assert.Equal(t, params.newStatus, payload.Status)
				assert.Equal(t, createTime(t, timeString), payload.Time)
			}),
		)
		// Prepare an event
		event := DeviceLifecycleEvent{
			DeviceID:  deviceID,
			EventType: params.newStatus,
			Timestamp: timestamp,
		}
		// Run the test
		assert.Nil(t, app.RunJob(nil, event))
	}
}

func TestConnectedEvents(t *testing.T) {
	const (
		deviceID   = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		timestamp  = 0
		accountID  = "c6d62b30-00ac-49c4-9268-88559a46889f"
		timeString = "1970/01/01 00:00:00"
	)
	testParams := []struct {
		currentStatus string
		newStatus     string
	}{
		{ // Transition from connected to disconnected
			currentStatus: shadow.CONNECTION_STATUS_DISCONNECTED,
			newStatus:     shadow.CONNECTION_STATUS_CONNECTED,
		},
	}
	for _, params := range testParams {
		// Create app under test
		app, mockUpdater, mockShadowClient, _ := getStubbedApp(t)
		// Prepare a shadow to return
		gomock.InOrder(
			// We expect to always update the transient status
			mockShadowClient.EXPECT().UpdateConnectionTransientID(deviceID, gomock.Any()),
			// We also always get the shadow
			mockShadowClient.EXPECT().Get(deviceID).Return(&shadow.Shadow{Connection: shadow.ConnectionShadow{
				Status: params.currentStatus,
			}}, nil),
			// Expect a call to update status
			mockUpdater.EXPECT().UpdateConnectionStatus(deviceID, createTime(t, timeString), shadow.CONNECTION_STATUS_CONNECTED),
		)
		// Prepare an event
		event := DeviceLifecycleEvent{
			DeviceID:  deviceID,
			EventType: params.newStatus,
			Timestamp: timestamp,
		}
		// Run the test
		assert.Nil(t, app.RunJob(nil, event))
	}
}

func TestIgnoreConnectionEvents(t *testing.T) {
	const (
		deviceID   = "792ac520-0733-4ffe-8137-8aba3ca446d7"
		timestamp  = 0
		timeString = "1970/01/01 00:00:00"
	)
	testParams := []struct {
		currentStatus string
		newStatus     string
	}{
		{ // Transition from connected to connected
			currentStatus: shadow.CONNECTION_STATUS_CONNECTED,
			newStatus:     shadow.CONNECTION_STATUS_CONNECTED,
		},
		{ // Transition from disconnected to disconnected
			currentStatus: shadow.CONNECTION_STATUS_DISCONNECTED,
			newStatus:     shadow.CONNECTION_STATUS_DISCONNECTED,
		},
	}
	for _, params := range testParams {
		// Create app under test
		app, _, mockShadowClient, _ := getStubbedApp(t)
		// Configure call to fetching device
		device := shadow.Shadow{Connection: shadow.ConnectionShadow{
			Status: params.currentStatus,
		}}
		gomock.InOrder(
			// We expect to always update the transient status
			mockShadowClient.EXPECT().UpdateConnectionTransientID(deviceID, gomock.Any()),
			// We also always get the shadow
			mockShadowClient.EXPECT().Get(deviceID).Return(&device, nil),
			// Nothing is enqueued
		)
		// Prepare an event
		event := DeviceLifecycleEvent{
			DeviceID:  deviceID,
			EventType: params.newStatus,
			Timestamp: timestamp,
		}
		// Run the test
		assert.Nil(t, app.RunJob(nil, event))
	}
}

func getStubbedApp(t *testing.T) (*app, *MockConnectionUpdater, *MockShadowClient, *MockSQSClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock updater client
	updater := NewMockConnectionUpdater(ctrl)
	// Create mock sqs
	sqs := NewMockSQSClient(ctrl)
	// Create mock shadow client
	shadow := NewMockShadowClient(ctrl)
	// Bundle up into an app
	return &app{updater: updater, shadow: shadow, sqs: sqs}, updater, shadow, sqs
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
