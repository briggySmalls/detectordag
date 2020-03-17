//+build mage

package main

import (
	"github.com/briggysmalls/detectordag/shared/mage"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var helper mage.Lambda

const apiSpecFile = "api.yaml"

func init() {
	helper = mage.New(".aws-sam/build/", "./tools/tools.go")
}

// Starts the API locally
func StartApi() error {
	mg.Deps(Build)
	return helper.StartApi()
}

// Build the project
func Build() error {
	mg.Deps(Generate)
	return helper.Build()
}

func Delve() error {
	return helper.BuildDelve()
}

func InstallTools() error {
	return helper.InstallTools()
}

func Generate() error {
	return sh.Run("swagger-codegen", "generate", "-i", "api.yaml", "--lang", "go-server")
}

func Test() error {
	mg.Deps(Mock)
	return sh.Run("go", "test", "./...")
}

func Mock() error {
	target.Path("./mocks/db.go", "../shared/database/db.go")
	return sh.Run("mockgen", "-destination", "mocks/db.go", "-package", "mocks", "github.com/briggysmalls/detectordag/shared/database", "Client")
}
