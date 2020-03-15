//+build mage

package main

import (
	"github.com/briggysmalls/detectordag/shared"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var helper shared.Lambda

const apiSpecFile = "api.yaml"

func init() {
	helper = shared.New(".aws-sam/build/", "./tools/tools.go")
}

type Swagger mg.Namespace

// Starts the API locally
func StartApi() error {
	return helper.StartApi()
}

// Build the project
func Build() error {
	return helper.Build()
}

func Delve() error {
	return helper.BuildDelve()
}

func InstallTools() error {
	return helper.InstallTools()
}

func (Swagger) Validate() error {
	return sh.Run("swagger", "validate", apiSpecFile)
}

func (Swagger) Generate() error {
	return sh.Run("swagger-codegen", "generate", "-i", "api.yaml", "--lang", "go-server")
}
