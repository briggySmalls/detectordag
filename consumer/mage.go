//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"regexp"
)

type Build mg.Namespace

// Runs dep ensure and then installs the binary.
func (Build) Production() error {
	return sh.Run(
		"go", "build",
		"-o", "consumer",
		"main.go")
}

// Builds a debug version of the build (with debugging)
func (Build) Debug() error {
	return sh.Run(
		"go", "build",
		"-gcflags", "all=-N -l",
		"-o", "consumer",
		"main.go")
}

func Delve() error {
	return sh.Run("env", "GO111MODULE=off", "go", "build", "-o", "./delve/dlv", "$GOPATH/src/github.com/go-delve/delve/cmd/dlv")
}

func InstallTools() error {
	// Read the tools file
	toolsFile, err := ioutil.ReadFile("./tools/tools.go")
	if err != nil {
		return err
	}

	// Parse the tools
	re := regexp.MustCompile(`_ \"(.*?)\"`)
	matches := re.FindAllSubmatch(toolsFile, -1)

	// Install
	for _, match := range matches {
		err := sh.Run("go", "install", string(match[1]))
		if err != nil {
			return err
		}
	}
	return nil
}
