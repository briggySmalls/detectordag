//+build mage

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	policyName        = "dd-edge-policy"
	belanaAppName     = "detectordag-edge"
	accountsTableName = "accounts"
)

type createThingResponse struct {
	ThingName string `json:""`
	ThingArn  string `json:""`
	ThingId   string `json:""`
}

type keyPair struct {
	Public  string `json:"PublicKey"`
	Private string `json:"PrivateKey"`
}

type createCertificateResponse struct {
	Arn     string  `json:"certificateArn"`
	Id      string  `json:"certificateId"`
	Pem     string  `json:"certificatePem"`
	KeyPair keyPair `json:""`
}

type endpointDescriptionResponse struct {
	Address string `json:"endpointAddress"`
}

func CreateThing() error {
	// Create a new cert/key
	createCertificateResponse, err := createCertificate()
	if err != nil {
		return err
	}
	// Decide on a new device ID
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	// Create a new thing
	err = createThing(id, createCertificateResponse.Id)
	if err != nil {
		return err
	}
	// Add the thing to the database
	err = createDbEntry(id)
	// Create balena device
	err = createDevice(id)
	if err != nil {
		return err
	}
	// Set certificates
	err = setCertificates(id, createCertificateResponse.Pem, createCertificateResponse.KeyPair.Private)
	if err != nil {
		return err
	}
	return nil
}

// CreatePolicy creates a policy for the edge devices
func CreatePolicy() error {
	return sh.Run("aws", "iot", "create-policy",
		"--policy-name", policyName,
		"--policy-document", "file://config/policy.json")
}

// CreateRule creates a rule to fire a lambda function
func CreateRule() error {
	return sh.Run("aws", "iot", "create-topic-rule",
		"--rule-name", "power_status_changed",
		"--topic-rule-payload", "file://config/topicRule.json")
}

func CreateTables() error {
	// Create accounts table
	return sh.Run("aws", "dynamodb", "create-table", "--table-name", accountsTableName, "--cli-json-input", "file://db/accounts.json")
}

// CreateThing creates a new certificate
func createCertificate() (*createCertificateResponse, error) {
	// Create a new certificate
	output, err := sh.Output("aws", "iot", "create-keys-and-certificate")
	if err != nil {
		return nil, err
	}
	// Parse the JSON
	var response createCertificateResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return nil, err
	}
	log.Printf("Certificate created: %s", response.Id)
	return &response, nil
}

func createDevice(id string) error {
	err := sh.Run("balena", "device", "register", belanaAppName, "--uuid", id)
	if err != nil {
		return err
	}
	log.Printf("Created device %s", id)
	return nil
}

// Encode a string in base64
func encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func setCertificates(id, cert, key string) error {
	// Convert the certificates to base64
	envs := map[string]string{
		"AWS_THING_CERT": encode(cert),
		"AWS_THING_KEY":  encode(key),
	}
	// Set the variables
	for key, value := range envs {
		err := sh.Run("balena", "env", "add", "--device", id, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDbEntry(id string) error {
	return sh.Run("aws", "dynamodb", "put-item",
		"--table-name", "devices",
		"--item", fmt.Sprintf("{\"device-id\": {\"S\": \"%s\"}}", id))
}

// createThing makes a new thing in AWS
func createThing(thingName, certificateId string) error {
	// Create a new thing
	output, err := sh.Output("aws", "iot", "register-thing",
		"--template-body", "file://config/thing.json",
		"--parameters", fmt.Sprintf(
			"ThingName=%s,CertificateId=%s,PolicyName=%s",
			thingName, certificateId, policyName))
	if err != nil {
		return err
	}
	// Parse the JSON
	var response createThingResponse
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return err
	}
	log.Printf("Thing created: %s", response.ThingName)
	return nil
}
