package db

type SlackUser struct {
	UserID       string `json:"user_id" dynamodbav:"UserID"`
	UserName     string `json:"user_name" dynamodbav:"UserName"`
	ChannelID    string `json:"channel_id" dynamodbav:"ChannelID"`
	TeamID       string `json:"team_id" dynamodbav:"TeamID"`
	EnterpriseID string `json:"enterprise_id" dynamodbav:"EnterpriseID"`
}

func NewSlackUser(userID, userName, channelID, teamID, enterpriseID string) *SlackUser {
	return &SlackUser{userID, userName, channelID, teamID, enterpriseID}
}
