.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/registerStravaWebhook handlers/register_strava_webhook.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleStravaActivityEvent handlers/handle_strava_activity_event.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleSlackInteractionEvent handlers/handle_slack_interaction_event.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleSlackEvent handlers/handle_slack_event.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleSlackSlashCommand handlers/handle_slack_slash_command.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleStravaAuthenticate handlers/handle_strava_authenticate.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
