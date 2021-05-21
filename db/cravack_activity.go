package db

type CravackActivityEvent struct {
	StravaEvent
	SlackChannelId        string `dynamodbav:"SlackChannelID"`
	SlackMessageTimestamp string `dynamodbav:"SlackMessageTimestamp"`
}

func NewCravackActivityEvent(event *StravaEvent, channelId, slackMsgTs string) *CravackActivityEvent {
	return &CravackActivityEvent{*event, channelId, slackMsgTs}
}
