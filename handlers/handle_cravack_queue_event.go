package main

import (
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roryhow/cravack/db"
	"github.com/roryhow/cravack/services"
)

func messageAttributeToString(msg map[string]events.SQSMessageAttribute, param string) string {
	return *msg[param].StringValue
}

func messageAttributeToInt(msg map[string]events.SQSMessageAttribute, param string) int {
	result, _ := strconv.Atoi(*msg[param].StringValue)
	return result
}

func handleStravaEventMessage(message map[string]events.SQSMessageAttribute) error {
	stravaEvent := &db.StravaEvent{
		ObjectType:     messageAttributeToString(message, "ObjectType"),
		ObjectID:       messageAttributeToInt(message, "ObjectID"),
		AspectType:     messageAttributeToString(message, "AspectType"),
		AthleteID:      messageAttributeToInt(message, "AthleteID"),
		SubscriptionID: messageAttributeToInt(message, "SubscriptionID"),
		EventTime:      messageAttributeToInt(message, "EventTime"),
	}

	// Don't handle anything other than creates and updates for now
	if stravaEvent.AspectType != "create" && stravaEvent.AspectType != "update" {
		return nil
	}

	// get the user auth details from the db
	cravackUser, err := services.GetAuthenticatedUser(stravaEvent.AthleteID)
	if err != nil {
		return err
	}

	stravaUser := cravackUser.StravaUser

	// if the user token has expired, refresh it
	if int64(stravaUser.ExpiresAt) < time.Now().Unix() {
		refreshToken, err := services.GetStravaUserRefreshToken(stravaUser.RefreshToken)
		if err != nil {
			return err
		}

		cravackUser, err = services.UpdateCravackStravaToken(refreshToken, stravaUser.AthleteID)
		if err != nil {
			return err
		}
	}

	// Fetch the corresponding event from strava api
	activity, err := services.GetStravaActivityForUser(stravaEvent, stravaUser)
	if err != nil {
		return err
	}

	var channelId, ts string
	if stravaEvent.AspectType == "create" {
		// Send the event to slack
		channelId, ts, err = services.PostActivityToChannel(activity, cravackUser)
		if err != nil {
			return err
		}
	} else if stravaEvent.AspectType == "update" {
		// Get the previous event, update the slack message
		cravackActivityEvent, err := services.GetCravackActivityEvent(stravaEvent)
		if err != nil {
			return err
		}

		// FIXME I don't think this is working
		channelId, ts, err = services.UpdateActivityMessage(activity, cravackUser, cravackActivityEvent)
		if err != nil {
			return err
		}
	}

	// Create a database entry for the actibity event
	_, err = services.PutStravaActivityEvent(stravaEvent, channelId, ts)

	return err
}

func Handler(sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		log.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		eventSource := *message.MessageAttributes["EventSource"].StringValue
		if eventSource == "strava" {
			err := handleStravaEventMessage(message.MessageAttributes) // not handling errors for now
			if err != nil {
				return err
			}
		}
	}

	log.Println("All SQS events handled successfully")
	return nil
}

func main() {
	lambda.Start(Handler)
}
