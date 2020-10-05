package mage

import (
	"errors"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"regexp"
	"fmt"
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

// Vets
func Vet() error {
	// Get the go directories
	packages, err := nonToolsDirs()
	if err != nil {
		return err
	}
	args := append([]string{"vet"}, packages...)
	return sh.RunV("go", args...)
}

// Runs the tests
func Test() error {
	// We need mocks for tests
	mg.Deps(Generate)
	// Get the go directories
	packages, err := nonToolsDirs()
	if err != nil {
		return err
	}
	args := append([]string{"test", "-v"}, packages...)
	return sh.RunV("go", args...)
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

func nonToolsDirs() ([]string, error) {
	// Get the go directories
	output, err := sh.Output("go", "list", "./...")
	if err != nil {
		return nil, err
	}
	// Exclude tools
	re := regexp.MustCompile(fmt.Sprintf("(?m)[\r\n]+^.*%s.*$", toolsDir))
	packages := strings.Split(re.ReplaceAllString(output, ""), "\n")
	return packages, nil
}
