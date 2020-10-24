package app

//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package app -mock_names Client=MockIoTClient github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_visibility.go -package app -mock_names Client=MockVisibilityEmailClient github.com/briggysmalls/detectordag/visibility EmailClient
//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package app -mock_names Client=MockShadowClient github.com/briggysmalls/detectordag/shared/shadow Client

import (
	"errors"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
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
