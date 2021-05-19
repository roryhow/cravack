package main

import (
	"log"
	"net/http/httputil"

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

	payload, ok := req.Form["payload"]
	if !ok || len(payload) < 1 {
		log.Printf("Missing payload from Slack response")
		return events.APIGatewayProxyResponse{Body: "Missing payload", StatusCode: 400}, nil
	}

	text, err := httputil.DumpRequest(req, true)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	log.Println(string(text))
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
