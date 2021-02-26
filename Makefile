.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/registerStravaWebhook registerStravaWebhook/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleStravaActivityEvent handleStravaActivityEvent/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleSlackInteractionEvent handleSlackInteractionEvent/main.go handleSlackInteractionEvent/request.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleStravaAuthenticate handleStravaAuthenticate/main.go handleStravaAuthenticate/strava_auth.go handleStravaAuthenticate/strava_users_dao.go handleStravaAuthenticate/strava_user.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose
