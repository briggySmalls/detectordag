package sqs

//go:generate go run github.com/golang/mock/mockgen -destination mock_sqs.go -package sqs github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
)

const (
	QueueUrl = "myqueuename"
)

func TestSend(t *testing.T) {
	// Create the unit under test
	client, isqs := createUnitAndMocks(t)
	// Set some test parameters
	const (
		deviceId = "573b0564-12f1-47fb-adf5-2d0906b39123"
		eventID  = "186b8a97-bd3e-43fc-ade1-e1d4f66bcb18"
		status   = "connected"
	)
	// Configure mock to expect a call
	isqs.EXPECT().SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(fmt.Sprintf(
			`{"deviceId":"%s","status":"%s","time":"1970-01-01T00:00:00Z","id":"%s"}`,
			deviceId, status, eventID)),
		QueueUrl: aws.String(QueueUrl),
	}).Return(nil, nil)
	// Make the call
	client.QueueConnectionEvent(ConnectionEventPayload{
		DeviceID: deviceId,
		Time:     time.Unix(0, 0).UTC(),
		Status:   status,
		ID:       eventID,
	})
}

func createUnitAndMocks(t *testing.T) (Client, *MockSQSAPI) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock SQSAPI
	mock := NewMockSQSAPI(ctrl)
	// Create the unit under test
	sender := client{
		sqs:      mock,
		queueUrl: QueueUrl,
	}
	// Create the new iot client
	return &sender, mock
}
