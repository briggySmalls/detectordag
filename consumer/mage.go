//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const binDir = "bin/"

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
	return invoke("-d", "5986", "--debugger-path", getBinFile("delve"), "--debug-args", "-delveAPI=2")
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
	// Ensure delve is installed
	mg.Deps(InstallTools)
	// Find Delve
	pattern := path.Join(
		os.ExpandEnv("$GOPATH"),
		"pkg/mod/github.com/go-delve/delve@*/cmd/dlv")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return fmt.Errorf("delve not found")
	}
	// Build our version of delve
	return sh.RunWith(
		map[string]string{
			"GOARCH": "amd64",
			"GOOS":   "linux",
		},
		"go", "build",
		"-o", getBinFile("delve/dlv"),
		matches[0])
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
	combined := []string{"build"}
	combined = append(combined, extraArgs...)
	combined = append(combined, "-o", getBinFile("consumer"), "main.go")
	return sh.RunWith(
		map[string]string{
			"GOARCH": "amd64",
			"GOOS":   "linux",
		},
		"go", combined...)
}

func invoke(extraArgs ...string) error {
	combined := []string{"local", "invoke", "consumer", "-e", "event.json"}
	combined = append(combined, extraArgs...)
	return sh.Run("sam", combined...)
}

func getBinFile(file string) string {
	return path.Join(binDir, file)
}
