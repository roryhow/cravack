# Cravack - A serverless connector from Strava to Slack.
CR (Columbia Road) - ava (Strava) - ck (slack).

Cravack is a connection layer between Strava and Slack. It publishes all kinds of activity data for a given user within a designated slack channel.

## The Architecture.
This project is built using Go, and is centred around a collection of serverless functions which each provide a piece of functionality. Most of the serverless functions are triggered by HTTP requests, either from Cravack users, or by AQS SQS (a managed pub/sub queue service).

### AWS
This architecture was created to minimise AWS running costs. To this end, all services are selected due to the fact that they exist within the AWS free tier. Thus, Cravack does not have any monthly billing costs (assuming that it doesn't start to receive high amounts of traffic). The primary AWS services used for Cravack are as follows:
- *AWS Lambda:* Used for each of the functions that are called when Cravack is interacted with.
- *AWS DynamoDB:* Used to store Cravack user data, and Cravack users' activity data.
- *AWS SQS:* Used to asynchronously process event data which can be published from any source. There is also a dead letter queue which stores events are unprocessable by Cravack.
- *AWS Parameter Store:*: Used to store application secrets for Strava and Slack.

### Slack
This application consists of `n` slack applications, where `n` is the number of environments that you would wish to use for running this service. Currently, Cravack uses 2 environments.
A template file for the slack application is included in this repository for demonstrative purposes.

### Strava
A Strava application is needed to be able to run Cravack for yourself. A way to create this application and documentation can be found [on Strava's developer portal](https://developers.strava.com).
This application uses Strava webhooks for a couple reasons:
1. The rate limits for Strava's APIs is ridiculously low (1000 requests daily)
2. Webhooks allow for the lambda's to be called fewer times, keeping costs down.

However, a drawback of this means that each individual user needs to authenticate Cravack to access their activity data. This is a necessary evil (and I suppose makes this whole thing GDPR-friendly?)

### Serverless
This application uses the [Serverless framework](https://github.com/serverless/serverless) for transpiling infrastructure to Cloudformation templates. Aside from this it isn't really required, and it could be removed in the future.

## How to install
This project is not intended to be run locally. Build times are low enough that deploying this directly takes less than a minute, which is enough for development to not be completely painful. However, it would be possible to create some sample events that would allow for each lambda to be invoked locally. This is not in the immediate horizon for development.

If you want to run this application for yourself, you will need the following:
- Golang
- An AWS account
- A Slack bot (following the template supplied in this repo)
- A Strava application
- The serverless framework
- [nvm](https://github.com/nvm-sh/nvm). The latest version of node doesn't work with serverless right now, so I am using node V14 LTS.

You will also need to add parameters into your AWS parameter store in order for your application to work. These are the following:
- `/cravack/<environment>/slack-api-key`
- `/cravack/<environment>/strava-client-id`
- `/cravack/<environment>/strava-client-secret`
- `/cravack/<environment>/strava-webhook-verify-token`
These are relatively self explanatory. You can find each of the values needed from the respective slack / strava apps that you create.

You can then deploy this application by running `make`. It's that easy!

## Contribution
You can contribute to the development of Cravack in a few ways:
1. If you encounter a bug with Cravack (and do not want to develop the fix), then please file an issue in this repository.
2. If you have a suggestion on how to expand the functionalities of Cravack, please leave an issue in this repository.
3. If you would like to develop Cravack, then please fork, develop your changes, and make a pull request to this repository.

## Roadmap
### v0.0.1
- [x] Install slack APIs Golang
- [x] Create endpoints to handle interactions from slack, via messages etc.
- [x] Create database that stores users authentication credentials, and refreshes when outdated
- [x] Upon subscribing, bot posts events into a set channel when a user records an event.
- [x] Fetch Strava activity data for user when activity is received via webhook
- [x] Bot provides link in channel upon entry to "subscribe"

### V0.1
- [x] Handle Strava events asynchronously, pass to a SQS and consume in another lambda
- [x] Database to remember channels to post into, can be any channel within any organisation
- [x] Record past history of events, to be able to update messages
- [x] Post message only shown to user when they have authenticated
- [x] Include slack bot template file in repository.

### Nice to have
- [ ] Unit tests for business logic areas of the application
- [x] Proper dev environment, migrate actual users to production environment
- [ ] Interaction to choose with bot which event types to post
- [ ] Render images in slack messages, if uploaded in Strava activity
- [ ] Interaction with bot to choose which data to post (organisation specific settings?)
