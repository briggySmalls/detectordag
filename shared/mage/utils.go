package mage

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	toolsFile = "./tools/tools.go"
)

// Installs the module tools
func InstallTools() error {
	// Read the tools file
	toolsFile, err := ioutil.ReadFile(toolsFile)
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

// Builds the delve debugger
func BuildDelve() error {
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
	src := matches[0]
	// State our inputs and outputs
	tgt := getBinFile("delve/dlv")
	isNew, err := target.Path(tgt, src)
	if err != nil {
		return err
	}
	if isNew {
		// Build our version of delve
		return sh.RunWith(
			map[string]string{
				"GOARCH": "amd64",
				"GOOS":   "linux",
			},
			"go", "build",
			"-o", tgt,
			src)
	}
	// Nothing to build
	return nil
}

func getBinFile(file string) string {
	return path.Join(buildDir, file)
}
