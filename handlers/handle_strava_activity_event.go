package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roryhow/cravack/db"
	"github.com/roryhow/cravack/services"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var bodyRequest services.StravaEvent
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		log.Printf("unable to decode JSON payload")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	// Don't handle anything other than creates for now
	if bodyRequest.AspectType != "create" {
		return events.APIGatewayProxyResponse{StatusCode: 204}, nil
	}

	// get the user auth details from the db
	cravackUser, err := db.GetAuthenticatedUser(bodyRequest.AthleteID)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	stravaUser := cravackUser.StravaUser

	// if the user token has expired, refresh it
	if int64(stravaUser.ExpiresAt) < time.Now().Unix() {
		refreshToken, err := services.GetStravaUserRefreshToken(stravaUser.RefreshToken)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}

		stravaUser, err = db.UpdateCravackStravaToken(refreshToken, stravaUser.AthleteID)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
	}

	// Fetch the corresponding event from strava api
	activity, err := services.GetStravaActivityForUser(&bodyRequest, stravaUser)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	host := request.Headers["host"]
	if len(host) <= 0 {
		log.Println("Host missing in request header, unable to post to channel")
		return events.APIGatewayProxyResponse{Body: "Host missing in header", StatusCode: 500}, nil
	}

	// Send the event to slack
	services.PostActivityToChannel(activity, stravaUser, cravackUser.SlackUser.ChannelID, host)

	// marshall the request back into a json response
	response, err := json.Marshal(&activity)
	if err != nil {
		log.Printf("unable to parse JSON response")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
