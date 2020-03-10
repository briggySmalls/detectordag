package app

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
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

// dbInit sets up the session and define table name, primary key, and sort key
func DbInit(session *session.Session) error {
	// Create Amazon DynamoDB client
	db = dynamodb.New(session)
	if db == nil {
		return fmt.Errorf("Failed to create database client")
	}
	return nil
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
