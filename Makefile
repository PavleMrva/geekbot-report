build:
	export $$(grep -v '^#' .env | xargs) && \
 	go build -ldflags \
 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.slackOauthToken=$$SLACK_OAUTH_TOKEN' \
 	-X 'main.geekBotChannel=$$GEEKBOT_CHANNEL_ID'"
install:
	export $$(grep -v '^#' .env | xargs) && \
 	go install -ldflags \
 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.slackOauthToken=$$SLACK_OAUTH_TOKEN' \
 	-X 'main.geekBotChannel=$$GEEKBOT_CHANNEL_ID'"

