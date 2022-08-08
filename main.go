package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// get command line flags
// if not defined, use ENV variables
var tempoOauthToken string
var jiraOauthToken string
var jiraUsername string

type Worklog struct {
	Issue struct {
		Self string `json:"self"`
	} `json:"issue"`
}

type TempoResponse struct {
	Results []Worklog `json:"results"`
}

type JiraResponse struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

func makeTempoRequest() []string {
	url := "https://api.tempo.io/core/3/worklogs"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	todaysDate := time.Now().Weekday()

	var formattedDate string
	if todaysDate == time.Monday {
		formattedDate = time.Now().Add(-72 * time.Hour).Format("2006-01-02")
	} else {
		formattedDate = time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	}

	requestQuery := req.URL.Query()

	requestQuery.Add("from", formattedDate)
	requestQuery.Add("to", formattedDate)
	req.URL.RawQuery = requestQuery.Encode()

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	authHeader := fmt.Sprintf("Bearer %s", tempoOauthToken)
	req.Header.Add("Authorization", authHeader)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var response TempoResponse
	if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Cannot unmarshal JSON")
	}

	var issueLinks []string
	for _, worklog := range response.Results {
		issueLinks = append(issueLinks, worklog.Issue.Self)
	}
	return issueLinks
}

func makeJiraRequest(issue string) string {
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

// TODO: Implement Slack Bot to communicate with Slack Geekbot directly
func main() {
	if tempoOauthToken == "" || jiraOauthToken == "" || jiraUsername == "" {
		_ = godotenv.Load()
		flag.StringVar(&tempoOauthToken, "tempoOauthToken", os.Getenv("TEMPO_OAUTH_TOKEN"), "a string")
		flag.StringVar(&jiraOauthToken, "jiraOauthToken", os.Getenv("JIRA_OAUTH_TOKEN"), "a string")
		flag.StringVar(&jiraUsername, "jiraUsername", os.Getenv("JIRA_USERNAME"), "a string")
		flag.Parse()
	}
	issues := makeTempoRequest()

	report := ""
	for index, issue := range issues {
		issueStr := makeJiraRequest(issue)
		report += issueStr

		if index != len(issues)-1 {
			report += "\n"
		}
	}
	fmt.Println(report)
	clipboard.WriteAll(report)
}
