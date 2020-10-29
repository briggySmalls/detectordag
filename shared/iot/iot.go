package iot

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iot/iotiface"
	"log"
	"strconv"
)

const (
	accountIDAttributeName  = "account-id"
	nameAttributeName       = "name"
	thingType               = "detectordag"
	thingGroup              = "detectordag"
	visibilityAttributeName = "visibility"
)

type client struct {
	iot iotiface.IoTAPI
}

type Client interface {
	GetThing(id string) (*Device, error)
	GetThingsByVisibility(status bool) ([]*Device, error)
	GetThingsByAccount(id string) ([]*Device, error)
	RegisterThing(accountID, deviceID, name string) (*Device, *Certificates, error)
	SetVisibiltyState(deviceID string, state bool) error
}

// Device holds the non-state properties of a device
type Device struct {
	Name       string
	DeviceId   string
	AccountId  string
	Visibility bool
}

type Certificates struct {
	Certificate string
	Public      string
	Private     string
}

// New gets a new Client
func New(sesh *session.Session) (Client, error) {
	// Create Amazon IoT client
	iot := iot.New(sesh)
	if iot == nil {
		return nil, errors.New("Failed to create database client")
	}
	// Create our client wrapper
	return &client{
		iot: iot,
	}, nil
}

// GetThing gets a thing with the given name from the AWS IoT registry
func (c *client) GetThing(id string) (*Device, error) {
	// Fetch the specified thing
	thing, err := c.iot.DescribeThing(&iot.DescribeThingInput{ThingName: aws.String(id)})
	if err != nil {
		return nil, fmt.Errorf("Get thing failure for '%s': %w", id, err)
	}
	// Convert the response to a 'Device'
	d := describeThingOutput{thing}
	return d.ToDevice()
}

// GetThings returns all things which have the specified visiblity status
func (c *client) GetThingsByVisibility(status bool) ([]*Device, error) {
	return c.getPaginatedDevices(&iot.ListThingsInput{
		AttributeName:  aws.String(visibilityAttributeName),
		AttributeValue: aws.String(strconv.FormatBool(status)),
	})
}

// GetThingsByAccount returns all things which are associated with the specified accountg
func (c *client) GetThingsByAccount(id string) ([]*Device, error) {
	// Search for things
	return c.getPaginatedDevices(&iot.ListThingsInput{
		AttributeName:  aws.String(accountIDAttributeName),
		AttributeValue: aws.String(id),
	})
}

// RegisterThing creates a new thing and provides certificates for it to communicate
func (c *client) RegisterThing(accountID, deviceID, name string) (*Device, *Certificates, error) {
	// Create a new certificate
	certsResponse, err := c.createCertificate()
	if err != nil {
		return nil, nil, err
	}
	// Create a new thing
	_, err = c.registerThing(deviceID, *certsResponse.CertificateId, name, accountID)
	// Check if we failed to create the thing
	if err != nil {
		log.Printf("Failed to RegisterThing: %v", err)
		// Try our best to delete the certificate we created
		_, delErr := c.iot.DeleteCertificate(&iot.DeleteCertificateInput{
			CertificateId: certsResponse.CertificateId,
			ForceDelete:   aws.Bool(true),
		})
		if delErr != nil {
			log.Printf("Failed to delete certificate: %v", delErr)
		}
		return nil, nil, err
	}
	// We're all done!
	d := Device{
		DeviceId:  deviceID,
		Name:      name,
		AccountId: accountID,
	}
	certs := Certificates{
		Certificate: *certsResponse.CertificatePem,
		Public:      *certsResponse.KeyPair.PublicKey,
		Private:     *certsResponse.KeyPair.PrivateKey,
	}
	// Activate certificate now we're happy all is well
	_, err = c.iot.UpdateCertificate(&iot.UpdateCertificateInput{
		CertificateId: certsResponse.CertificateId,
		NewStatus:     aws.String("ACTIVE"),
	})
	return &d, &certs, nil
}

// SetVisibilityState sets attribute indicating if the device is lost
func (c *client) SetVisibiltyState(deviceID string, state bool) error {
	// Set the attribute
	_, err := c.iot.UpdateThing(&iot.UpdateThingInput{
		ThingName:     aws.String(deviceID),
		ThingTypeName: aws.String(thingType),
		AttributePayload: &iot.AttributePayload{
			Attributes: map[string]*string{
				visibilityAttributeName: aws.String(strconv.FormatBool(state)),
			},
			Merge: aws.Bool(true), // Don't nuke the other attributes
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to update thing '%s': %w", deviceID, err)
	}
	return err
}

// createCertificate creates a new certificate
func (c *client) createCertificate() (*iot.CreateKeysAndCertificateOutput, error) {
	return c.iot.CreateKeysAndCertificate(&iot.CreateKeysAndCertificateInput{
		SetAsActive: aws.Bool(false),
	})
}

func (c *client) registerThing(deviceId, certificateID, name, accountID string) (*iot.RegisterThingOutput, error) {
	return c.iot.RegisterThing(&iot.RegisterThingInput{
		// Use the template for provisioning a new device
		TemplateBody: aws.String(provisioningTemplate),
		// Set the parameters used in the template
		Parameters: map[string]*string{
			"DeviceId":      aws.String(deviceId),
			"ThingGroup":    aws.String(thingGroup),
			"ThingType":     aws.String(thingType),
			"DeviceName":    aws.String(name),
			"AccountId":     aws.String(accountID),
			"CertificateId": aws.String(certificateID),
		},
	})
}

func (c *client) getPaginatedDevices(input *iot.ListThingsInput) ([]*Device, error) {
	// Search for things
	things := []*iot.ThingAttribute{}
	var err error
	things, err = c.getPaginatedThings(input, nil, things)
	if err != nil {
		return nil, err
	}
	// Wrap up each thing
	wrappedThings := make([]*Device, len(things))
	for i, thing := range things {
		t := thingAttribute{thing}
		device, err := t.ToDevice()
		if err != nil {
			return nil, err
		}
		wrappedThings[i] = device
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
		return nil, fmt.Errorf("Failed to list things: %w", err)
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
