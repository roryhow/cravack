package main

import (
	"errors"
	"fmt"
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

	if stravaEvent.AspectType == "delete" {
		// delete the cravack event, return the data before it was deleted
		cravackActivityEvent, err := services.DeleteCravackActivityEvent(stravaEvent)
		if err != nil {
			return err
		}

		// Use the activity event information to delete from Slack
		err = services.DeleteActivityMessage(cravackActivityEvent)
		if err != nil {
			return err
		}

		// delete completed successfully - nothing left to do
		return nil
	}

	// get the user auth details from the db
	cravackUser, err := services.GetCravackUser(stravaEvent.AthleteID)
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

	if stravaEvent.AspectType == "create" {
		// Send the event to slack
		channelId, ts, err := services.PostActivityMessage(activity, cravackUser)
		if err != nil {
			return err
		}

		// Create a database entry for the activity event
		_, err = services.PutCravackActivityEvent(stravaEvent, channelId, ts)
		return err
	} else if stravaEvent.AspectType == "update" {
		// Get the previous event, update the slack message
		cravackActivityEvent, err := services.GetCravackActivityEvent(stravaEvent)
		if err != nil {
			return err
		}

		// Update the existing activity message in Slack
		channelId, ts, err := services.UpdateActivityMessage(activity, cravackUser, cravackActivityEvent)
		if err != nil {
			return err
		}
		// Create a database entry for the actibity event
		_, err = services.PutCravackActivityEvent(stravaEvent, channelId, ts)
		return err
	}

	return errors.New(fmt.Sprintf("Unknown AspectType %s for Strava Event %+v", stravaEvent.AspectType, stravaEvent))
}

func Handler(sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		log.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		eventSource := *message.MessageAttributes["EventSource"].StringValue
		if eventSource == "strava" {
			err := handleStravaEventMessage(message.MessageAttributes) // not handling errors for now
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
