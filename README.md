# Geekbot Report Script

This script retrieves Tempo worklogs from Jira of previous work-day and sends the daily standup report. 

## Installation

Download all dependencies:
```bash
go get -d ./...
```

## How to use

### Main Method:

For easier setup copy `.env-example` to `.env` and fill in the necessary data ([Tempo v4 API key](https://autaut.atlassian.net/plugins/servlet/ac/io.tempo.jira/tempo-app#!/configuration/api-integration), Jira Username/Email, [Jira API Token](https://id.atlassian.com/manage-profile/security/api-tokens) and [Geekbot api key](https://app.geekbot.com/dashboard/api-webhooks))
Setting Jira User ID is necessary only if you are seeing tickets from the other people in your worklog. Otherwise you can leave it empty. 
```bash
cp .env-example .env
```

Next, the `make build` command should be run. It will extract defined ENV variables and
the Go program will be compiled with these variables which will be the default values of the following variables:
- `tempoOauthToken`
- `jiraOauthToken`
- `jiraUsername`
- `geekBotAPIKey`
- `jiraUserID`


```bash
# background commands:
#	export $$(grep -v '^#' .env | xargs) && \
# 	go build -ldflags \
# 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
# 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
# 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
# 	-X 'geekbot-report/external.geekBotAPIKey=$$GEEKBOT_API_KEY' \
#   -X 'geekbot-report/external.jiraUserID=$$JIRA_USER_ID'"
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
# 	-X 'geekbot-report/external.geekBotAPIKey=$$GEEKBOT_API_KEY' \
#   -X 'geekbot-report/external.jiraUserID=$$JIRA_USER_ID'"
make install
```
> **_NOTE_** You should have `$GOPATH/bin` inside your `$PATH` in order to run installed Go program