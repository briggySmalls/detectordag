//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"log"
	"os"
	// mage:import
	_ "github.com/briggysmalls/detectordag/shared/mage"
)

var path string

func init() {
	var err error
	path, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateSpec() error {
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", path), "quay.io/goswagger/swagger", "generate", "spec", "-w", "/app/api", "-o", "/app/api.yml")
}

func GenerateLib() error {
	// Remove any existing content
	const libDir = "frontend/lib/client"
	err := sh.Run("rm", "-rf", libDir)
	if err != nil {
		return err
	}
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/local", path), "openapitools/openapi-generator-cli", "generate", "-i", "/local/api.yml", "-g", "typescript-axios", "-o", fmt.Sprintf("/local/%s", libDir))
}
