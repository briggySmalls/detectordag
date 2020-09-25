package api

import (
	"github.com/briggysmalls/detectordag/shared/mage"
	"github.com/magefile/mage/mg"
)

type Api mg.Namespace

var cmds = []string{"start-api"}

// Runs the lambda function locally
func (Api) Run() error {
	return mage.Sam(cmds, false)
}

// Debugs the lambda function locally
func (Api) Debug() error {
	return mage.Sam(cmds, true)
}
