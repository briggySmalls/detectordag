//+build mage

package main

import (
	// mage:import
	_ "github.com/briggysmalls/detectordag/shared/mage"
	// mage:import
	_ "github.com/briggysmalls/detectordag/shared/mage/api"
	"github.com/magefile/mage/sh"
)

func Clean() error {
	return sh.RunV("find", ".", "-name", "mock_*.go", "-delete")
}
