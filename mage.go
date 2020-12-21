//+build mage

package main

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	// mage:import
	_ "github.com/briggysmalls/detectordag/shared/mage"
)

// Build constants
const (
	// Directory for build outputs
	buildDir                = "./build"
	applicationName         = "detectordag-edge"
	awsProvisioningTemplate = "./config/thing.json"
	balenaVersion           = "v2.54.2+rev1"
	deviceType              = "raspberrypi"
	certFile                = "thing.cert.pem"
	keyFile                 = "thing.private.key"
	deviceIDEnvVar          = "DDAG_DEVICE_ID"
	imageFile               = "detectordag-edge.img"
)

// Build nearly-constants (derived from constants)
var vanillaImageFile = fmt.Sprintf("%s/detectordag-edge.img", buildDir)

type Generate mg.Namespace
type Provision mg.Namespace

var path string

func init() {
	var err error
	path, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

func createBuildDir() error {
	// Create a new one
	return sh.Run("mkdir", "-p", buildDir)
}

func deviceBuildDir() error {
	mg.Deps(createBuildDir)
	// Get configuration argument
	deviceID, err := getEnvVar(deviceIDEnvVar)
	if err != nil {
		return err
	}
	// Ensure build directory is present
	return sh.Run("mkdir", "-p", fmt.Sprintf("%s/%s", buildDir, deviceID))
}

// Generates the OpenAPI specification from the api
func (Generate) Spec() error {
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", path), "quay.io/goswagger/swagger", "generate", "spec", "-w", "/app/api", "-o", "/app/api.yml")
}

func ValidateSpec() error {
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/local", path), "openapitools/openapi-generator-cli", "validate", "-i", "/local/api.yml")
}

// Generates the javascript API client from the OpenAPI specification
func (Generate) Lib() error {
	// Remove any existing content
	const libDir = "frontend/lib/client"
	err := sh.Run("rm", "-rf", libDir)
	if err != nil {
		return err
	}
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/local", path), "openapitools/openapi-generator-cli", "generate", "-i", "/local/api.yml", "-g", "typescript-axios", "-o", fmt.Sprintf("/local/%s", libDir))
}

// Generates documentation from the OpenAPI specification
func (Generate) Docs() error {
	// Remove any existing content
	const docsDir = "build/docs"
	err := sh.Run("rm", "-rf", docsDir)
	if err != nil {
		return err
	}
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/local", path), "broothie/redoc-cli", "bundle", "/local/api.yml", "-o", fmt.Sprintf("/local/%s/index.html", docsDir))
}

// Start a mock API from the OpenAPI specification
func MockApi() error {
	return sh.Run("docker", "run", "--init", "--rm", "-v", fmt.Sprintf("%s:/local", path), "-p", "3000:4010", "stoplight/prism:4", "mock", "-h", "0.0.0.0", "/local/api.yml")
}

// Download the BalenaOS image for an edge device
func DownloadImg() error {
	// Ensure we have a build directory
	mg.Deps(createBuildDir)
	// Download the OS image
	return sh.Run("balena", "os", "download", deviceType, "--version", balenaVersion, "--output", vanillaImageFile)
}

// Configure the BalenaOS image for use in the detectordag application
func ConfigureImg() error {
	mg.Deps(
		deviceBuildDir, // We need a build dir to output into
		RegisterBalena, // We need the device to be registered
	)
	var err error
	// Ensure the vanilla image is present
	if _, err = os.Stat(vanillaImageFile); os.IsNotExist(err) {
		mg.Deps(DownloadImg)
	}
	// Get configuration argument
	deviceID, err := getEnvVar(deviceIDEnvVar)
	if err != nil {
		return err
	}
	// Copy it in preparation for a new device
	err = sh.Copy(getDeviceImageFile(deviceID), vanillaImageFile)
	if err != nil {
		return err
	}
	// Create application configuration
	configFile := fmt.Sprintf("%s/config.json", getDeviceBuildDir(deviceID))
	err = sh.Run("balena", "config", "generate",
		"--version", balenaVersion,
		"--device", deviceID,
		"--network", "ethernet",
		"--appUpdatePollInterval", "10",
		"--output", configFile,
	)
	if err != nil {
		return err
	}
	// Apply the application configuration to it
	return sh.Run("balena", "os", "configure", "--config-network", "ethernet", "--config", configFile, "--device", deviceID, getDeviceImageFile(deviceID))
}

// Register a new 'thing' on AWS
func RegisterAWS() error {
	mg.Deps(deviceBuildDir)
	// Get some configuration
	name, err := getEnvVar("DDAG_DEVICE_NAME")
	if err != nil {
		return err
	}
	deviceID, err := getEnvVar(deviceIDEnvVar)
	if err != nil {
		return err
	}
	accountID, err := getEnvVar("DDAG_ACCOUNT_ID")
	if err != nil {
		return err
	}
	// Create the IoT client
	sesh := shared.CreateSession(aws.Config{})
	client, err := iot.New(sesh)
	if err != nil {
		return err
	}
	// Register the new device
	_, certificates, err := client.RegisterThing(accountID, deviceID, name)
	if err != nil {
		return err
	}
	// Write the certificates to the build directory
	if err := writeFile(fmt.Sprintf("%s/%s/%s", buildDir, deviceID, certFile), certificates.Certificate); err != nil {
		return err
	}
	if err := writeFile(fmt.Sprintf("%s/%s/%s", buildDir, deviceID, keyFile), certificates.Private); err != nil {
		return err
	}
	return nil
}

// Register a device on BalenaCloud
func RegisterBalena() error {
	deviceID, err := getEnvVar(deviceIDEnvVar)
	if err != nil {
		return err
	}
	if err := sh.Run("balena", "device", "register", applicationName, "--uuid", deviceID); err != nil {
		return err
	}
	// Sleep a bit to make sure the device is available
	time.Sleep(100 * time.Millisecond)
	return nil
}

func SetEnvVars() error {
	mg.Deps(
		RegisterBalena, // We need the device to exist
		RegisterAWS,    // We need the certificates
	)
	// We'll need the device ID
	deviceID, err := getEnvVar(deviceIDEnvVar)
	if err != nil {
		return err
	}
	certFileText, err := readCertFile(deviceID, certFile)
	if err != nil {
		return err
	}
	keyFileText, err := readCertFile(deviceID, keyFile)
	if err != nil {
		return err
	}
	// Set certificate environment variables
	envVars := map[string]string{
		"AWS_THING_CERT": certFileText,
		"AWS_THING_KEY":  keyFileText,
		"AWS_THING_NAME": deviceID,
	}
	for key, value := range envVars {
		// Add the environment variable
		err = sh.Run("balena", "env", "add", "--device", deviceID, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Provision a device on Balena, AWS, and create an image to burn
func ProvisionDevice() error {
	mg.Deps(
		RegisterBalena, // Create the balena device
		ConfigureImg,   // Configure an image using the device
		RegisterAWS,    // Register an AWS device with the same name
		SetEnvVars,     // Set the environment variables of the device
	)
	return nil
}

// Write the BalenaOS image to an external drive
func WriteImage() error {
	// Assume caller has set drive
	drv, err := getEnvVar("DDAG_DRIVE")
	if err != nil {
		return err
	}
	return sh.Run("balena", "os", "initialize", imageFile, "--type", deviceType, "--drive", drv, "--yes")
}

func getEnvVar(varName string) (string, error) {
	env := os.Getenv(varName)
	if env == "" {
		return "", fmt.Errorf("%s not set", varName)
	}
	return env, nil
}

func GenerateDeviceID() error {
	return sh.Run("openssl", "rand", "-hex", "16")
}

func getDeviceBuildDir(deviceID string) string {
	return fmt.Sprintf("%s/%s", buildDir, deviceID)
}

func readCertFile(deviceID, file string) (string, error) {
	// Build the filename
	fileName := fmt.Sprintf("%s/%s", getDeviceBuildDir(deviceID), file)
	// Read the file
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	// Convert to base64
	return string(b64.StdEncoding.EncodeToString(dat)), nil
}

func getDeviceImageFile(deviceID string) string {
	return fmt.Sprintf("%s/%s", getDeviceBuildDir(deviceID), imageFile)
}

func writeFile(file, content string) error {
	// Create a new file
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	// Write all the content to it
	_, err2 := f.WriteString(content)
	if err2 != nil {
		return err
	}
	return nil
}
