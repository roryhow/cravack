.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/getBin getFolder/getExample.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/postBin postFolder/postExample.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getQueryBin getFolder/getQueryExample.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

