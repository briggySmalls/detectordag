package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package app -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_connection_updater.go -package app github.com/briggysmalls/detectordag/visibility ConnectionUpdater

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	// "github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInvalidPayload(t *testing.T) {
	testParams := []struct {
		event string
	}{
		{event: "other"},
		{event: `{"dummy": "text"}`},
		{event: `{"deviceId:"e35238bb-ca2c-4e2b-88da-3d305ffe904c"}`},
		{event: `{"deviceId":"not-uuid","time":"2020-12-12T19:58:16+00:00"}`},
		{event: `{"deviceId:"e35238bb-ca2c-4e2b-88da-3d305ffe904c","time":"some-bad-time"}`},
	}
	// Run the test
	for _, params := range testParams {
		// Create app under test
		app, _, _, _ := getStubbedApp(t)
		// Prepare an event
		event := events.SQSEvent{Records: []events.SQSMessage{{Body: params.event}}}
		// Run the test
		assert.NotNil(t, app.Handler(nil, event))
	}
}

func TestConnectionStatusLookupFailed(t *testing.T) {
	const (
		deviceID = "e35238bb-ca2c-4e2b-88da-3d305ffe904c"
	)
	// Create app under test
	app, _, shadow, _ := getStubbedApp(t)
	// Construct the event
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{Body: `{"deviceId":"e35238bb-ca2c-4e2b-88da-3d305ffe904c","time":"2020-12-12T19:58:16+00:00"}`},
		},
	}
	// Expect a call to shadow
	shadow.EXPECT().GetConnectionStatus(deviceID).Return(nil, errors.New("Something went wrong"))
	// Run the test
	assert.NotNil(t, app.Handler(nil, event))
}

func TestStaleEvent(t *testing.T) {
	const (
		deviceID = "e35238bb-ca2c-4e2b-88da-3d305ffe904c"
	)
	// Create app under test
	app, _, mockShadow, _ := getStubbedApp(t)
	// Construct the event
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{Body: `{"deviceId":"e35238bb-ca2c-4e2b-88da-3d305ffe904c","time":"2020-12-12T19:58:16+00:00"}`},
		},
	}
	// Expect a call to shadow
	mockShadow.EXPECT().GetConnectionStatus(deviceID).Return(
		&shadow.ConnectionState{State: true, Updated: createTime(t, "2020/12/12 19:58:17")},
		nil,
	)
	// Run the test
	assert.Nil(t, app.Handler(nil, event))
}

// func TestDeviceLookupFailed(t *testing.T) {
// 	// Create app under test
// 	app, mockIoT, _ := getStubbedApp(t)
// 	// Prepare some test parameters
// 	const (
// 		deviceID = "b6d62b30-00ac-49c4-9268-88559a46889f"
// 	)
// 	// Prepare an event
// 	event := DeviceLifecycleEvent{
// 		DeviceID:  deviceID,
// 		Timestamp: 0,
// 		EventType: "connected",
// 	}
// 	// Configure lookup to fail
// 	mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(nil, errors.New("Something went wrong"))
// 	// Run test
// 	assert.NotNil(t, app.RunJob(nil, event))
// }

// func TestEmailsSent(t *testing.T) {
// 	testParams := []struct {
// 		event     string
// 		status    bool
// 		timestamp int64
// 		time      time.Time
// 	}{
// 		{event: "connected", status: true, timestamp: 0, time: createTime(t, "1970/01/01 00:00:00")},
// 		{event: "disconnected", status: false, timestamp: 0, time: createTime(t, "1970/01/01 00:00:00")},
// 	}
// 	const (
// 		deviceID = "792ac520-0733-4ffe-8137-8aba3ca446d7"
// 	)
// 	// Run the test
// 	for _, params := range testParams {
// 		// Create app under test
// 		app, mockIoT, mockEmail := getStubbedApp(t)
// 		// Configure lookup to pass
// 		device := iot.Device{
// 			Name:      "1",
// 			DeviceId:  deviceID,
// 			AccountId: "af6f3796-c446-4850-9fda-65936cab9b6d",
// 		}
// 		mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(&device, nil)
// 		// Configure call to update visibility state
// 		mockIoT.EXPECT().SetVisibiltyState(gomock.Eq(device.DeviceId), gomock.Eq(params.status))
// 		// Configure call to send emails
// 		mockEmail.EXPECT().SendVisibilityStatus(gomock.Eq(&device), gomock.Eq(params.time), gomock.Eq(params.status))
// 		// Prepare an event
// 		event := DeviceLifecycleEvent{DeviceID: deviceID, EventType: params.event, Timestamp: params.timestamp}
// 		// Run the test
// 		assert.Nil(t, app.RunJob(nil, event))
// 	}
// }

func getStubbedApp(t *testing.T) (*app, *MockIoTClient, *MockShadowClient, *MockConnectionUpdater) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock shadow
	shadow := NewMockShadowClient(ctrl)
	// Create mock iot
	iot := NewMockIoTClient(ctrl)
	// Create mock connection updater
	updater := NewMockConnectionUpdater(ctrl)
	// Bundle up into an app
	return &app{iot: iot, shadow: shadow, updater: updater}, iot, shadow, updater
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
