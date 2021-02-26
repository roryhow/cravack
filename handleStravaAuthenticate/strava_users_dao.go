package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutAuthenticatedUser(user *AuthenticatedStravaUser) (*dynamodb.PutItemOutput, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Printf("Error when trying to marshal map")
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("STRAVA_USER_AUTH_TABLE")),
	}
	output, err := svc.PutItem(input)

	return output, err
}
