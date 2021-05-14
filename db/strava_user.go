package db

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type StravaRefreshToken struct {
	TokenType    string `json:"token_type" dynamodbav:":t"`
	AccessToken  string `json:"access_token" dynamodbav:":a" validate:"required"`
	ExpiresAt    int    `json:"expires_at" dynamodbav:":ea"`
	ExpiresIn    int    `json:"expires_in" dynamodbav:":ei"`
	RefreshToken string `json:"refresh_token" dynamodbav:":r" validate:"required"`
}

type StravaUser struct {
	TokenType     string
	ExpiresAt     int
	ExpiresIn     int
	RefreshToken  string `validate:"required"`
	AccessToken   string `validate:"required"`
	AthleteID     int    `validate:"require"`
	Username      string `validate:"required"`
	FirstName     string
	LastName      string
	ProfileMedium string
}

func (a *StravaUser) UnmarshalJSON(buf []byte) error {
	var tmp struct {
		TokenType    string `json:"token_type"`
		ExpiresAt    int    `json:"expires_at"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
		Athlete      struct {
			AthleteID     int    `json:"id"`
			Username      string `json:"username"`
			FirstName     string `json:"firstname"`
			LastName      string `json:"lastname"`
			ProfileMedium string `json:"profile_medium"`
		}
	}

	if err := json.Unmarshal(buf, &tmp); err != nil {
		return errors.Wrap(err, "Error when pasring data into AuthenticatedStravaUser")
	}

	a.TokenType = tmp.TokenType
	a.ExpiresAt = tmp.ExpiresAt
	a.ExpiresIn = tmp.ExpiresIn
	a.RefreshToken = tmp.RefreshToken
	a.AccessToken = tmp.AccessToken
	a.AthleteID = tmp.Athlete.AthleteID
	a.Username = tmp.Athlete.Username
	a.FirstName = tmp.Athlete.FirstName
	a.LastName = tmp.Athlete.LastName
	a.ProfileMedium = tmp.Athlete.ProfileMedium

	return nil
}
