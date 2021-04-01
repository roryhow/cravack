package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SlackEvent struct {
	Token        string                 `json:"token"`
	TeamID       string                 `json:"team_id"`
	APIAppID     string                 `json:"api_app_id"`
	EventContext string                 `json:"event_context"`
	EventID      string                 `json:"event_id"`
	EventType    string                 `json:"type"`
	Event        map[string]interface{} `json:"event"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event SlackEvent

	if err := json.Unmarshal([]byte(request.Body), &event); err != nil {
		log.Printf("Unable to parse JSON payload:\n%s", request.Body)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	log.Printf("Received Slack event:\n%s", request.Body)
	return events.APIGatewayProxyResponse{
		Body:       request.Body,
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
