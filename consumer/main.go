package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/briggysmalls/detectordag/consumer/app"
	"log"
)

func init() {
	// Add file/line number to the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(app.HandleRequest)
}
