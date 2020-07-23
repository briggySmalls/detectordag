package shared

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

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

func LogErrorAndReturn(err error) {
	log.Printf("%+v\n", wrapError(err))
}

func LogErrorAndExit(err error) {
	log.Fatalf("%+v\n", wrapError(err))
}

func wrapError(err error) error {
	return errors.Wrap(err, "detectordag error:")
}
