package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-lambda-go/lambda"
)

type PowerStatusChangedEvent struct {
	DeviceId  string `json:""`
	Timestamp string `json:""`
	Version   string `json:""`
	Status    bool   `json:""`
}

//init set up the session and define table name, primary key, and sort key
func init(tn string, pk string, sk string) *DynamoDb {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	dbSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create Amazon DynamoDB client
	db := dynamodb.New(sess)
	return db
}

//It is a best practice to instanciate the Amazon DynamoDB client outside
//of the AWS Lambda function handler.
//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.Lambda.BestPracticesWithDynamoDB.html
var db = init(DB_TABLE_CONFIG_NAME, DB_TABLE_CONFIG_PK, DB_TABLE_CONFIG_SK)

func HandleRequest(ctx context.Context, event PowerStatusChangedEvent) (string, error) {
	// Request for the device associated with the ID
	device, err := db.GetItem(db.GetItemInput{
		TableName: "devices",
		Key: map[string]*dynamodb.AttributeValueevent{"device-id": event.DeviceId}
	})
	if err != nil {
		log.Fatal()
	}
	// Request for the account associated with the device
	account, err := db.GetItem(db.GetItemInput{
		TableName: "accounts",
		Key: map[string]*dynamodb.AttributeValueevent{"account-id": device.account-id}
	})
	if err != nil {
		log.Fatal()
	}
	// Pull out the emails
	emails, err := account.emails
	// Send 'power status updated' emails
	for email := range emails {

	}
	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
