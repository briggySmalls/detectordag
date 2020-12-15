package shared

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"log"
)

var Validate *validator.Validate

func init() {
	// Create a global validator
	Validate = validator.New()
}

func CreateSession(config aws.Config) *session.Session {
	// Create a new session just for emailing (we have to use a different region)
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            config,
	})
	if err != nil {
		log.Fatal(err)
	}
	return sesh
}

func LogErrorAndReturn(err error) error {
	err = wrapError(err)
	log.Printf("%+v\n", err)
	return err
}

func LogErrorAndExit(err error) {
	log.Fatalf("%+v\n", wrapError(err))
}

func wrapError(err error) error {
	return errors.Wrap(err, "detectordag error:")
}
