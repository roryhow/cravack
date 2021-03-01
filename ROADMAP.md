# Cravack roadmap

## Features

### MVP
- [x] Install slack APIs Golang
- [x] Create endpoints to handle interactions from slack, via messages etc.
- [x] Create database that stores users authentication credentials, and refreshes when outdated
- [x] Upon subscribing, bot posts events into a set channel when a user records an event.
- [ ] Fetch Strava activity data for user when activity is received via webhook
- [ ] Bot provides link in channel upon entry to "subscribe"

### Nice to have
- [ ] Basic DB to remember channels to post into, can be any channel within any organisation
- [ ] Interaction to choose with bot which event types to post
- [ ] Interaction with bot to choose which data to post
