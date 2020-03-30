package database

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	ACCOUNTS_TABLE    = "accounts"
	DEVICES_TABLE     = "devices"
	ACCOUNTS_GSI_NAME = "username-index"
	DEVICES_GSI_NAME  = "account-id-index"
)

type client struct {
	db *dynamodb.DynamoDB
}

// Client is a client for interfacing with a detectordag database
type Client interface {
	GetDeviceById(id string) (*Device, error)
	GetAccountById(id string) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)
	GetDevicesByAccount(id string) ([]Device, error)
	UpdateAccountEmails(accountId string, emails []string) (*Account, error)
}

// account represents an 'accounts' table entry
type Account struct {
	AccountId string   `dynamodbav:"account-id"` // TODO: unmarshal into our own UUID type
	Emails    []string `dynamodbav:"emails"`
	Username  string   `dynamodbav:"username"`
	Password  string   `dynamodbav:"password"`
}

// device is a 'device' table row
type Device struct {
	Name      string `dynamodbav:"name"`
	DeviceId  string `dynamodbav:"device-id"`
	AccountId string `dynamodbav:"account-id"`
}

// New gets a new Client
func New(sesh *session.Session) (Client, error) {
	// Create Amazon DynamoDB client
	db := dynamodb.New(sesh)
	if db == nil {
		return nil, errors.New("Failed to create database client")
	}
	// Create our client wrapper
	client := client{
		db: db,
	}
	return &client, nil
}

func (d *client) GetDeviceById(id string) (*Device, error) {
	// Request for the device associated with the ID
	result, err := d.db.GetItem(&dynamodb.GetItemInput{
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
	kc := expression.Key("account-id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(kc).Build()
	if err != nil {
		return nil, err
	}
	// Request for the devices associated with the account
	result, err := d.db.Query(&dynamodb.QueryInput{
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
	return devices, nil
}

func (d *client) GetAccountById(id string) (*Account, error) {
	// Request for the account associated with the device
	result, err := d.db.GetItem(&dynamodb.GetItemInput{
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
	result, err := d.db.Query(&dynamodb.QueryInput{
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

func (d *client) UpdateAccountEmails(accountId string, emails []string) (*Account, error) {
	// Build an update expression
	update := expression.Set(
		expression.Name("emails"),
		expression.Value(emails),
	)
	// Create the DynamoDB expression from the Update.
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, err
	}
	// Update the emails (request updated response)
	result, err := d.db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:                 aws.String(ACCOUNTS_TABLE),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       map[string]*dynamodb.AttributeValue{"account-id": {S: aws.String(accountId)}},
		UpdateExpression:          expr.Update(),
		ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
	})
	if err != nil {
		return nil, err
	}
	return unmarshalAccount(result.Attributes)
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
