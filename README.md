# Geekbot Report Script

This script retrieves Tempo worklogs from Jira of previous work-day and sends the daily standup report. 

## Installation

Download all dependencies:
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
- `geekBotApiKey`


```bash
# background commands:
#	export $$(grep -v '^#' .env | xargs) && \
# 	go build -ldflags \
# 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
# 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
# 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
# 	-X 'geekbot-report/external.geekBotApiKey=$$GEEKBOT_API_KEY'"
make build
```

After that run the binary and the geekbot report will be generated and sent:
```bash
./geekbot-report
```

If everything works well, finally install the Go program using the following command:
```bash
# background commands:
#	export $$(grep -v '^#' .env | xargs) && \
# 	go install -ldflags \
# 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
# 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
# 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
# 	-X 'geekbot-report/external.geekBotApiKey=$$GEEKBOT_API_KEY'"
make install
```
> **_NOTE_** You should have `$GOPATH/bin` inside your `$PATH` in order to run installed Go program