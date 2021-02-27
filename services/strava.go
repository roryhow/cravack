package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/roryhow/cravack/db"
)

func AuthenticateStravaUser(userAuthCode string) (*db.AuthenticatedStravaUser, error) {
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

	var stravaResponse db.AuthenticatedStravaUser
	if err := json.NewDecoder(resp.Body).Decode(&stravaResponse); err != nil {
		return nil, err
	}

	log.Printf("Created user:\n%+v", stravaResponse)
	return &stravaResponse, nil
}
