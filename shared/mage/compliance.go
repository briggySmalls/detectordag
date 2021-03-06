package mage

import (
	"errors"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"regexp"
	"strings"
)

type Format mg.Namespace

// Directory to ignore because it only contains tools
const toolsDir = "tools"

// Fixes formatting
func (Format) Fix() error {
	return format("-w")
}

// Checks that formatting is valid
func (Format) Check() error {
	return format("-d")
}

// Runs the tests
func Test() error {
	// Get the go directories
	output, err := sh.Output("go", "list", "./...")
	if err != nil {
		return err
	}
	// Exclude tools
	re := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*%s.*$", toolsDir))
	packages := strings.Split(re.ReplaceAllString(output, ""), "\n")
	args := append([]string{"test", "-v"}, packages...)
	return sh.RunV("go", args...)
}

// Cleans
func Clean() error {
	return sh.RunV("find", ".", "-name", "'mock_*.go'", "-delete")
}

// Generates the mocks
func Generate() error {
	return sh.RunV("go", "generate", "./...")
}

func format(args ...string) error {
	// Combine the command and the arguments
	combinedArgs := append([]string{"-l", "-e", "-s"}, args...)
	// Add the directory
	combinedArgs = append(combinedArgs, ".")
	// Run the command
	output, err := sh.Output("gofmt", combinedArgs...)
	if err != nil {
		return err
	}
	// Check if there were any changes
	if output != "" {
		return errors.New(output)
	}
	return nil
}
