package main

//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package main -mock_names Client=MockIoTClient -self_package github.com/briggysmalls/detectordag/visibility/findlost github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_visibility.go -package main -mock_names Client=MockVisibilityEmailClient -self_package github.com/briggysmalls/detectordag/visibility/findlost github.com/briggysmalls/detectordag/visibility EmailClient
//go:generate go run github.com/golang/mock/mockgen -destination mock_shadow.go -package main -mock_names Client=MockShadowClient -self_package github.com/briggysmalls/detectordag/visibility/findlost github.com/briggysmalls/detectordag/shared/shadow Client

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeviceLookupFailed(t *testing.T) {
	// Create app under test
	app, mockIoT, _, _ := getStubbedApp(t)
	// Create some test parameters
	const deviceID = "f88948e6-5f93-4f11-8d58-15d48075069d"
	// Configure lookup to fail
	mockIoT.EXPECT().GetThingsByVisibility(gomock.Eq(true)).Return(nil, errors.New("Something went wrong"))
	// Run test
	assert.NotNil(t, app.runJob(nil))
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
	// Bundle up into an app
	return &app{iot: iot, email: email, shadow: shadow}, iot, email, shadow
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
