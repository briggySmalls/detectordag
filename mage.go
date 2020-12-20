//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/mg"
	"log"
	"os"
	// mage:import
	_ "github.com/briggysmalls/detectordag/shared/mage"
)

const (
	// Directory for build outputs
	buildDir = "./build"
	applicationName = "detectordag-edge"
	imgConfigFile = "./provisioning/detectordag-edge.json"
	balenaVersion = "v2.54.2+rev1"
	deviceType = "raspberrypi"
)

var imageFile = fmt.Sprintf("%s/detectordag-edge.img", buildDir)

type Generate mg.Namespace

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

func MockApi() error {
	return sh.Run("docker", "run", "--init", "--rm", "-v", fmt.Sprintf("%s:/local", path), "-p", "3000:4010", "stoplight/prism:4", "mock", "-h", "0.0.0.0", "/local/api.yml")
}

// Download the BalenaOS image for an edge device
func DownloadImg() error {
	// Ensure we have a build directory
	mg.Deps(createBuildDir)
	// Download the OS image
	return sh.Run("balena", "os", "download", deviceType, "--version", balenaVersion, "--output", imageFile)
}

// Configure the BalenaOS image for use in the detectordag application
func ConfigureImg() error {
	// Download the image
	mg.Deps(DownloadImg)
	// Apply the application configuration to it
	return sh.Run("balena", "os", "configure", "--application", applicationName, "--config", imgConfigFile, imageFile)
}
