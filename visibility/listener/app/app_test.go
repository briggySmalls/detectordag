package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_sqs.go -package app -mock_names Client=MockSQSClient github.com/briggysmalls/detectordag/shared/sqs Client

import (
	"errors"
	"github.com/briggysmalls/detectordag/shared/iot"
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

func TestDeviceLookupFailed(t *testing.T) {
	// Create app under test
	app, mockIoT, _ := getStubbedApp(t)
	// Prepare some test parameters
	const (
		deviceID = "b6d62b30-00ac-49c4-9268-88559a46889f"
	)
	// Prepare an event
	event := DeviceLifecycleEvent{
		DeviceID:  deviceID,
		Timestamp: 0,
		EventType: "connected",
	}
	// Configure lookup to fail
	mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(nil, errors.New("Something went wrong"))
	// Run test
	assert.NotNil(t, app.RunJob(nil, event))
}

func TestMessageSent(t *testing.T) {
	testParams := []struct {
		event     string
		status    bool
		timestamp int64
		time      time.Time
	}{
		{event: "connected", status: true, timestamp: 0, time: createTime(t, "1970/01/01 00:00:00")},
		{event: "disconnected", status: false, timestamp: 0, time: createTime(t, "1970/01/01 00:00:00")},
	}
	const (
		deviceID = "792ac520-0733-4ffe-8137-8aba3ca446d7"
	)
	// Run the test
	for _, params := range testParams {
		// Create app under test
		app, mockIoT, mockSQS := getStubbedApp(t)
		// Configure lookup to pass
		device := iot.Device{
			Name:      "1",
			DeviceId:  deviceID,
			AccountId: "af6f3796-c446-4850-9fda-65936cab9b6d",
		}
		mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(&device, nil)
		// Configure call to update visibility state
		mockIoT.EXPECT().SetVisibiltyState(gomock.Eq(device.DeviceId), gomock.Eq(params.status))
		// Configure call to send emails
		mockSQS.EXPECT().SendMessage(gomock.Eq(sqs.ConnectionStatusPayload{
			DeviceID:  deviceID,
			Connected: params.status,
			Time:      params.time,
		}))
		// Prepare an event
		event := DeviceLifecycleEvent{DeviceID: deviceID, EventType: params.event, Timestamp: params.timestamp}
		// Run the test
		assert.Nil(t, app.RunJob(nil, event))
	}
}

func getStubbedApp(t *testing.T) (*app, *MockIoTClient, *MockSQSClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock iot
	iot := NewMockIoTClient(ctrl)
	// Create mock sqs
	sqs := NewMockSQSClient(ctrl)
	// Bundle up into an app
	return &app{iot: iot, sqs: sqs}, iot, sqs
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}