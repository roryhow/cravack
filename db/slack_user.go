package db

type SlackUser struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	ChannelID    string `json:"channel_id"`
	TeamID       string `json:"team_id"`
	EnterpriseID string `json:"enterprise_id"`
}

func NewSlackUser(userID, userName, channelID, teamID, enterpriseID string) *SlackUser {
	return &SlackUser{userID, userName, channelID, teamID, enterpriseID}
}
