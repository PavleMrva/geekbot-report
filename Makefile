.DEFAULT_GOAL := build

build:
	export $$(grep -v '^#' .env | xargs) && \
 	go build -ldflags \
 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.geekBotAPIKey=$$GEEKBOT_API_KEY' \
 	-X 'geekbot-report/external.jiraUserID=$$JIRA_USER_ID'"
install:
	export $$(grep -v '^#' .env | xargs) && \
 	go install -ldflags \
 	"-X 'geekbot-report/external.jiraUsername=$$JIRA_USERNAME' \
 	-X 'geekbot-report/external.jiraOauthToken=$$JIRA_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.tempoOauthToken=$$TEMPO_OAUTH_TOKEN' \
 	-X 'geekbot-report/external.geekBotAPIKey=$$GEEKBOT_API_KEY' \
 	-X 'geekbot-report/external.jiraUserID=$$JIRA_USER_ID'"

