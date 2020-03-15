package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/briggysmalls/detectordag/api/app"
)

// main is the entrypoint to the lambda function
func main() {
	lambda.Start(app.HandleRequest)
}
