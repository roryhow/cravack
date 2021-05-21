package db

type CravackActivityEvent struct {
	*StravaEvent
	SlackMessageTimestamp string `dynamodbav:"SlackMessageTimestamp"`
}
