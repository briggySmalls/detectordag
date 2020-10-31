package sqs

//go:generate mockgen -destination mock_sqs.go -package sqs github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI

import (
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
	// Configure mock to expect a call
	isqs.EXPECT().SendMessage(gomock.Not(gomock.Nil())).Do(func(input *sqs.SendMessageInput) {
		assert.Equal(t, `{"connected":true,"time":"1970-01-01T00:00:00Z"}`, *input.MessageBody)
		assert.Equal(t, QueueUrl, *input.QueueUrl)
	}).Return(nil, nil)
	// Make the call
	client.SendMessage(ConnectionStatusPayload{
		Connected: true,
		Time:      time.Unix(0, 0).UTC(),
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
