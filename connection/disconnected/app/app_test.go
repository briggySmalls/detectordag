package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package app -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_connection_updater.go -package app github.com/briggysmalls/detectordag/connection ConnectionUpdater

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInvalidPayload(t *testing.T) {
	testParams := []struct {
		event string
	}{
		{event: "other"},
		{event: `{"dummy":"text"}`},
		{event: `{"deviceId:"e35238bb-ca2c-4e2b-88da-3d305ffe904c"}`},
		{event: `{"deviceId":"not-uuid","time":"2020-12-12T19:58:16+00:00"}`},
		{event: `{"deviceId:"e35238bb-ca2c-4e2b-88da-3d305ffe904c","time":"some-bad-time"}`},
	}
	// Run the test
	for _, params := range testParams {
		// Create app under test
		app, _, _ := getStubbedApp(t)
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
	app, shadow, _ := getStubbedApp(t)
	// Construct the event
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{Body: `{"deviceId":"e35238bb-ca2c-4e2b-88da-3d305ffe904c","time":"2020-12-12T19:58:16+00:00"}`},
		},
	}
	// Expect a call to shadow
	shadow.EXPECT().Get(deviceID).Return(nil, errors.New("Something went wrong"))
	// Run the test
	assert.NotNil(t, app.Handler(nil, event))
}

func TestStaleEvent(t *testing.T) {
	const (
		deviceID = "e35238bb-ca2c-4e2b-88da-3d305ffe904c"
	)
	// Create app under test
	app, mockShadow, _ := getStubbedApp(t)
	// Construct the event
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{Body: `{
				"deviceId":"e35238bb-ca2c-4e2b-88da-3d305ffe904c",
				"id":"4e0a66f2-c928-4ad5-8870-b0c72ded0ae4",
				"time":"2020-12-12T19:58:16+00:00",
				"type":"connected"
			}`},
		},
	}
	// Expect a call to shadow
	mockShadow.EXPECT().Get(deviceID).Return(
		&shadow.Shadow{Connection: shadow.ConnectionShadow{
			Status:      shadow.CONNECTION_STATUS_CONNECTED,
			TransientID: "52068a06-f89d-4256-9b64-48fa990088d9",
			Updated:     createTime(t, "2020/12/12 19:58:17"),
		}},
		nil,
	)
	// Run the test
	assert.Nil(t, app.Handler(nil, event))
}

func TestEmailsSent(t *testing.T) {
	// Prepare some test parameters
	const (
		deviceID     = "b6d62b30-00ac-49c4-9268-88559a46889f"
		accountID    = "c6d62b30-00ac-49c4-9268-88559a46889f"
		eventTimeStr = "2020-12-12T19:58:16+00:00"
		transientID  = "52068a06-f89d-4256-9b64-48fa990088d9"
	)
	connectionStatus := shadow.CONNECTION_STATUS_CONNECTED
	// Create app under test
	app, mockShadow, mockUpdater := getStubbedApp(t)
	eventTime, err := time.Parse(time.RFC3339, eventTimeStr)
	assert.Nil(t, err)
	gomock.InOrder(
		// Expect a call to shadow
		mockShadow.EXPECT().Get(deviceID).Return(
			// Indicate the status hasn't been updated for a while
			&shadow.Shadow{Connection: shadow.ConnectionShadow{
				Status:      shadow.CONNECTION_STATUS_CONNECTED,
				TransientID: transientID,
				Updated:     createTime(t, "2020/12/12 00:00:00"),
			}},
			nil,
		),
		// Expect a call to update status
		mockUpdater.EXPECT().UpdateConnectionStatus(deviceID, eventTime, connectionStatus),
	)
	// Run test
	// Prepare an event
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{Body: fmt.Sprintf(`{
				"deviceId":"%s",
				"id":"%s",
				"type":"%s",
				"time":"%s"
			}`, deviceID, transientID, connectionStatus, eventTimeStr)},
		},
	}
	assert.Nil(t, app.Handler(nil, event))
}

func getStubbedApp(t *testing.T) (*app, *MockShadowClient, *MockConnectionUpdater) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock shadow
	shadow := NewMockShadowClient(ctrl)
	// Create mock connection updater
	updater := NewMockConnectionUpdater(ctrl)
	// Bundle up into an app
	return &app{shadow: shadow, updater: updater}, shadow, updater
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
