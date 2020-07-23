package iot

//go:generate mockgen -destination mock_iot.go -package iot -mock_names Client=MockIoT github.com/aws/aws-sdk-go/service/iot/iotiface IoTAPI

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetThing(t *testing.T) {
	// Define some test parameters
	const (
		deviceID   = "261f3f87-84bb-4c0e-91bc-ba41c3bc0668"
		deviceName = "Testing"
		accountID  = "9962902c-f7e7-417d-bea0-dc2eb0bc67d7"
	)
	// Create unit under test and mocks
	mock, c := createUnitAndMocks(t)
	// Expect iot calls
	mock.EXPECT().DescribeThing(gomock.Not(gomock.Nil())).Do(func(input *iot.DescribeThingInput) {
		// Assert that the device ID is as expected
		assert.Equal(t, deviceID, *input.ThingName)
	}).Return(&iot.DescribeThingOutput{
		ThingName: aws.String(deviceID),
		Attributes: map[string]*string{
			accountIDAttributeName: aws.String(accountID),
			nameAttributeName:      aws.String(deviceName),
		},
	}, nil)
	// Query for a known device
	device, err := c.GetThing(deviceID)
	assert.NoError(t, err)
	// Assert it has expected fields
	assert.Equal(t, accountID, device.AccountId)
}

func TestGetThings(t *testing.T) {
	const (
		accountID     = "9962902c-f7e7-417d-bea0-dc2eb0bc67d7"
		nextToken     = "1a13f6f2-13e5-408d-a184-1ce292320175"
		deviceOne     = "4fa62730-dd7a-421b-91b9-ec1f20ad265b"
		deviceOneName = "One"
		deviceTwo     = "70c3e40a-fbc2-40d7-9cb3-7f7637f85cb4"
		deviceTwoName = "Two"
	)
	// Create unit under test and mocks
	mock, c := createUnitAndMocks(t)
	// Expect a call to ListDevices
	gomock.InOrder(
		mock.EXPECT().ListThings(gomock.Not(gomock.Nil())).Do(func(input *iot.ListThingsInput) {
			// Assert that the search is setting the correct parameters
			assert.Equal(t, thingType, *input.ThingTypeName)
			assert.Nil(t, input.AttributeName)
			assert.Nil(t, input.AttributeValue)
		}).Return(&iot.ListThingsOutput{
			Things: []*iot.ThingAttribute{
				{
					ThingName: aws.String(deviceOne),
					Attributes: map[string]*string{
						accountIDAttributeName: aws.String(accountID),
						nameAttributeName:      aws.String(deviceOneName),
					},
				},
			},
			NextToken: aws.String(nextToken), // Indicate there are more things to come
		}, nil),
		mock.EXPECT().ListThings(gomock.Not(gomock.Nil())).Do(func(input *iot.ListThingsInput) {
			// Assert that the search is setting the correct parameters
			assert.Equal(t, thingType, *input.ThingTypeName)
			assert.Nil(t, input.AttributeName)
			assert.Nil(t, input.AttributeValue)
			assert.Equal(t, nextToken, *input.NextToken)
		}).Return(&iot.ListThingsOutput{
			Things: []*iot.ThingAttribute{
				{
					ThingName: aws.String(deviceTwo),
					Attributes: map[string]*string{
						accountIDAttributeName: aws.String(accountID),
						nameAttributeName:      aws.String(deviceTwoName),
					},
				},
			},
		}, nil),
	)
	// Query for devices
	devices, err := c.GetThings()
	assert.NoError(t, err)
	// Assert the returned devices
	expectedDevices := []Device{
		{DeviceId: deviceOne, Name: deviceOneName, AccountId: accountID},
		{DeviceId: deviceTwo, Name: deviceTwoName, AccountId: accountID},
	}
	assert.Len(t, devices, len(expectedDevices))
	for i, device := range devices {
		assert.Equal(t, expectedDevices[i], *device)
	}
}

func TestGetThingsByAccount(t *testing.T) {
	const (
		accountID     = "9962902c-f7e7-417d-bea0-dc2eb0bc67d7"
		nextToken     = "1a13f6f2-13e5-408d-a184-1ce292320175"
		deviceOne     = "4fa62730-dd7a-421b-91b9-ec1f20ad265b"
		deviceOneName = "One"
		deviceTwo     = "70c3e40a-fbc2-40d7-9cb3-7f7637f85cb4"
		deviceTwoName = "Two"
	)
	// Create unit under test and mocks
	mock, c := createUnitAndMocks(t)
	// Expect a call to ListDevices
	gomock.InOrder(
		mock.EXPECT().ListThings(gomock.Not(gomock.Nil())).Do(func(input *iot.ListThingsInput) {
			// Assert that the search is setting the correct parameters
			assert.Equal(t, accountIDAttributeName, *input.AttributeName)
			assert.Equal(t, accountID, *input.AttributeValue)
		}).Return(&iot.ListThingsOutput{
			Things: []*iot.ThingAttribute{
				{
					ThingName: aws.String(deviceOne),
					Attributes: map[string]*string{
						accountIDAttributeName: aws.String(accountID),
						nameAttributeName:      aws.String(deviceOneName),
					},
				},
			},
			NextToken: aws.String(nextToken), // Indicate there are more things to come
		}, nil),
		mock.EXPECT().ListThings(gomock.Not(gomock.Nil())).Do(func(input *iot.ListThingsInput) {
			// Assert that the search is setting the correct parameters
			assert.Equal(t, accountIDAttributeName, *input.AttributeName)
			assert.Equal(t, accountID, *input.AttributeValue)
			assert.Equal(t, nextToken, *input.NextToken)
		}).Return(&iot.ListThingsOutput{
			Things: []*iot.ThingAttribute{
				{
					ThingName: aws.String(deviceTwo),
					Attributes: map[string]*string{
						accountIDAttributeName: aws.String(accountID),
						nameAttributeName:      aws.String(deviceTwoName),
					},
				},
			},
		}, nil),
	)
	// Query for devices associated with an account
	devices, err := c.GetThingsByAccount(accountID)
	assert.NoError(t, err)
	// Assert the returned devices
	expectedDevices := []Device{
		{DeviceId: deviceOne, Name: deviceOneName, AccountId: accountID},
		{DeviceId: deviceTwo, Name: deviceTwoName, AccountId: accountID},
	}
	assert.Len(t, devices, len(expectedDevices))
	for i, device := range devices {
		assert.Equal(t, expectedDevices[i], *device)
	}
}

func TestRegisterDevice(t *testing.T) {
	const (
		accountID             = "aac45d02-c97d-442c-8431-336d578fdcf7"
		deviceID              = "f80103e1-ba55-4b55-b80e-b24f5dd518bb"
		deviceName            = "Annex"
		certificateID         = "d5c29c58-5a69-4b46-908e-13d2ad5b21a6"
		certificatePem        = "THIS IS A PEM"
		certificatePrivateKey = "THIS IS A PRIVATE KEY"
		certificatePublicKey  = "THIS IS A PUBLIC KEY"
	)
	// Create unit under test and mocks
	mock, c := createUnitAndMocks(t)
	// Configure mock to 'create' a certificate successfully
	mock.EXPECT().CreateKeysAndCertificate(gomock.Not(gomock.Nil())).Do(func(input *iot.CreateKeysAndCertificateInput) {
		assert.False(t, *input.SetAsActive)
	}).Return(&iot.CreateKeysAndCertificateOutput{
		CertificateId:  aws.String(certificateID),
		CertificatePem: aws.String(certificatePem),
		KeyPair: &iot.KeyPair{
			PublicKey:  aws.String(certificatePublicKey),
			PrivateKey: aws.String(certificatePrivateKey),
		},
	}, nil)
	// Configure mock to create a device successfully
	mock.EXPECT().RegisterThing(gomock.Not(gomock.Nil())).Do(func(input *iot.RegisterThingInput) {
		assert.Equal(t, provisioningTemplate, *input.TemplateBody)
		assert.Equal(t, thingGroup, *input.Parameters["ThingGroup"])
		assert.Equal(t, thingType, *input.Parameters["ThingType"])
		assert.Equal(t, deviceName, *input.Parameters["DeviceName"])
		assert.Equal(t, certificateID, *input.Parameters["CertificateId"])
		assert.Equal(t, accountID, *input.Parameters["AccountId"])
		assert.Equal(t, deviceID, *input.Parameters["DeviceId"])
	}).Return(nil, nil)
	// Configure mock to expect activation of certificate
	mock.EXPECT().UpdateCertificate(gomock.Not(gomock.Nil())).Do(func(input *iot.UpdateCertificateInput) {
		assert.Equal(t, certificateID, *input.CertificateId)
		assert.Equal(t, "ACTIVE", *input.NewStatus)
	}).Return(nil, nil)
	// Query for devices associated with an account
	device, certs, err := c.RegisterThing(accountID, deviceID, deviceName)
	assert.NoError(t, err)
	// Assert device has expected fields
	assert.Equal(t, deviceName, device.Name)
	assert.Equal(t, accountID, device.AccountId)
	assert.Equal(t, deviceID, device.DeviceId)
	// Assert certs has expected fields
	assert.Equal(t, certificatePem, certs.Certificate)
	assert.Equal(t, certificatePublicKey, certs.Public)
	assert.Equal(t, certificatePrivateKey, certs.Private)
}

func TestSetVisibiltyState(t *testing.T) {
	// Define some test parameters
	const (
		deviceID   = "261f3f87-84bb-4c0e-91bc-ba41c3bc0668"
		deviceName = "Testing"
		accountID  = "9962902c-f7e7-417d-bea0-dc2eb0bc67d7"
	)
	// Create unit under test and mocks
	mock, c := createUnitAndMocks(t)
	// Create a helper for asserting the attribute update request
	updateThingAssertion := func(state string) *gomock.Call {
		return mock.EXPECT().UpdateThing(gomock.Not(gomock.Nil())).Do(func(input *iot.UpdateThingInput) {
			assert.Equal(t, deviceID, *input.ThingName)
			assert.Equal(t, thingType, *input.ThingTypeName)
			assert.Equal(t, state, *input.AttributePayload.Attributes[visibilityAttributeName])
		})
	}
	// Configure mock to set the state
	gomock.InOrder(
		updateThingAssertion("true"),
		updateThingAssertion("false"),
	)
	// Create a device to pass in
	device := Device{
		DeviceId:  deviceID,
		Name:      deviceName,
		AccountId: accountID,
	}
	// Run the test
	c.SetVisibiltyState(&device, true)
	c.SetVisibiltyState(&device, false)
}

func createUnitAndMocks(t *testing.T) (*MockIoTAPI, Client) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	// Create mock IoTAPI
	mock := NewMockIoTAPI(ctrl)
	// Create the unit under test
	iot := client{
		iot: mock,
	}
	// Create the new iot client
	return mock, &iot
}
