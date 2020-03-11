package app

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strconv"
	"time"
)

const (
	ACCOUNTS_TABLE = "accounts"
	DEVICES_TABLE  = "devices"
)

// account represents an 'accounts' table entry
type Account struct {
	AccountId int      `dynamodbav:"account-id"`
	Emails    []string `dynamodbav:"emails"`
}

// device is a 'device' table row
type Device struct {
	DeviceId  string `dynamodbav:"device-id"`
	AccountId int    `dynamodbav:"account-id"`
}

//It is a best practice to instanciate the Amazon DynamoDB client outside
//of the AWS Lambda function handler.
//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.Lambda.BestPracticesWithDynamoDB.html
var db *dynamodb.DynamoDB

// init sets up the session
func init() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	var err error
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Create Amazon DynamoDB client
	db = dynamodb.New(sesh)
	if db == nil {
		log.Fatal("Failed to create database client")
	}
}

func updateDevice(update PowerStatusChangedEvent) (*Device, error) {
	// Create an expression for updating the row
	input := &dynamodb.UpdateItemInput{
		// Look up the device of interest
		TableName: aws.String("devices"),
		Key: map[string]*dynamodb.AttributeValue{
			"device-id": {
				S: aws.String(update.DeviceId),
			},
		},
		// Update the 'status' and 'last-updated' fields
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("status"),
			"#T": aws.String("last-update"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				BOOL: aws.Bool(update.Status),
			},
			":t": {
				S: aws.String(update.Timestamp.Format(time.RFC3339)),
			},
		},
		UpdateExpression: aws.String("SET #S = :s, #T = :t"),
		// Only update if this is more recent (or never set at all)
		ConditionExpression: aws.String("attribute_not_exists(#T) or #T < :t"),
		// Return all the attributes (we will use them to look up account)
		ReturnValues: aws.String("ALL_NEW"),
	}
	// Run the update operation
	result, err := db.UpdateItem(input)
	if err != nil {
		return nil, err
	}
	log.Printf("Updated device: %s", result)
	// Pull out the device attributes
	device := Device{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &device)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func getAccount(id int) (*Account, error) {
	// Request for the account associated with the device
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ACCOUNTS_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"account-id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	// Check we got exactly one account
	if result.Item == nil {
		return nil, fmt.Errorf("Unknown account: %d", id)
	}
	// Unmarshal the account
	account := Account{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
