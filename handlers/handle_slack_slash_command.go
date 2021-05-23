package main

import (
	"encoding/json"
	"log"
	"strings"

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

	host := req.Header["Host"]
	if len(host) <= 0 {
		log.Printf("Host header required in order to form callback")
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	slashCommand := services.NewSlashCommandFromForm(&req.Form)
	text := strings.TrimSpace(slashCommand.Text)

	if text == "connect" {
		msg, _ := slashCommand.GetStravaConnectResponse(host[0])
		response, _ := json.Marshal(msg)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(response),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	} else if text == "disconnect" {
		// Get user info from DB
		cravackUser, err := services.GetCravackUserBySlackID(slashCommand.UserID)
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		// remove user information
		_, err = services.DeleteCravackUser(cravackUser)
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		// deauthorise strava
		err = services.DeauthorizeStravaForCravackUser(cravackUser)
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		// Get deauthorisation message
		msg, err := slashCommand.GetDeauthorizeCravackResponse(host[0])
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		response, err := json.Marshal(msg)
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(response),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	} else {
		// We can't handle the command supplied
		msg, err := slashCommand.GetUnknownCommandResponse(host[0])
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		response, err := json.Marshal(msg)
		if err != nil {
			return services.HandleErrorAndLambdaReturn(err, 500)
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(response),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
