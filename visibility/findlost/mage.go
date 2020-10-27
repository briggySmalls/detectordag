//+build mage

package main

import (
	// mage:import
	"github.com/briggysmalls/detectordag/shared/mage/lambda"
)

func init() {
	lambda.LambdaName = "findlost"
}
