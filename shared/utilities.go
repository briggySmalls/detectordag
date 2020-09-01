package shared

import (
	"github.com/pkg/errors"
	"log"
)

func LogErrorAndReturn(err error) {
	log.Printf("%+v\n", wrapError(err))
}

func LogErrorAndExit(err error) {
	log.Fatalf("%+v\n", wrapError(err))
}

func wrapError(err error) error {
	return errors.Wrap(err, "detectordag error:")
}
