package iot

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

const (
	accountIdAttributeName = "account-id"
)

type client struct {
	iot *iot.IoT
}

type Client interface {
	GetThing(id string) (*iot.DescribeThingOutput, error)
	GetThingsByAccount(id string) ([]*iot.ThingAttribute, error)
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
func (c *client) GetThing(id string) (*iot.DescribeThingOutput, error) {
	return c.iot.DescribeThing(&iot.DescribeThingInput{ThingName: aws.String(id)})
}

// GetThingsByAccount returns all things which are associated with the specified accountg
func (c *client) GetThingsByAccount(id string) ([]*iot.ThingAttribute, error) {
	// Search for things
	things := []*iot.ThingAttribute{}
	return c.getPaginatedThings(&iot.ListThingsInput{
		AttributeName:  aws.String(accountIdAttributeName),
		AttributeValue: aws.String(id),
	}, nil, things)
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
