package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/roryhow/cravack/db"
)

var validate *validator.Validate

func AuthenticateStravaUser(userAuthCode string) (*db.StravaUser, error) {
	clientId := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	client := http.DefaultClient
	req, err := http.NewRequest("POST", "https://www.strava.com/api/v3/oauth/token", nil)
	if err != nil {
		log.Printf("Failure to build request to auth: %s", err.Error())
		return nil, err
	}
	q := req.URL.Query()
	q.Add("client_id", clientId)
	q.Add("client_secret", clientSecret)
	q.Add("code", userAuthCode)
	q.Add("grant_type", "authorization_code") // this function will only be possible for initial authenticate
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failure when generating Strava OAuth token:\n%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var stravaResponse db.StravaUser
	if err := json.NewDecoder(resp.Body).Decode(&stravaResponse); err != nil {
		return nil, err
	}

	validate = validator.New()
	if err := validate.Struct(&stravaResponse); err != nil {
		log.Printf("Error when validating Strava user:\n%+v", stravaResponse)
		return nil, err
	}

	return &stravaResponse, nil
}

// Fetch Strava refresh token from Strava API
func GetStravaUserRefreshToken(refreshToken string) (*db.StravaRefreshToken, error) {
	clientId := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	client := http.DefaultClient
	req, err := http.NewRequest("POST", "https://www.strava.com/api/v3/oauth/token", nil)
	if err != nil {
		log.Printf("Failure to build request to auth: %s", err.Error())
		return nil, err
	}
	q := req.URL.Query()
	q.Add("client_id", clientId)
	q.Add("client_secret", clientSecret)
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failure when generating Strava OAuth token:\n%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var stravaResponse db.StravaRefreshToken
	if err := json.NewDecoder(resp.Body).Decode(&stravaResponse); err != nil {
		return nil, err
	}

	log.Printf("Strava Response: %+v", stravaResponse)
	validate = validator.New()
	if err := validate.Struct(&stravaResponse); err != nil {
		return nil, err
	}

	return &stravaResponse, nil

}

func GetStravaActivityForUser(event *db.StravaEvent, user *db.StravaUser) (*db.StravaEventFull, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.strava.com/api/v3/activities/%d", event.ObjectID), nil)
	if err != nil {
		log.Printf("Error when building request to strava to fetch activity:\n%s", err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failure when fetching Strava activity")
	}

	var stravaEventFull db.StravaEventFull
	if err := json.NewDecoder(resp.Body).Decode(&stravaEventFull); err != nil {
		return nil, err
	}

	return &stravaEventFull, nil
}
