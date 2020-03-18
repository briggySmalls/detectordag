package database

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"log"
)

const (
	ACCOUNTS_TABLE    = "accounts"
	DEVICES_TABLE     = "devices"
	ACCOUNTS_GSI_NAME = "username-index"
	ACCOUNTS_GSI_NAME = "account-id-index"
)

type client struct {
}

// Client is a client for interfacing with a detectordag database
type Client interface {
	GetDeviceById(id string) (*Device, error)
	GetAccountById(id string) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)
}

// account represents an 'accounts' table entry
type Account struct {
	AccountId string   `dynamodbav:"account-id"`
	Emails    []string `dynamodbav:"emails"`
	Username  string   `dynamodbav:"username"`
	Password  string   `dynamodbav:"password"`
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

// New gets a new Client
func New() Client {
	client := client{}
	return &client
}

func (d *client) GetDeviceById(id string) (*Device, error) {
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

func (d *client) GetDevicesByAccount(id string) ([]Device, error) {
	// Build an expression
	kc := expression.Key("#acc").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(kc).Build()
	if err != nil {
		return nil, err
	}
	// Request for the devices associated with the account
	result, err := db.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(DEVICES_TABLE),
		IndexName:                 aws.String(DEVICES_GSI_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Select:                    aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		return nil, err
	}
	// Unmarshal the devices
	var devices []Device
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &devices)
	if err != nil {
		return nil, err
	}
	return devices
}

func (d *client) GetAccountById(id string) (*Account, error) {
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
		return nil, fmt.Errorf("Unknown account: %s", id)
	}
	// Unmarshal the account
	return unmarshalAccount(result.Item)
}

func (d *client) GetAccountByUsername(username string) (*Account, error) {
	// Build an expression
	kc := expression.Key("username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithKeyCondition(kc).Build()
	if err != nil {
		return nil, err
	}
	// Request for the account associated with the username
	result, err := db.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(ACCOUNTS_TABLE),
		IndexName:                 aws.String(ACCOUNTS_GSI_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Select:                    aws.String("ALL_PROJECTED_ATTRIBUTES"),
	})
	if err != nil {
		return nil, err
	}
	// Check we got exactly one account
	if len(result.Items) != 1 {
		return nil, fmt.Errorf("Unknown account: %s", username)
	}
	return unmarshalAccount(result.Items[0])
}

func unmarshalAccount(item map[string]*dynamodb.AttributeValue) (*Account, error) {
	// Unmarshal the account
	account := Account{}
	err := dynamodbattribute.UnmarshalMap(item, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
