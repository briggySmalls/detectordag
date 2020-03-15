//+build mage

package main

import (
	"github.com/briggysmalls/detectordag/shared"
)

var helper shared.Lambda

func init() {
	helper = shared.New(".aws-sam/build/", "./tools/tools.go")
}

// Invokes the lambda function locally
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
