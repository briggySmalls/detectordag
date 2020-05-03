package iot

//go:generate mockgen -destination mock_iot.go -package iot -mock_names Client=MockIoT github.com/aws/aws-sdk-go/service/iot/iotiface IoTAPI

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
	var deviceID string
	mock.EXPECT().RegisterThing(gomock.Not(gomock.Nil())).Do(func(input *iot.RegisterThingInput) {
		assert.Equal(t, provisioningTemplate, *input.TemplateBody)
		assert.Equal(t, thingGroup, *input.Parameters["ThingGroup"])
		assert.Equal(t, thingType, *input.Parameters["ThingType"])
		assert.Equal(t, deviceName, *input.Parameters["DeviceName"])
		assert.Equal(t, certificateID, *input.Parameters["CertificateId"])
		assert.Equal(t, accountID, *input.Parameters["AccountId"])
		// Check the UUID was generated
		deviceID = *input.Parameters["DeviceId"]
		_, err := uuid.Parse(deviceID)
		assert.NoError(t, err)
	}).Return(nil, nil)
	// Configure mock to expect activation of certificate
	mock.EXPECT().UpdateCertificate(gomock.Not(gomock.Nil())).Do(func(input *iot.UpdateCertificateInput) {
		assert.Equal(t, certificateID, *input.CertificateId)
		assert.Equal(t, "ACTIVE", *input.NewStatus)
	}).Return(nil, nil)
	// Query for devices associated with an account
	device, certs, err := c.RegisterThing(accountID, deviceName)
	assert.NoError(t, err)
	// Assert device has expected fields
	assert.Equal(t, deviceName, device.Name)
	assert.Equal(t, accountID, device.AccountId)
	_, err = uuid.Parse(device.DeviceId)
	assert.NoError(t, err)
	// Assert certs has expected fields
	assert.Equal(t, certificatePem, certs.Certificate)
	assert.Equal(t, certificatePublicKey, certs.Public)
	assert.Equal(t, certificatePrivateKey, certs.Private)
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
