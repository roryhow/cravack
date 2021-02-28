package db

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type StravaRefreshToken struct {
	TokenType    string `json:"token_type" dynamodbav:":t"`
	AccessToken  string `json:"access_token" dynamodbav:":a"`
	ExpiresAt    int    `json:"expires_at" dynamodbav:":ea"`
	ExpiresIn    int    `json:"expires_int" dynamodbav:":ei"`
	RefreshToken string `json:"refresh_token" dynamodbav:":r"`
}

type AuthenticatedStravaUser struct {
	TokenType     string
	ExpiresAt     int
	ExpiresIn     int
	RefreshToken  string
	AccessToken   string
	AthleteID     int
	Username      string
	FirstName     string
	LastName      string
	ProfileMedium string
}

func (a *AuthenticatedStravaUser) UnmarshalJSON(buf []byte) error {
	var tmp struct {
		TokenType    string `json:"token_type"`
		ExpiresAt    int    `json:"expires_at"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
		Athlete      struct {
			AthleteID     int    `json:"id"`
			Username      string `json:"username"`
			FirstName     string `json:"firstname"`
			LastName      string `json:"lastname"`
			ProfileMedium string `json:"profile_medium"`
		}
	}

	if err := json.Unmarshal(buf, &tmp); err != nil {
		return errors.Wrap(err, "Error when pasring data into AuthenticatedStravaUser")
	}

	a.TokenType = tmp.TokenType
	a.ExpiresAt = tmp.ExpiresAt
	a.ExpiresIn = tmp.ExpiresIn
	a.RefreshToken = tmp.RefreshToken
	a.AccessToken = tmp.AccessToken
	a.AthleteID = tmp.Athlete.AthleteID
	a.Username = tmp.Athlete.Username
	a.FirstName = tmp.Athlete.FirstName
	a.LastName = tmp.Athlete.LastName
	a.ProfileMedium = tmp.Athlete.ProfileMedium

	return nil
}

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

func GetAuthenticatedUser(athleteID int) (*AuthenticatedStravaUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("STRAVA_USER_AUTH_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"AthleteID": {
				N: aws.String(strconv.Itoa(athleteID)),
			},
		},
	})

	if err != nil {
		log.Printf("Error when fetching from database\n%s", err.Error())
		return nil, err
	}

	stravaUser := AuthenticatedStravaUser{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &stravaUser)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into AuthenticatedStravaUser\n%s", err.Error())
		return nil, err
	}

	return &stravaUser, nil
}

func UpdateStravaUserToken(refreshedUser *StravaRefreshToken, athleteID int) (*AuthenticatedStravaUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	expr, err := dynamodbattribute.MarshalMap(refreshedUser)
	if err != nil {
		log.Printf("Error when marshalling refresh token:\n%s", err.Error())
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(os.Getenv("STRAVA_USER_AUTH_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"AthleteID": {
				N: aws.String(strconv.Itoa(athleteID)),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("set TokenType = :t, AccessToken = :a, ExpiresIn = :ei, ExpiresAt = :ea, RefreshToken = :r"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		log.Printf("Error when user token in database for athelete: %d\n%s", athleteID, err.Error())
		return nil, err
	}

	var updatedAthlete AuthenticatedStravaUser
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedAthlete)
	if err != nil {
		log.Printf("Error when unmarshalling results from dynamodb update into AuthenticatedStravaUser\n%s", err.Error())
		return nil, err
	}

	return &updatedAthlete, nil
}
