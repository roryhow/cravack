package services

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/roryhow/cravack/db"
)

func SendStravaActivityEventMessage(event *db.StravaEvent) (string, error) {
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)

	queueUrlOutput, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(os.Getenv("CRAVACK_ACTIVITY_QUEUE")),
	})
	if err != nil {
		log.Printf("Error when getting SQS queue URL:\n%s", err.Error())
		return "", err
	}

	sendMessageOutput, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"EventSource": {
				DataType:    aws.String("String"),
				StringValue: aws.String("strava"),
			},
			"ObjectType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.ObjectType),
			},
			"ObjectID": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(event.ObjectID)),
			},
			"AspectType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.AspectType),
			},
			"AthleteID": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(event.AthleteID)),
			},
			"SubscriptionID": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(event.SubscriptionID)),
			},
			"EventTime": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.Itoa(event.EventTime)),
			},
		},
		MessageBody: aws.String(fmt.Sprintf("Strava Event %d", event.ObjectID)),
		QueueUrl:    queueUrlOutput.QueueUrl,
	})
	if err != nil {
		log.Printf("Error when sending message to SQS queue:\n%s", err.Error())
		return "", err
	}

	return *sendMessageOutput.MessageId, nil
}
