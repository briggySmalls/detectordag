package mage

import (
	"errors"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Format mg.Namespace

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
	// We need mocks for tests
	mg.Deps(Generate)
	return sh.RunV("go", "test", "-v", "./...")
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
