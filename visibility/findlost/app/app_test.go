package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_visibility.go -package app -mock_names Client=MockVisibilityEmailClient github.com/briggysmalls/detectordag/visibility EmailClient

import (
	"errors"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const lastSeenDurationHours = 24

func TestDeviceLookupFailed(t *testing.T) {
	// Create app under test
	app, mockIoT, _, _ := getStubbedApp(t)
	// Configure lookup to fail
	mockIoT.EXPECT().GetThingsByVisibility(gomock.Eq(true)).Return(nil, errors.New("Something went wrong"))
	// Run test
	assert.NotNil(t, app.RunJob(nil))
}

func TestNoDevices(t *testing.T) {
	// Create app under test
	app, mockIoT, _, _ := getStubbedApp(t)
	// Configure lookup to fail
	mockIoT.EXPECT().GetThingsByVisibility(gomock.Eq(true)).Return([]*iot.Device{}, nil)
	// Run test
	assert.Nil(t, app.RunJob(nil))
}

func TestRaceConditionDevices(t *testing.T) {
	// Create app under test
	app, mockIoT, _, _ := getStubbedApp(t)
	// Configure lookup to return a device marked as lost
	devices := []*iot.Device{
		{Visibility: false},
		{Visibility: false},
		{Visibility: false},
		{Visibility: false},
	}
	mockIoT.EXPECT().GetThingsByVisibility(gomock.Eq(true)).Return(devices, nil)
	// Run test
	assert.Nil(t, app.RunJob(nil))
}

func TestOkDevice(t *testing.T) {
	// Create app under test
	app, mockIoT, _, mockShadow := getStubbedApp(t)
	// Configure lookup to return a device marked as visible
	device := iot.Device{
		Name:       "1",
		DeviceId:   "0688f31b-7ba4-4a4d-8302-55e49d0393e7",
		AccountId:  "af6f3796-c446-4850-9fda-65936cab9b6d",
		Visibility: true, // Device is 'visible'
	}
	mockIoT.EXPECT().GetThingsByVisibility(gomock.Eq(true)).Return([]*iot.Device{&device}, nil)
	// Configure shadow to indicate device was seen recently
	shadow := shadow.Shadow{
		State: shadow.State{
			Reported: map[string]interface{}{"status": true},
		},
		Metadata: shadow.Metadata{
			Reported: map[string]shadow.MetadataEntry{
				"status": {Timestamp: shadow.Timestamp{Time: time.Now().Add(-1 * (lastSeenDurationHours*time.Hour - time.Second))}},
			},
		},
	}
	mockShadow.EXPECT().Get(gomock.Eq(device.DeviceId)).Return(&shadow, nil)
	// Run test
	assert.Nil(t, app.RunJob(nil))
}

func TestLostDevice(t *testing.T) {
	// Create app under test
	app, mockIoT, mockEmail, mockShadow := getStubbedApp(t)
	// Configure lookup to return a device marked as visible
	device := iot.Device{
		Name:       "1",
		DeviceId:   "0688f31b-7ba4-4a4d-8302-55e49d0393e7",
		AccountId:  "af6f3796-c446-4850-9fda-65936cab9b6d",
		Visibility: true, // Device is 'visible'
	}
	mockIoT.EXPECT().GetThingsByVisibility(gomock.Eq(true)).Return([]*iot.Device{&device}, nil)
	// Configure shadow to indicate device hasn't been seen for a while
	tooLongAgoTime := time.Now().Add(-1 * (lastSeenDurationHours*time.Hour + time.Second))
	shadow := shadow.Shadow{
		State: shadow.State{
			Reported: map[string]interface{}{"status": true},
		},
		Metadata: shadow.Metadata{
			Reported: map[string]shadow.MetadataEntry{
				"status": {Timestamp: shadow.Timestamp{Time: tooLongAgoTime}},
			},
		},
	}
	mockShadow.EXPECT().Get(gomock.Eq(device.DeviceId)).Return(&shadow, nil)
	// Configure call to update status to "lost"
	mockIoT.EXPECT().SetVisibiltyState(gomock.Eq(device.DeviceId), false).Return(nil)
	// Configure call to send emails
	mockEmail.EXPECT().SendVisibilityStatus(gomock.Eq(&device), gomock.Eq(tooLongAgoTime), gomock.Eq(false))
	// Run test
	assert.Nil(t, app.RunJob(nil))
}

func getStubbedApp(t *testing.T) (*app, *MockIoTClient, *MockEmailClient, *MockShadowClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock iot
	iot := NewMockIoTClient(ctrl)
	// Create mock visibility
	email := NewMockEmailClient(ctrl)
	// Create mock shadow
	shadow := NewMockShadowClient(ctrl)
	// Create a 'last seen duration'
	lastSeenDuration, err := time.ParseDuration(fmt.Sprintf("%dh", lastSeenDurationHours))
	assert.Nil(t, err)
	// Bundle up into an app
	return &app{iot: iot, email: email, shadow: shadow, lastSeenDuration: lastSeenDuration}, iot, email, shadow
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
