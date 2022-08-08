build:
	export $(grep -v '^#' .env | xargs) && go build -ldflags "-X 'main.jiraUsername=$$JIRA_USERNAME' -X 'main.jiraOauthToken=$$JIRA_OAUTH_TOKEN' -X 'main.tempoOauthToken=$$TEMPO_OAUTH_TOKEN'"
install:
	mv ./geekbot-report $$GOPATH
