.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/registerStravaWebhook registerStravaWebhook/get.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handleStravaActivityEvent handleStravaActivityEvent/post.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

