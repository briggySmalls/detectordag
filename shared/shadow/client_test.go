package shadow

//go:generate go run github.com/golang/mock/mockgen -destination mock_iotdataplane.go -package shadow github.com/aws/aws-sdk-go/service/iotdataplane/iotdataplaneiface IoTDataPlaneAPI

import (
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
			payload: `{
				"metadata":{"reported":{
					"status":{"timestamp":1584803414}
				}},
				"state":{"reported":{
					"name":"hello world",
					"connection":{
						"current":"connected",
						"transientId":"efb3ed5f-5357-4ebd-843c-6f8e79b74eae",
						"updated":1584803417
					},
					"status":"off"
				}},
				"timestamp":1584810789,"version":50
			}`,
			error: nil,
			shadow: Shadow{
				Name:    "hello world",
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: ConnectionShadow{
					Status:      CONNECTION_STATUS_CONNECTED,
					TransientID: "efb3ed5f-5357-4ebd-843c-6f8e79b74eae",
					Updated:     time.Unix(1584803417, 0),
				},
				Power: PowerShadow{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
		},
		{ // Missing a name
			deviceID: "63eda5eb-7f56-417f-88ed-44a9eb9e5f67",
			payload: `{"metadata":{"reported":{
				"status":{"timestamp":1584803414}
			}},
			"state":{"reported":{
				"connection":{
					"current":"connected",
					"transientId":"f5dc1874-5ba1-4727-8366-35d8278ea3e4",
					"updated":1584803417
				},
				"status":"off"
			}},"timestamp":1584810789,"version":50}`,
			error: nil,
			shadow: Shadow{
				Name:    "",
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: ConnectionShadow{
					Status:      CONNECTION_STATUS_CONNECTED,
					Updated:     time.Unix(1584803417, 0),
					TransientID: "f5dc1874-5ba1-4727-8366-35d8278ea3e4",
				},
				Power: PowerShadow{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
		},
	}
	// Cycle through the tests
	for _, params := range testParams {
		// Create mocks
		client, mock := createStubbedClient(t)
		// Configure expectations
		mock.EXPECT().GetThingShadow(&iotdataplane.GetThingShadowInput{
			ThingName: aws.String(params.deviceID),
		}).Return(&iotdataplane.GetThingShadowOutput{Payload: []byte(params.payload)}, params.error)
		// Run the test
		shadow, err := client.Get(params.deviceID)
		// Assert the outcome
		assert.Equal(t, params.error, err)
		assert.Equal(t, params.shadow, *shadow)
	}
}

// A helper for executing UpdateConnectionStatus without arguments
func updateConnectionStatusFactory(id, status string, time time.Time) func(Client) (*Shadow, error) {
	return func(client Client) (*Shadow, error) {
		return client.UpdateConnectionStatus(id, status, time)
	}
}

// A helper for executing UpdateName without arguments
func updateNameFactory(id, name string) func(Client) (*Shadow, error) {
	return func(client Client) (*Shadow, error) {
		return client.UpdateName(id, name)
	}
}

func TestUpdateShadow(t *testing.T) {
	// Create some test iterations
	testParams := []struct {
		deviceID      string
		status        string
		payload       string
		returnPayload string
		shadow        *Shadow
		testFunc      func(Client) (*Shadow, error)
	}{
		{ // Update connection to 'connected'
			testFunc: updateConnectionStatusFactory(
				"eb49b2e7-fd3a-4c03-b47f-b819281475e5",
				CONNECTION_STATUS_CONNECTED,
				time.Unix(1584803417, 0),
			),
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:  `{"state":{"reported":{"connection":{"current":"connected","updated":1584803417}}}}`,
			returnPayload: `{"metadata":{"reported":{
					"status":{"timestamp":1584803414}
				}},
				"state":{"reported":{
					"name":"my dag",
					"connection":{
						"current":"connected",
						"transientId":"619eb763-d0ab-4513-aeeb-8ff6ad8a500e",
						"updated":1584803417
					},
					"status":"off"
				}},"timestamp":1584810789,"version":50}`,
			shadow: &Shadow{
				Name:    "my dag",
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: ConnectionShadow{
					Status:      CONNECTION_STATUS_CONNECTED,
					TransientID: "619eb763-d0ab-4513-aeeb-8ff6ad8a500e",
					Updated:     time.Unix(1584803417, 0),
				},
				Power: PowerShadow{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
		},
		{ // Update connection to 'disconnected'
			testFunc: updateConnectionStatusFactory(
				"eb49b2e7-fd3a-4c03-b47f-b819281475e5",
				CONNECTION_STATUS_DISCONNECTED,
				time.Unix(1584803417, 0),
			),
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:  `{"state":{"reported":{"connection":{"current":"disconnected","updated":1584803417}}}}`,
			returnPayload: `{"metadata":{"reported":{
				"status":{"timestamp":1584803414}
			}},
			"state":{"reported":{
				"name":"Annex",
				"connection":{
					"current":"disconnected",
					"transientId":"f5dc1874-5ba1-4727-8366-35d8278ea3e4",
					"updated":1584803417
				},
				"status":"off"
			}},"timestamp":1584810789,"version":50}`,
			shadow: &Shadow{
				Name:    "Annex",
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: ConnectionShadow{
					Status:      CONNECTION_STATUS_DISCONNECTED,
					Updated:     time.Unix(1584803417, 0),
					TransientID: "f5dc1874-5ba1-4727-8366-35d8278ea3e4",
				},
				Power: PowerShadow{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
		},
		{ // Update name to 'Hello'
			testFunc: updateNameFactory(
				"eb49b2e7-fd3a-4c03-b47f-b819281475e5",
				"Hello",
			),
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:  `{"state":{"reported":{"name":"Hello"}}}`,
			returnPayload: `{"metadata":{"reported":{
				"status":{"timestamp":1584803414}
			}},
			"state":{"reported":{
				"name":"Hello",
				"connection":{
					"current":"connected",
					"transientId":"18592df0-ecc9-4e44-acd6-1b63872a8cf3",
					"updated":1584803417
				},
				"status":"off"
			}},"timestamp":1584810789,"version":50}`,
			shadow: &Shadow{
				Name:    "Hello",
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: ConnectionShadow{
					Status:      CONNECTION_STATUS_CONNECTED,
					Updated:     time.Unix(1584803417, 0),
					TransientID: "18592df0-ecc9-4e44-acd6-1b63872a8cf3",
				},
				Power: PowerShadow{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
		},
		{ // Update name to "My Dag"
			testFunc: updateNameFactory(
				"eb49b2e7-fd3a-4c03-b47f-b819281475e5",
				"My Dag",
			),
			deviceID: "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:  `{"state":{"reported":{"name":"My Dag"}}}`,
			returnPayload: `{"metadata":{"reported":{"status":{"timestamp":1584803414}}},
			"state":{"reported":{
				"name":"My Dag",
				"connection":{
					"current":"disconnected",
					"transientId":"9e9b59ac-b6b6-491b-8c55-f2d502f653b9",
					"updated":1584803417
				},
				"status":"off"
			}},"timestamp":1584810789,"version":50}`,
			shadow: &Shadow{
				Name:    "My Dag",
				Time:    time.Unix(1584810789, 0),
				Version: 50,
				Connection: ConnectionShadow{
					Status:      CONNECTION_STATUS_DISCONNECTED,
					Updated:     time.Unix(1584803417, 0),
					TransientID: "9e9b59ac-b6b6-491b-8c55-f2d502f653b9",
				},
				Power: PowerShadow{
					Value:   POWER_STATUS_OFF,
					Updated: time.Unix(1584803414, 0),
				},
			},
		},
	}
	// Iterate the tests
	for i, params := range testParams {
		log.Printf("Test iteration: %d", i)
		// Create mocks
		client, mock := createStubbedClient(t)
		// Configure expectations
		mock.EXPECT().UpdateThingShadow(&iotdataplane.UpdateThingShadowInput{
			ThingName: aws.String(params.deviceID),
			Payload:   []byte(params.payload),
		})
		mock.EXPECT().GetThingShadow(&iotdataplane.GetThingShadowInput{
			ThingName: aws.String(params.deviceID),
		}).Return(
			&iotdataplane.GetThingShadowOutput{
				Payload: []byte(params.returnPayload),
			}, nil)
		// Run the test
		shadow, err := params.testFunc(client)
		assert.Nil(t, err)
		assert.Equal(t, params.shadow, shadow)
	}
}

func TestUpdateTransientID(t *testing.T) {
	// Create some test iterations
	testParams := []struct {
		deviceID    string
		payload     string
		transientID string
	}{
		{
			deviceID:    "eb49b2e7-fd3a-4c03-b47f-b819281475e5",
			payload:     `{"state":{"reported":{"connection":{"transientId":"9e9b59ac-b6b6-491b-8c55-f2d502f653b9"}}}}`,
			transientID: "9e9b59ac-b6b6-491b-8c55-f2d502f653b9",
		},
	}
	// Iterate the tests
	for i, params := range testParams {
		log.Printf("Test iteration: %d", i)
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
		err := client.UpdateConnectionTransientID(params.deviceID, params.transientID)
		assert.Nil(t, err)
	}
}

func TestRequestStatusUpdate(t *testing.T) {
	// Create mocks
	client, mock := createStubbedClient(t)
	// Expect a call
	const (
		deviceID = "eb49b2e7-fd3a-4c03-b47f-b819281475e5"
		topic    = "dags/eb49b2e7-fd3a-4c03-b47f-b819281475e5/status/request"
	)
	mock.EXPECT().Publish(&iotdataplane.PublishInput{
		Qos:     aws.Int64(1),
		Topic:   aws.String(topic),
		Payload: []byte("{}"),
	}).Return(nil, nil)
	// Run the test
	assert.Nil(t, client.RequestStatusUpdate(deviceID))
}

func createStubbedClient(t *testing.T) (Client, *MockIoTDataPlaneAPI) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock database client
	mock := NewMockIoTDataPlaneAPI(ctrl)
	// Create the unit under test
	client := client{
		dp: mock,
	}
	return &client, mock
}
