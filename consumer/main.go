package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

const ACCOUNTS_TABLE = "accounts"
const DEVICES_TABLE = "devices"

type PowerStatusChangedEvent struct {
	DeviceId  string `json:""`
	Timestamp string `json:""`
	Version   string `json:""`
	Status    bool   `json:""`
}

//init set up the session and define table name, primary key, and sort key
func dbInit() *dynamodb.DynamoDB {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	dbSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create Amazon DynamoDB client
	db := dynamodb.New(dbSession)
	return db
}

//It is a best practice to instanciate the Amazon DynamoDB client outside
//of the AWS Lambda function handler.
//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.Lambda.BestPracticesWithDynamoDB.html
var db = dbInit()

func HandleRequest(ctx context.Context, event PowerStatusChangedEvent) {
	// Request for the device associated with the ID
	device, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(DEVICES_TABLE),
		Key:       map[string]*dynamodb.AttributeValue{"device-id": event.DeviceId},
	})
	if err != nil {
		log.Fatal()
	}
	// Request for the account associated with the device
	account, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ACCOUNTS_TABLE),
		Key:       map[string]*dynamodb.AttributeValue{"account-id": device.account - id},
	})
	if err != nil {
		log.Fatal()
	}
	// Pull out the emails
	emails, err := account.emails
	// Send 'power status updated' emails
	for email := range emails {
		log.Printf("Send email to: %s", email)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
