package main

import (
	"fmt"
	"log"
	"strings"

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

	// State is [UserID,UserName,ChannelID,TeamID,EnterpriseID]
	state := request.QueryStringParameters["state"]
	stateSlice := strings.Split(state, ",")

	if len(stateSlice) != 5 {
		return events.APIGatewayProxyResponse{
			Body:       "Error when parsing state in authentication step",
			StatusCode: 500,
		}, nil
	}
	slackUser := db.NewSlackUser(
		stateSlice[0],
		stateSlice[1],
		stateSlice[2],
		stateSlice[3],
		stateSlice[4],
	)

	userAuthCode := request.QueryStringParameters["code"]
	stravaUser, err := services.AuthenticateStravaUser(userAuthCode)
	if err != nil {
		log.Printf("Error when authenticating Strava user:\n%s", err.Error())
		return events.APIGatewayProxyResponse{
			Body:       "Error when authenticating Strava user",
			StatusCode: 500,
		}, nil
	}

	u := db.NewCravackUser(stravaUser, slackUser)
	_, err = db.PutCravackUser(u)
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
