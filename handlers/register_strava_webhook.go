package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ValidationResponse struct {
	ChallengeResponse string `json:"hub.challenge"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mode := request.QueryStringParameters["hub.mode"]
	challenge := request.QueryStringParameters["hub.challenge"]
	verifyToken := request.QueryStringParameters["hub.verify_token"]

	if mode == "subscribe" && verifyToken == os.Getenv("STRAVA_WEBHOOK_VERIFY_TOKEN") {
		validationResponse := ValidationResponse{
			ChallengeResponse: challenge,
		}

		response, err := json.Marshal(&validationResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}

		responseHeaders := map[string]string{
			"Content-Type": "application/json",
		}

		return events.APIGatewayProxyResponse{
			Body:       string(response),
			Headers:    responseHeaders,
			StatusCode: 200,
		}, nil
	}

	return events.APIGatewayProxyResponse{Body: "Invalid Request", StatusCode: 400}, nil
}

func main() {
	lambda.Start(Handler)
}
