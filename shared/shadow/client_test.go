package shadow

//go:generate go run github.com/golang/mock/mockgen -destination mock_iotdataplane.go -package shadow github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface IoTDataPlaneAPI

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
			payload:  `{"metadata":{"reported":{"connection":{"timestamp":1584803417},"status":{"timestamp":1584803414}}},"state":{"reported":{"connection":"connected","status":"off"}},"timestamp":1584810789,"version":50}`,
			error:    nil,
			shadow: Shadow{
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: StringShadowField{
					Value:   CONNECTION_STATUS_CONNECTED,
					Updated: time.Unix(1584803417, 0),
				},
				Power: StringShadowField{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
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
