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
	return helper.Invoke(false, "")
}

// Invokes the lambda function locally, running the debug server
func (Invoke) Debug() error {
	return helper.Invoke(true, "")
}

// InstallTools installs project tools
func InstallTools() error {
	return helper.InstallTools()
}
