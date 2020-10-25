package iot

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/iot"
	"strconv"
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
	// Get the visibility state
	visibilityStr, err := t.getAttribute(visibilityAttributeName)
	if err != nil {
		return nil, err
	}
	// Convert visibility state to a boolean
	visibility, err := strconv.ParseBool(visibilityStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Visibility attribute to bool: '%s'", visibilityStr)
	}
	// Get the device ID
	deviceID := t.ThingName
	return &Device{
		Name:       name,
		DeviceId:   *deviceID,
		AccountId:  accountID,
		Visibility: visibility,
	}, nil
}

func (t *thingAttribute) getAttribute(key string) (string, error) {
	if accountID, ok := t.ThingAttribute.Attributes[key]; ok {
		return *accountID, nil
	}
	return "", fmt.Errorf("Attribute %s was missing", key)
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
	// Get the visibility state
	visibilityStr, err := t.getAttribute(visibilityAttributeName)
	if err != nil {
		return nil, err
	}
	// Convert visibility state to a boolean
	visibility, err := strconv.ParseBool(visibilityStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Visibility attribute to bool: '%s'", visibilityStr)
	}
	// Get the device ID
	deviceID := t.ThingName
	return &Device{
		Name:       name,
		DeviceId:   *deviceID,
		AccountId:  accountID,
		Visibility: visibility,
	}, nil
}

func (t *describeThingOutput) getAttribute(key string) (string, error) {
	if accountID, ok := t.DescribeThingOutput.Attributes[key]; ok {
		return *accountID, nil
	}
	return "", fmt.Errorf("Attribute '%s' was missing", key)
}
