package db

import (
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CravackUser struct {
	UserID     int         `json:"user_id" dynamodbav:"UserID"`
	StravaUser *StravaUser `json:"strava_user" dynamodbav:"StravaUser"`
	SlackUser  *SlackUser  `json:"slack_user" dynamodbav:"SlackUser"`
}

func NewCravackUser(stravaUser *StravaUser, slackUser *SlackUser) *CravackUser {
	return &CravackUser{
		UserID:     stravaUser.AthleteID,
		StravaUser: stravaUser,
		SlackUser:  slackUser,
	}
}

func PutCravackUser(user *CravackUser) (*dynamodb.PutItemOutput, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Printf("Error when trying to marshal map")
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("CRAVACK_USER_TABLE")),
	}
	output, err := svc.PutItem(input)

	return output, err
}

func GetAuthenticatedUser(athleteID int) (*CravackUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("CRAVACK_USER_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(strconv.Itoa(athleteID)),
			},
		},
	})

	if err != nil {
		log.Printf("Error when fetching from database\n%s", err.Error())
		return nil, err
	}

	user := CravackUser{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &user, nil
}

func UpdateCravackStravaToken(refreshedUser *StravaRefreshToken, athleteID int) (*StravaUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	expr, err := dynamodbattribute.MarshalMap(refreshedUser)
	if err != nil {
		log.Printf("Error when marshalling refresh token:\n%s", err.Error())
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(os.Getenv("CRAVACK_USER_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"AthleteID": {
				N: aws.String(strconv.Itoa(athleteID)),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
		UpdateExpression: aws.String(`
set StravaUser.TokenType = :t,
StravaUser.AccessToken = :a,
StravaUser.ExpiresIn = :ei,
StravaUser.ExpiresAt = :ea,
StravaUser.RefreshToken = :r`),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		log.Printf("Error when user token in database for athelete: %d\n%s", athleteID, err.Error())
		return nil, err
	}

	var updatedAthlete StravaUser
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedAthlete)
	if err != nil {
		log.Printf("Error when unmarshalling results from dynamodb update into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &updatedAthlete, nil
}
