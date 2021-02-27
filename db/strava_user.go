package db

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

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
				N: aws.String(string(athleteID)),
			},
		},
	})

	if err != nil {
		log.Printf("Error when fetching from database")
		return nil, errors.Wrap(err, "Error when fetching from database")
	}

	stravaUser := AuthenticatedStravaUser{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &stravaUser)
	if err != nil {
		errmsg := "Error when unmarshalling result from DB into AuthenticatedStravaUser"
		log.Printf(errmsg)
		return nil, errors.Wrap(err, errmsg)
	}

	return &stravaUser, nil
}
