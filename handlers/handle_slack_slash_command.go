package main

import (
	"encoding/json"
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

	host := req.Header["Host"]
	if len(host) <= 0 {
		log.Printf("Host header required in order to form callback")
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	slashCommand := services.NewSlashCommandFromForm(&req.Form)
	msg, _ := slashCommand.GetStravaConnectResponse(host[0])

	response, _ := json.Marshal(msg)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
