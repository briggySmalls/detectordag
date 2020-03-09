//+build mage

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/magefile/mage/sh"
	"log"
	"strings"
)

const policyName = "dd-edge-policy"
const belanaAppName = "detectordag-edge"

type createThingResponse struct {
	ThingName string `json:""`
	ThingArn  string `json:""`
	ThingId   string `json:""`
}

type keyPair struct {
	Publickey  string `json:""`
	Privatekey string `json:""`
}

type createCertificateResponse struct {
	Arn  string  `json:"certificateArn"`
	Id   string  `json:"certificateId"`
	Pem  string  `json:"certificatePem"`
	Pair keyPair `json:"certificatePair"`
}

func CreateThing() error {
	// Create a new cert/key
	createCertificateResponse, err := createCertificate()
	if err != nil {
		return err
	}
	// Create a new thing
	err = createThing("dd-edge-1", createCertificateResponse.Id)
	if err != nil {
		return err
	}
	// Create balena device
	id, err := createDevice()
	if err != nil {
		return err
	}
	// Set certificates
	err = setCertificates(*id, createCertificateResponse.Pem, createCertificateResponse.Pair.Privatekey)
	if err != nil {
		return err
	}
	return nil
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

func createDevice() (*string, error) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	err := sh.Run("balena", "device", "register", belanaAppName, "--uuid", id)
	if err != nil {
		return nil, err
	}
	log.Printf("Created device %s", id)
	return &id, nil
}

// createThing makes a new thing in AWS
func createThing(thingName, certificateId string) error {
	// Create a new thing
	output, err := sh.Output("aws", "iot", "register-thing",
		"--template-body", "file://thing.json",
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

func setCertificates(id, cert, key string) error {
	// Convert the certificates to base64
	envs := map[string]string{
		"AWS_THING_CERT": encode(cert),
		"AWS_THING_KEY":  encode(key),
	}
	// Set the variables
	for key, value := range envs {
		err := sh.Run("balena", "env", "add", "--device", id, key, value)
		log.Printf("Set env var: %s = %s", key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreatePolicy creates a policy for the edge devices
func CreatePolicy() error {
	return sh.Run("aws", "iot", "create-policy",
		"--policy-name", policyName,
		"--policy-document", "file://policy.json")
}

// Encode a string in base64
func encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
