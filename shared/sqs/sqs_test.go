package sqs

//go:generate mockgen -destination mock_sqs.go -package sqs github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	)
	// Configure mock to expect a call
	isqs.EXPECT().SendMessage(gomock.Not(gomock.Nil())).Do(func(input *sqs.SendMessageInput) {
		assert.Equal(t, fmt.Sprintf(`{"deviceId":"%s","time":"1970-01-01T00:00:00Z"}`, deviceId), *input.MessageBody)
		assert.Equal(t, QueueUrl, *input.QueueUrl)
	}).Return(nil, nil)
	// Make the call
	client.QueueDisconnectedEvent(DisconnectedPayload{
		DeviceID: deviceId,
		Time:     time.Unix(0, 0).UTC(),
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
