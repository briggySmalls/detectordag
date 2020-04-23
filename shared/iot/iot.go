package iot

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

const (
	accountIDAttributeName = "account-id"
)

var (
	ErrAccountIDMissing = errors.New("The account-id attribute was missing")
)

type client struct {
	iot *iot.IoT
}

type Client interface {
	GetThing(id string) (Device, error)
	GetThingsByAccount(id string) ([]Device, error)
}

type thingAttribute struct {
	*iot.ThingAttribute
}

type describeThingOutput struct {
	*iot.DescribeThingOutput
}

type Device interface {
	AccountID() (string, error)
}

// New gets a new Client
func New(sesh *session.Session) (Client, error) {
	// Create Amazon DynamoDB client
	iot := iot.New(sesh)
	if iot == nil {
		return nil, errors.New("Failed to create database client")
	}
	// Create our client wrapper
	client := client{
		iot: iot,
	}
	return &client, nil
}

// GetThing gets a thing with the given name from the AWS IoT registry
func (c *client) GetThing(id string) (Device, error) {
	thing, err := c.iot.DescribeThing(&iot.DescribeThingInput{ThingName: aws.String(id)})
	return &describeThingOutput{thing}, err
}

// GetThingsByAccount returns all things which are associated with the specified accountg
func (c *client) GetThingsByAccount(id string) ([]Device, error) {
	// Search for things
	things := []*iot.ThingAttribute{}
	var err error
	things, err = c.getPaginatedThings(&iot.ListThingsInput{
		AttributeName:  aws.String(accountIDAttributeName),
		AttributeValue: aws.String(id),
	}, nil, things)
	if err != nil {
		return nil, err
	}
	// Wrap up each thing
	wrappedThings := make([]Device, len(things))
	for i, thing := range things {
		wrappedThings[i] = &thingAttribute{thing}
	}
	return wrappedThings, nil
}

func (c *client) getPaginatedThings(input *iot.ListThingsInput, output *iot.ListThingsOutput, things []*iot.ThingAttribute) ([]*iot.ThingAttribute, error) {
	// Request the things
	var err error
	if output == nil {
		// This is the first request so just use the input
		output, err = c.iot.ListThings(input)
	} else {
		// We are making a paginated request, so use the 'next token'
		output, err = c.iot.ListThings(input.SetNextToken(*output.NextToken))
	}
	// Return if there is an error
	if err != nil {
		return nil, err
	}
	// Add the things
	things = append(things, output.Things...)
	// Short circuit if there are no more requests to make
	if output.NextToken == nil {
		return things, nil
	}
	// Recursively request more things
	return c.getPaginatedThings(input, output, things)
}

// AccountID gets the account ID for the thing
func (t *thingAttribute) AccountID() (string, error) {
	return t.getAttribute(accountIDAttributeName)
}

func (t *thingAttribute) getAttribute(key string) (string, error) {
	if accountID, ok := t.ThingAttribute.Attributes[key]; ok {
		return *accountID, nil
	}
	return "", ErrAccountIDMissing
}

// AccountID gets the account ID for the thing
func (t *describeThingOutput) AccountID() (string, error) {
	return t.getAttribute(accountIDAttributeName)
}

func (t *describeThingOutput) getAttribute(key string) (string, error) {
	if accountID, ok := t.DescribeThingOutput.Attributes[key]; ok {
		return *accountID, nil
	}
	return "", ErrAccountIDMissing
}
