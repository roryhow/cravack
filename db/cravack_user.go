package db

type CravackUser struct {
	UserID     int         `json:"user_id" dynamodbav:"UserID"`
	StravaUser *StravaUser `json:"strava_user" dynamodbav:"StravaUser"`
	SlackUser  *SlackUser  `json:"slack_user" dynamodbav:"SlackUser"`
}

func NewCravackUser(stravaUser *StravaUser, slackUser *SlackUser) *CravackUser {
	return &CravackUser{
		UserID:     stravaUser.AthleteID,
		StravaUser: stravaUser,
		SlackUser:  slackUser,
	}
}
