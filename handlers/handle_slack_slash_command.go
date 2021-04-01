package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roryhow/cravack/services"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := services.APIGatewayProxyRequestToHTTPRequest(request)
	if err != nil {
		log.Printf("Error when converting Lambda event to HTTP Request: %s", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	if err := req.ParseForm(); err != nil {
		log.Printf("Error when parsing form in HTTP request: %s", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	log.Printf("Received Slack event:\n%s", req)

	services.SendSlackConnectMessage()

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
