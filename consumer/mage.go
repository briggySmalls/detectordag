//+build mage

package main

import (
	"github.com/briggysmalls/detectordag/shared/mage"
	"github.com/magefile/mage/mg"
)

var helper mage.Lambda

func init() {
	helper = mage.New(".aws-sam/build/", "./tools/tools.go")
}

type Invoke mg.Namespace

// Invokes the lambda function locally
func (Invoke) Production() error {
	return helper.Invoke()
}

// Invokes the lambda function locally, running the debug server
func (Invoke) Debug() error {
	return helper.InvokeDebug()
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
