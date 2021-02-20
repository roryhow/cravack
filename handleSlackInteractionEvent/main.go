package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type InteractiveResponse struct{}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var parsedResponse InteractiveResponse
	if err := json.Unmarshal([]byte(request.Body), &parsedResponse); err != nil {
		log.Printf("JSON payload:\n%s", request.Body)
		log.Printf("Is payload Base64 encoded? %t", request.IsBase64Encoded)
		log.Printf("Encountered an error when parsing JSON payload: %s", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	log.Printf("%+v\n", parsedResponse)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
