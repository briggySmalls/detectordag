package app

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

// HandleRequest handles a lambda call
func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// For now just print out details
	log.Print("Request body: ", event)
	log.Print("Context: ", ctx)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
