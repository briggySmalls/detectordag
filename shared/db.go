package shared

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const (
	ACCOUNTS_TABLE = "accounts"
	DEVICES_TABLE  = "devices"
)

// account represents an 'accounts' table entry
type Account struct {
	AccountId string   `dynamodbav:"account-id"`
	Emails    []string `dynamodbav:"emails"`
}

// device is a 'device' table row
type Device struct {
	DeviceId  string `dynamodbav:"device-id"`
	AccountId string `dynamodbav:"account-id"`
}

//It is a best practice to instantiate the Amazon DynamoDB client outside
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

func GetDevice(id string) (*Device, error) {
	// Request for the device associated with the ID
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(DEVICES_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"device-id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	// Check we got exactly one device
	if result.Item == nil {
		return nil, fmt.Errorf("Unknown device: %s", id)
	}
	// Unmarshal the device
	device := Device{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &device)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func GetAccount(id string) (*Account, error) {
	// Request for the account associated with the device
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ACCOUNTS_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"account-id": {
				S: aws.String(id),
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
