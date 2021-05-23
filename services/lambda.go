package services

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func HandleErrorAndLambdaReturn(err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	log.Printf(err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       err.Error(),
	}, nil
}
