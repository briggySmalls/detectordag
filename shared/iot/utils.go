package iot

import (
	"github.com/aws/aws-sdk-go/service/iot"
)

const (
	accountIDAttributeName = "account-id"
	nameAttributeName      = "name"
)

type thingAttribute struct {
	*iot.ThingAttribute
}

type describeThingOutput struct {
	*iot.DescribeThingOutput
}

type deviceSource interface {
	ToDevice() (*Device, error)
}

func (t *thingAttribute) ToDevice() (*Device, error) {
	// Get the name
	name, err := t.getAttribute(nameAttributeName)
	if err != nil {
		return nil, err
	}
	// Get the account ID
	accountID, err := t.getAttribute(accountIDAttributeName)
	if err != nil {
		return nil, err
	}
	// Get the device ID
	deviceID := t.ThingName
	return &Device{
		Name:      name,
		DeviceId:  *deviceID,
		AccountId: accountID,
	}, nil
}

func (t *thingAttribute) getAttribute(key string) (string, error) {
	if accountID, ok := t.ThingAttribute.Attributes[key]; ok {
		return *accountID, nil
	}
	return "", ErrAccountIDMissing
}

func (t *describeThingOutput) ToDevice() (*Device, error) {
	// Get the name
	name, err := t.getAttribute(nameAttributeName)
	if err != nil {
		return nil, err
	}
	// Get the account ID
	accountID, err := t.getAttribute(accountIDAttributeName)
	if err != nil {
		return nil, err
	}
	// Get the device ID
	deviceID := t.ThingName
	return &Device{
		Name:      name,
		DeviceId:  *deviceID,
		AccountId: accountID,
	}, nil
}

func (t *describeThingOutput) getAttribute(key string) (string, error) {
	if accountID, ok := t.DescribeThingOutput.Attributes[key]; ok {
		return *accountID, nil
	}
	return "", ErrAccountIDMissing
}
