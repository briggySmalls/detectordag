package lambda

import (
	"github.com/briggysmalls/detectordag/shared/mage"
	"github.com/magefile/mage/mg"
)

const (
	cmd = "invoke"
)

type Lambda mg.Namespace

var LambdaName string
var cmds []string{}

func init() {
	cmds = []string{"invoke", LambdaName}
}

// Runs the lambda function locally
func (Invoke) Run() error {
	return mage.Sam(cmds, false)
}

// Debugs the lambda function locally
func (Invoke) Debug() error {
	return mage.Sam(cmds, true)
}
