# Geekbot Report Script

This script retrieves Tempo worklogs from Jira of previous work-day and copies the report to clipboard.
It was implemented out of pure laziness, to avoid dull copy-pasting.

## Installation

Download all dependencies
```bash
go get -d ./...
```

## How to use

### Main Method:

For easier setup copy `.env-example` to `.env` and fill in the necessary data (Tempo OAuth Token, Jira Username/Email and Jira OAuth Token)
```bash
cp .env-example .env
```

Next, the `make build` command should be run. It will extract defined ENV variables and
the Go program will be compiled with these variables which will be the default values of the following variables:
- `tempoOauthToken`
- `jiraOauthToken`
- `jiraUsername`
- `slackOauthToken`
- `geekBotChannel`


```bash
# background commands:
#	export $$(grep -v '^#' .env | xargs) && \
# 	go build -ldflags \
# 	"-X 'main.jiraUsername=$$JIRA_USERNAME' \
# 	-X 'main.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
# 	-X 'main.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
# 	-X 'main.slackOauthToken=$$SLACK_OAUTH_TOKEN' \
# 	-X 'main.geekBotChannel=$$GEEKBOT_CHANNEL_ID'"
make build
```

After that just run the binary and the geekbot report will be generated and sent
```bash
./geekbot-report
```

If everything works well, you can finally install the Go program using the following command:
```bash
# background commands:
#	export $$(grep -v '^#' .env | xargs) && \
# 	go install -ldflags \
# 	"-X 'main.jiraUsername=$$JIRA_USERNAME' \
# 	-X 'main.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
# 	-X 'main.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
# 	-X 'main.slackOauthToken=$$SLACK_OAUTH_TOKEN' \
# 	-X 'main.geekBotChannel=$$GEEKBOT_CHANNEL_ID'"
make install
```
> **_NOTE_** You should have `$GOPATH/bin` inside your `$PATH`