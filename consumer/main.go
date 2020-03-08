package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strconv"
)

const ACCOUNTS_TABLE = "accounts"
const DEVICES_TABLE = "devices"

type PowerStatusChangedEvent struct {
	DeviceId  string `json:""`
	Timestamp string `json:""`
	Version   string `json:""`
	Status    bool   `json:""`
}

type Account struct {
	AccountId int
	Emails    []string
}

type Device struct {
	DeviceId  string
	AccountId int
}

//init set up the session and define table name, primary key, and sort key
func dbInit() *dynamodb.DynamoDB {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	dbSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	if dbSession == nil {
		log.Fatal("Failed to start session")
	}
	// Create Amazon DynamoDB client
	db := dynamodb.New(dbSession)
	if db == nil {
		log.Fatal("Failed to create database client")
	}
	return db
}

//It is a best practice to instanciate the Amazon DynamoDB client outside
//of the AWS Lambda function handler.
//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Streams.Lambda.BestPracticesWithDynamoDB.html
var db = dbInit()

func HandleRequest(ctx context.Context, event PowerStatusChangedEvent) {
	// Get the device ID
	device, err := getDevice(event.DeviceId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account ID: %d", device.AccountId)
	// Get the account
	account, err := getAccount(device.AccountId)
	if err != nil {
		log.Fatal(err)
	}
	// Send 'power status updated' emails
	for email := range account.Emails {
		log.Printf("Send email to: %s", email)
	}
}

func getDevice(id string) (*Device, error) {
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

func main() {
	lambda.Start(HandleRequest)
}
