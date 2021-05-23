package services

import (
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/roryhow/cravack/db"
)

func PutCravackUser(user *db.CravackUser) (*dynamodb.PutItemOutput, error) {
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

func GetCravackUser(athleteID int) (*db.CravackUser, error) {
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

	user := db.CravackUser{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &user, nil
}

func GetCravackUserBySlackID(slackUserID string) (*db.CravackUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("CRAVACK_USER_TABLE")),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String(slackUserID),
			},
		},
		FilterExpression: aws.String("SlackUser.UserID = :user"),
	})

	if err != nil {
		log.Printf("Error when fetching from database\n%s", err.Error())
		return nil, err
	}

	if *result.Count < 1 {
		return nil, errors.Errorf("No such user exists for slackID %s", slackUserID)
	}

	user := db.CravackUser{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &user, nil
}

func DeleteCravackUser(user *db.CravackUser) (*db.CravackUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("CRAVACK_EVENT_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(strconv.Itoa(user.UserID)),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	})
	if err != nil {
		return nil, err
	}

	cravackUser := db.CravackUser{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &cravackUser)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into CravackUser\n%s", err.Error())
		return nil, err
	}

	return &cravackUser, nil
}

func UpdateCravackStravaToken(refreshedUser *db.StravaRefreshToken, athleteID int) (*db.CravackUser, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	expr, err := dynamodbattribute.MarshalMap(refreshedUser)
	if err != nil {
		log.Printf("Error when marshalling refresh token:\n%s", err.Error())
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,

		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
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

	updatedAthlete := db.CravackUser{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedAthlete)
	if err != nil {
		log.Printf("Error when unmarshalling results from dynamodb update into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &updatedAthlete, nil
}

func PutCravackActivityEvent(event *db.StravaEvent, slackChannelId, slackMsgTs string) (*dynamodb.PutItemOutput, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	c := db.NewCravackActivityEvent(event, slackChannelId, slackMsgTs)
	av, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		log.Printf("Error when trying to marshal map")
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("CRAVACK_EVENT_TABLE")),
	}
	output, err := svc.PutItem(input)

	return output, err
}

func GetCravackActivityEvent(event *db.StravaEvent) (*db.CravackActivityEvent, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("CRAVACK_EVENT_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(strconv.Itoa(event.AthleteID)),
			},
			"EventID": {
				N: aws.String(strconv.Itoa(event.ObjectID)),
			},
		},
	})

	if err != nil {
		log.Printf("Error when fetching from database\n%s", err.Error())
		return nil, err
	}

	cravackEvent := db.CravackActivityEvent{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &cravackEvent)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &cravackEvent, nil
}

func DeleteCravackActivityEvent(event *db.StravaEvent) (*db.CravackActivityEvent, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("CRAVACK_EVENT_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(strconv.Itoa(event.AthleteID)),
			},
			"EventID": {
				N: aws.String(strconv.Itoa(event.ObjectID)),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	})
	if err != nil {
		return nil, err
	}

	cravackEvent := db.CravackActivityEvent{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &cravackEvent)
	if err != nil {
		log.Printf("Error when unmarshalling result from DB into StravaUser\n%s", err.Error())
		return nil, err
	}

	return &cravackEvent, nil
}
