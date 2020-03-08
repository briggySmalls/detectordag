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

// Invokes the lambda function locally
func (Invoke) Production() error {
	mg.Deps(Build.Production)
	return invoke()
}

// Invokes the lambda function locally, running the debug server
func (Invoke) Debug() error {
	mg.Deps(Build.Debug, Delve)
	return invoke("-d", "5986", "--debugger-path", "delve", "--debug-args", "-delveAPI=2")
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

func invoke(extraArgs ...string) error {
	combined := []string{"local", "invoke", "consumer", "-e", "event.json"}
	combined = append(combined, extraArgs...)
	return sh.Run("sam", combined...)
}
