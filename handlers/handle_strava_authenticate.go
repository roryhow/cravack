package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roryhow/cravack/db"
	"github.com/roryhow/cravack/services"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if error := request.QueryStringParameters["error"]; len(error) > 0 {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error when authenticating: %s", error),
			StatusCode: 500,
		}, nil
	}

	// Pull needed params from query and environment
	userScope := request.QueryStringParameters["scope"]
	if userScope != "read,activity:read" {
		return events.APIGatewayProxyResponse{
			Body:       "Incorrect user scope for authentication",
			StatusCode: 400,
		}, nil
	}

	userAuthCode := request.QueryStringParameters["code"]
	authInfo, err := services.AuthenticateStravaUser(userAuthCode)
	if err != nil {
		log.Printf("Error when authenticating Strava user:\n%s", err.Error())
		return events.APIGatewayProxyResponse{
			Body:       "Error when authenticating Strava user",
			StatusCode: 500,
		}, nil
	}

	_, err = db.PutAuthenticatedUser(authInfo)
	if err != nil {
		log.Printf("Error when adding authenticated user to database:\n%s", err.Error())
		return events.APIGatewayProxyResponse{
			Body:       "Error when adding authenticated user to database",
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Authenticated!",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
