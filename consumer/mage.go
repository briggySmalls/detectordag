//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

var extraFlags string

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
		"-gcflags", "\"all=-N -l\"",
		"-o", "consumer",
		"main.go")
}

func Delve() error {
	return sh.Run("env", "GO111MODULE=off", "go", "build", "-o", "./delve/dlv", "$GOPATH/src/github.com/go-delve/delve/cmd/dlv")
}
