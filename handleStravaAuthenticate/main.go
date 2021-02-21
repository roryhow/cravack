package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if error := request.QueryStringParameters["error"]; len(error) > 0 {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error when authenticating: %s", error),
			StatusCode: 500,
		}, nil
	}

	// Pull needed params from query and environment
	code := request.QueryStringParameters["code"]
	cravackClientId := os.Getenv("STRAVA_CLIENT_ID")
	cravackClientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	client := http.DefaultClient
	req, err := http.NewRequest("POST", "https://www.strava.com/api/v3/oauth/token", nil)
	if err != nil {
		log.Printf("Failure to build request to auth: %s", err.Error())
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: 500}, nil
	}
	q := req.URL.Query()
	q.Add("client_id", cravackClientId)
	q.Add("client_secret", cravackClientSecret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failure when generating Strava OAuth token:\n%s", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failure reading Strava Oauth response body:\n%s", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
