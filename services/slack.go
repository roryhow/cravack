package services

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/roryhow/cravack/db"
	"github.com/slack-go/slack"
)

func SendSlackConnectMessage() {
	api := slack.New(os.Getenv("SLACK_API_KEY"))

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
		"cr-half-marathon",
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

func getHeaderTextForActivityType(activityType string, name string) string {
	a := fmt.Sprintf("%s did a workout!", name)
	switch activityType {
	case "AlpineSki":
		a = fmt.Sprintf("%s went alpine skiing! :skier:", name)
	case "BackcountrySki":
		a = fmt.Sprintf("%s went canoeing! :canoe:", name)
	case "Crossfit":
		a = fmt.Sprintf("%s did some crossfit :runner:", name)
	case "EBikeRide":
		a = fmt.Sprintf("%s went on an e-bike ride? What is that? :bicyclist:", name)
	case "Elliptical":
		a = fmt.Sprintf("%s went on the elliptical! :runner:", name)
	case "Golf":
		a = fmt.Sprintf("%s played some golf :golfer:", name)
	case "Handcycle":
		a = fmt.Sprintf("%s went on their handcycle!", name)
	case "Hike":
		a = fmt.Sprintf("%s went for a hike! :hiking_boot:", name)
	case "IceSkate":
		a = fmt.Sprintf("%s did some ice skating! :ice_skate:", name)
	case "InlineSkate":
		a = fmt.Sprintf("%s did some inline skating! :ice_skate:", name)
	case "Kayaking":
		a = fmt.Sprintf("%s did some kite-surfing", name)
	case "NordicSki":
		a = fmt.Sprintf("%s did some cross-country skiing! :skier:", name)
	case "Ride":
		a = fmt.Sprintf("%s went for a bike ride! :bicyclist:", name)
	case "RockClimbing":
		a = fmt.Sprintf("%s went rock climbing! :person_climbing:", name)
	case "RollerSki":
		a = fmt.Sprintf("%s went roller skiing! :roller_skate:", name)
	case "Rowing":
		a = fmt.Sprintf("%s went rowing! :rowboat:", name)
	case "Run":
		a = fmt.Sprintf("%s went for a run! :runner:", name)
	case "Sail":
		a = fmt.Sprintf("%s went sailing! :sailboat:", name)
	case "Skateboard":
		a = fmt.Sprintf("%s went skateboarding! Cowabunga! :skateboard:", name)
	case "Snowboard":
		a = fmt.Sprintf("Perhaps not as fun as skiing, but %s went snowboarding! :snowboarder:", name)
	case "Snowshoe":
		a = fmt.Sprintf("%s went snowshoeing :snowflake:", name)
	case "Soccer":
		a = fmt.Sprintf("%s played some football! :soccer:", name)
	case "StairStepper":
		a = fmt.Sprintf("%s went stepping on some stairs! :foot:", name)
	case "StandUpPaddling":
		a = fmt.Sprintf("%s did some SUP boarding! :surfer:", name)
	case "Surfing":
		a = fmt.Sprintf("%s caught some gnarly waves and went surfing! :surfer:", name)
	case "Swin":
		a = fmt.Sprintf("%s went for a swim! :swimmer:", name)
	case "Velomobile":
		a = fmt.Sprintf("%s went on their velomobile... whatever that is? :shrug:", name)
	case "VirtualRide":
		a = fmt.Sprintf("%s went for a virtual ride! :bicyclist:", name)
	case "VirtualRun":
		a = fmt.Sprintf("%s went for a virtual run! Is that online... or? :globe_with_meridians::runner:", name)
	case "Walk":
		a = fmt.Sprintf("%s went for a leisurely walk :walking:", name)
	case "WeightTraining":
		a = fmt.Sprintf("%s is feeling the gains because they just went weight training! :weight_lifter:", name)
	case "Wheelchair":
		a = fmt.Sprintf("%s knows that sitting doesn't always have to be still because they just went on their wheelchair! :person_in_manual_wheelchair:", name)
	case "Windsurf":
		a = fmt.Sprintf("%s just went windsurfing! :surfer:", name)
	case "Workout":
		a = fmt.Sprintf("%s just did a workout! :runner:", name)
	case "Yoga":
		a = fmt.Sprintf("%s is feeling zen because they just did some yoga :person_in_lotus_position:", name)
	}

	return a
}

func metersPerSecondToMinutesPerKm(speed float64) string {
	pace := speed / (60 * 1000)
	leftover := math.Mod(pace, 1)
	minutes := int(pace - leftover)
	seconds := int(math.Round(leftover * 60))
	return fmt.Sprintf("%d:%d", minutes, seconds)
}

func PostActivityToChannel(activity *StravaEventFull, user *db.AuthenticatedStravaUser, channelID string) {
	api := slack.New(os.Getenv("SLACK_API_KEY"))

	// Title text
	headerText := slack.NewTextBlockObject(
		"mrkdwn",
		getHeaderTextForActivityType(activity.Type, user.FirstName),
		false,
		false,
	)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	subHeaderText := slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf(":speech_balloon: %s", activity.Name),
		false,
		false,
	)
	subHeaderSection := slack.NewContextBlock("", subHeaderText)

	// Components for stats block
	distanceText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Distance travelled:* %.2fkm", activity.Distance/1000), false, false)
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", activity.ElapsedTime))
	durationText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Duration:* %s", duration.String()), false, false)

	minsPerKm := metersPerSecondToMinutesPerKm(activity.AverageSpeed)
	paceText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Average Speed:* %smin/km", minsPerKm), false, false)
	elevationGainText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Elevation Gain:* %.2fm", activity.TotalElevationGain), false, false)
	statsSectionFields := []*slack.TextBlockObject{
		distanceText,
		durationText,
		paceText,
		elevationGainText,
	}
	statsSection := slack.NewSectionBlock(nil, statsSectionFields, nil)

	// Divider - purely visual
	divider := slack.NewDividerBlock()

	// Action buttons block
	authoriseBtnTxt := slack.NewTextBlockObject("plain_text", "Authorise Cravack to Strava", false, false)
	authCallbackUrl := "https://unepe1p44k.execute-api.eu-central-1.amazonaws.com/handleStravaAuthenticate"
	cravackClientID := os.Getenv("STRAVA_CLIENT_ID")
	stravaAuthUrl := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&approval_prompt=force&scope=read,activity:read", cravackClientID, authCallbackUrl)
	authoriseBtn := slack.ButtonBlockElement{
		Type: slack.METButton,
		Text: authoriseBtnTxt,
		URL:  stravaAuthUrl,
	}
	fullActivityBtnText := slack.NewTextBlockObject("plain_text", "View full activity on Strava", false, false)
	stravaFullActivityUrl := fmt.Sprintf("https://www.strava.com/activities/%d", activity.ID)
	fullActivityBtn := slack.ButtonBlockElement{
		Type: slack.METButton,
		Text: fullActivityBtnText,
		URL:  stravaFullActivityUrl,
	}
	authoriseActionBlock := slack.NewActionBlock("", fullActivityBtn, authoriseBtn)

	// Send the message to the channel
	channelID, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionBlocks(
			headerSection,
			subHeaderSection,
			statsSection,
			divider,
			authoriseActionBlock,
		),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
