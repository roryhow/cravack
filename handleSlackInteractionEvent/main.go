package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type InteractiveResponse struct{}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	decodedRequest, err := b64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Printf("Error when decoding request: %s", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	var parsedResponse InteractiveResponse
	if err := json.Unmarshal([]byte(decodedRequest), &parsedResponse); err != nil {
		log.Printf("JSON payload:\n%s", request.Body)
		log.Printf("Encountered an error when parsing JSON payload: %s", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	log.Printf("%+v\n", parsedResponse)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
