package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func getUsers(api *slack.Client) {
	users, err := api.GetUsers()

	if err != nil {
		fmt.Printf("Error when getting users: %s\n", err)
	}

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}
}

func postMessageToTest(api *slack.Client) {
	// header
	headerText := slack.NewTextBlockObject("mrkdwn", "Hi all! :wave:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// subheader
	subheaderText := slack.NewTextBlockObject("mrkdwn", "Thanks for inviting me into your secret running club. I'm the bot that will hook your Strava activity into this channel :runner::bicyclist::swimmer:", false, false)
	subheaderSection := slack.NewSectionBlock(subheaderText, nil, nil)

	divider := slack.NewDividerBlock()

	bodyText := slack.NewTextBlockObject("mrkdwn", "To be able to see your activity in this channel, you'll need to authorise the Cravack application to access your Strava account. You can do that by clicking the button below :point_down:", false, false)
	bodySection := slack.NewSectionBlock(bodyText, nil, nil)

	// Authorise button
	authoriseBtnTxt := slack.NewTextBlockObject("plain_text", "Authorise Cravack to Strava", false, false)

	// TODO establish this from request to lambda rather than hardcode
	authCallbackUrl := "https://unepe1p44k.execute-api.eu-central-1.amazonaws.com/handleStravaAuthenticate"
	cravackClientID := os.Getenv("STRAVA_CLIENT_ID")
	stravaAuthUrl := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&approval_prompt=force&scope=read,activity:read", cravackClientID, authCallbackUrl)

	authoriseBtn := slack.ButtonBlockElement{
		Type: slack.METButton,
		Text: authoriseBtnTxt,
		URL:  stravaAuthUrl,
	}
	authoriseActionBlock := slack.NewActionBlock("", authoriseBtn)

	channelID, timestamp, err := api.PostMessage(
		"cr-cravack-test",
		slack.MsgOptionBlocks(
			headerSection,
			subheaderSection,
			divider,
			bodySection,
			authoriseActionBlock,
		),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func main() {
	api := slack.New(os.Getenv("SLACK_API_KEY"))

	postMessageToTest(api)
}
