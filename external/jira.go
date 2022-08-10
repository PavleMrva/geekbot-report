package external

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type JiraResponse struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

var jiraUsername string
var jiraOauthToken string

func MakeJiraRequest(issue string) string {
	if jiraUsername == "" || jiraOauthToken == "" {
		jiraUsername = os.Getenv("JIRA_USERNAME")
		jiraOauthToken = os.Getenv("JIRA_OAUTH_TOKEN")
	}

	url := issue
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	authString := fmt.Sprintf("%s:%s", jiraUsername, jiraOauthToken)
	encodedb64String := b64.StdEncoding.EncodeToString([]byte(authString))
	authHeader := fmt.Sprintf("Basic %s", encodedb64String)

	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var response JiraResponse
	if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Cannot unmarshal JSON")
	}

	return fmt.Sprintf("%s - %s", response.Key, response.Fields.Summary)
}
