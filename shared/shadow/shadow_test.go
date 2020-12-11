package shadow

//go:generate mockgen -destination mock_iotdataplane.go -package shadow github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface IoTDataPlaneAPI

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetShadow(t *testing.T) {
	// Create some test iterations
	testParams := []struct {
		deviceID string
		payload  string
		error    error
		shadow   Shadow
	}{
		{
			deviceID: "63eda5eb-7f56-417f-88ed-44a9eb9e5f67",
			payload:  `{"state":{"desired":{"status":true},"reported":{"status":false},"delta":{"status":true}},"metadata":{"desired":{"status":{"timestamp":1584003580}},"reported":{"status":{"timestamp":1584803417}}},"version":50,"timestamp":1584810789}`,
			error:    nil,
			shadow: func() Shadow {
				// Create a shadow
				s := Shadow{}
				s.Timestamp = Timestamp{time.Unix(1584810789, 0)}
				s.State.Reported = map[string]interface{}{"status": false}
				s.Metadata.Reported = map[string]MetadataEntry{"status": {Timestamp: Timestamp{time.Unix(1584803417, 0)}}}
				return s
			}(),
		},
	}
	// Cycle through the tests
	for _, params := range testParams {
		// Create mock controller
		ctrl := gomock.NewController(t)
		// Create mock database client
		mock := NewMockIoTDataPlaneAPI(ctrl)
		// Create the unit under test
		client := client{
			dp: mock,
		}
		// Configure expectations
		mock.EXPECT().GetThingShadow(&iotdataplane.GetThingShadowInput{
			ThingName: aws.String(params.deviceID),
		}).Return(&iotdataplane.GetThingShadowOutput{Payload: []byte(params.payload)}, params.error)
		// Run the test
		shadow, err := client.Get(params.deviceID)
		if err != params.error {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		// Assert parts of the shadow
		if cmp.Equal(shadow, params.shadow) {
			t.Errorf("Unexpected shadow: %v", shadow)
		}
	}
}

func TestSetVisibilityStatus(t *testing.T) {
	// Create some test iterations
	testParams := []struct {
		deviceID string
		status   bool
		payload  string
	}{
		{
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			status:   true,
			payload:  `{"state":{"reported":{"connection":true}}}`,
		},
		{
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			status:   false,
			payload:  `{"state":{"reported":{"connection":false}}}`,
		},
	}
	// Iterate the tests
	for _, params := range testParams {
		// Create mock controller
		ctrl := gomock.NewController(t)
		// Create mock database client
		mock := NewMockIoTDataPlaneAPI(ctrl)
		// Create the unit under test
		client := client{
			dp: mock,
		}
		// Configure expectations
		mock.EXPECT().UpdateThingShadow(&iotdataplane.UpdateThingShadowInput{
			ThingName: aws.String(params.deviceID),
			Payload:   []byte(params.payload),
		})
		// Run the test
		err := client.UpdateConnectionStatus(params.deviceID, params.status)
		assert.Nil(t, err)
	}
}

func TestGetConnectionStatus(t *testing.T) {
	// Create some test iterations
	testParams := []struct {
		deviceID string
		payload  string
		status   bool
		time     time.Time
	}{
		{
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:  `{"state":{"reported":{"connection":true}},"metadata":{"reported":{"connection":{"timestamp":1584803417}}}}`,
			status:   true,
			time:     time.Unix(1584803417, 0),
		},
		{
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:  `{"state":{"reported":{"connection":false}},"metadata":{"reported":{"connection":{"timestamp":1584803417}}}}`,
			status:   false,
			time:     time.Unix(1584803417, 0),
		},
	}
	// Iterate the tests
	for _, params := range testParams {
		// Create mock controller
		ctrl := gomock.NewController(t)
		// Create mock database client
		mock := NewMockIoTDataPlaneAPI(ctrl)
		// Create the unit under test
		client := client{
			dp: mock,
		}
		// Configure expectations
		mock.EXPECT().GetThingShadow(&iotdataplane.GetThingShadowInput{
			ThingName: aws.String(params.deviceID),
		}).Return(&iotdataplane.GetThingShadowOutput{Payload: []byte(params.payload)}, nil)
		// Run the test
		state, err := client.GetConnectionStatus(params.deviceID)
		// Assert the result
		assert.Nil(t, err)
		assert.Equal(t, params.status, state.State)
		assert.Equal(t, params.time, state.Timestamp.Time)
	}
}
