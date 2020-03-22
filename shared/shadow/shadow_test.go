package shadow

//go:generate mockgen -destination mock_iotdataplane.go -package shadow github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface IoTDataPlaneAPI

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/golang/mock/gomock"
	"testing"
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
			shadow:   {Timestamp: time.Unix(1584810789, 0), Metadata: {Reported: {"status": {Timestamp: time.Unix(1584803417, 0)}}}, State: {Reported: {"status": false}}},
			error:    nil,
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
		mock.EXPECT().GetThingShadow(gomock.Eq(iotdataplane.GetThingShadowInput{
			ThingName: aws.String(params.deviceID),
		})).Return(&iotdataplane.GetThingShadowOutput{Payload: []byte(params.payload)}, params.error)
		// Run the test
		shadow, err := client.Get(params.deviceID)
		if err != params.error {
			t.Errorf("Unexpected error: %v", err)
		} else if shadow != params.payload {
			t.Errorf("Unexpected shadow: %v", shadow)
		}
	}
}
