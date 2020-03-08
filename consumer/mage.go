//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"regexp"
)

type Invoke mg.Namespace
type Build mg.Namespace

var invokeCmd = []string{"sam", "local", "invoke", "consumer", "-e", "event.json"}

// Invokes the lambda function locally
func (Invoke) Invoke() error {
	mg.Deps(Build.Production)
	return sh.Run(invokeCmd[0], invokeCmd[1:]...)
}

// Invokes the lambda function locally, running the debug server
func (Invoke) Debug() error {
	mg.Deps(Build.Debug, Delve)
	cmdWithDebugger := append(invokeCmd, "-d", "5986", "--debugger-path", "delve", "--debug-args", "-delveAPI=2")
	return sh.Run(
		cmdWithDebugger[0], cmdWithDebugger[1:]...,
	)
}

// Runs dep ensure and then installs the binary.
func (Build) Production() error {
	return build()
}

// Builds a debug version of the build (with debugging)
func (Build) Debug() error {
	return build("-gcflags", "all=-N -l")
}

func Delve() error {
	return sh.Run("env", "GO111MODULE=off", "GOARCH=amd64", "GOOS=linux",
		"go", "build",
		"-o", "./delve/dlv",
		"$GOPATH/src/github.com/go-delve/delve/cmd/dlv")
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

func build(extraArgs ...string) error {
	combined := []string{"GOARCH=amd64", "GOOS=linux", "go", "build"}
	combined = append(combined, extraArgs...)
	combined = append(combined, "-o", "consumer", "main.go")
	return sh.Run("env", combined...)
}
