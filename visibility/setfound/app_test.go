package main

//go:generate go run github.com/golang/mock/mockgen -destination mock_iot.go -package main -mock_names Client=MockIoTClient -self_package github.com/briggysmalls/detectordag/visibility/setfound github.com/briggysmalls/detectordag/shared/iot Client
//go:generate go run github.com/golang/mock/mockgen -destination mock_visibility.go -package main -mock_names Client=MockVisibilityEmailClient -self_package github.com/briggysmalls/detectordag/visibility/setfound github.com/briggysmalls/detectordag/visibility EmailClient

import (
	"errors"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeviceLookupFailed(t *testing.T) {
	// Create app under test
	app, mockIoT, _ := getStubbedApp(t)
	// Create some test parameters
	const deviceID = "f88948e6-5f93-4f11-8d58-15d48075069d"
	// Configure lookup to fail
	mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(nil, errors.New("Something went wrong"))
	// Run test
	assert.NotNil(t, app.handleRequest(nil, DeviceSeenEvent{
		DeviceId: deviceID,
		Updated: updated{
			Status: timestamp{
				Timestamp: 0,
			},
		},
	}))
}

func TestAlreadySeen(t *testing.T) {
	// Create app under test
	app, mockIoT, _ := getStubbedApp(t)
	// Create some test parameters
	const (
		deviceID   = "f88948e6-5f93-4f11-8d58-15d48075069d"
		accountID  = "e8a5a68a-13bc-4054-a3a1-7d3b0028e8dd"
		deviceName = "Alderney"
	)
	// Configure mocks
	// Successfully look up a 'thing' that is 'found'
	device := iot.Device{Visibility: true, Name: deviceName, AccountId: accountID, DeviceId: deviceID}
	mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(&device, nil)
	// Run the handler
	event := DeviceSeenEvent{
		DeviceId: deviceID,
		Updated: updated{
			Status: timestamp{
				Timestamp: 0,
			},
		},
	}
	assert.Nil(t, app.handleRequest(nil, event))
}

func TestSetFound(t *testing.T) {
	// Create app under test
	app, mockIoT, mockEmail := getStubbedApp(t)
	// Create some test parameters
	const (
		deviceID   = "f88948e6-5f93-4f11-8d58-15d48075069d"
		accountID  = "e8a5a68a-13bc-4054-a3a1-7d3b0028e8dd"
		deviceName = "Alderney"
		updateTime = 1584063049
	)
	// Configure mocks
	// Successfully look up a 'thing' that is 'lost'
	device := iot.Device{Visibility: true, Name: deviceName, AccountId: accountID, DeviceId: deviceID}
	mockIoT.EXPECT().GetThing(gomock.Eq(deviceID)).Return(&device, nil)
	// Successfully email
	mockEmail.EXPECT().SendVisibilityStatus(
		gomock.Eq(&device),
		gomock.Eq(createTime(t, "2020/03/13 01:30:49")),
		gomock.Eq(true),
	)
	// Update the visibility status
	mockIoT.EXPECT().SetVisibiltyState(gomock.Eq(deviceID), gomock.Eq(true))
	// Run the handler
	event := DeviceSeenEvent{
		DeviceId: deviceID,
		Updated: updated{
			Status: timestamp{
				Timestamp: updateTime,
			},
		},
	}
	assert.Nil(t, app.handleRequest(nil, event))
}

func getStubbedApp(t *testing.T) (*app, *MockIoTClient, *MockEmailClient) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock iot
	iot := NewMockIoTClient(ctrl)
	// Create mock visibility
	email := NewMockEmailClient(ctrl)
	// Bundle up into an app
	return &app{iot: iot, email: email}, iot, email
}

func createTime(t *testing.T, timeString string) time.Time {
	tme, err := time.Parse("2006/01/02 15:04:05", timeString)
	assert.NoError(t, err)
	return tme
}
