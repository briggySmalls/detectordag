//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

var extraFlags string

// Runs dep ensure and then installs the binary.
func (Build) Production() {
	extraFlags = ""
	mg.Deps(build)
}

// Builds a debug version of the build (with debugging)
func (Build) Debug() {
	extraFlags = "-gcflags='-N -l'"
	mg.Deps(build)
}

func build() error {
	return sh.Run(
		"env", "GOARCH=amd64", "GOOS=linux",
		"go", "build",
		extraFlags,
		"-o", "consumer",
		"main.go")
}
