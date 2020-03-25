//+build mage

package main

import (
	"github.com/briggysmalls/detectordag/shared/mage"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var helper mage.Lambda

const (
	envFile  = ""
	toolsDir = "./tools/tools.go"
	buildDir = ".aws-sam/"
)

func init() {
	helper = mage.New(buildDir, toolsDir)
}

// Starts the API locally
func StartApi() error {
	mg.Deps(Generate)
	return helper.StartApi(false, envFile)
}

// Starts the API locally, with debugging
func DebugApi() error {
	mg.Deps(Generate)
	return helper.StartApi(true, envFile)
}

// Build the project
func Build() error {
	mg.Deps(Generate)
	return helper.Build()
}

// InstallTools installs tools locally
func InstallTools() error {
	return helper.InstallTools()
}

// Test runs unit tests
func Test() error {
	// Generate mocks
	mg.Deps(Generate)
	return sh.Run("go", "test", "-v", "./swagger/...")
}

// Generate generates sources
func Generate() error {
	return sh.Run("go", "generate")
}
