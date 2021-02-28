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

	user, err := db.GetAuthenticatedUser(bodyRequest.AthleteID)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	if int64(user.ExpiresAt) < time.Now().Unix() {
		refreshToken, err := services.GetStravaUserRefreshToken(user.RefreshToken)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}

		user, err = db.UpdateStravaUserToken(refreshToken, user.AthleteID)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
	}

	// Fetch the corresponding event from strava api

	// marshall the request back into a json response
	response, err := json.Marshal(&user)
	if err != nil {
		log.Printf("unable to parse JSON response")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	log.Printf(string(response))
	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
