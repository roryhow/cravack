package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roryhow/cravack/db"
	"github.com/roryhow/cravack/services"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var bodyRequest db.StravaEvent
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		log.Printf("unable to decode JSON payload")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	_, err = services.SendStravaActivityEventMessage(&bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
