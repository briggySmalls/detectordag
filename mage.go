//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"log"
	"os"
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
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", path), "quay.io/goswagger/swagger", "generate", "spec", "-w", "/app/api", "-o", "api.yml")
}

func GenerateLib() error {
	return sh.Run("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/local", path), "openapitools/openapi-generator-cli", "generate", "-i", "/local/api/api.yml", "-g", "typescript", "-o", "/local/out/js")
}
