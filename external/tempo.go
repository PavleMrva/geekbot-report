package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Author struct {
	AccountId string `json:"accountId"`
}

type Worklog struct {
	Issue struct {
		Self string `json:"self"`
	} `json:"issue"`
	Author Author `json:"author"`
}

type TempoResponse struct {
	Results []Worklog `json:"results"`
}

var tempoOauthToken string
var jiraUserId string

func MakeTempoRequest(date string) ([]string, error) {
	if tempoOauthToken == "" {
		tempoOauthToken = os.Getenv("TEMPO_OAUTH_TOKEN")
	}
	if jiraUserId == "" {
		jiraUserId = os.Getenv("JIRA_USER_ID")
	}

	url := "https://api.tempo.io/core/3/worklogs"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}

	today := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	requestQuery := req.URL.Query()

	requestQuery.Add("from", date)
	requestQuery.Add("to", today)
	req.URL.RawQuery = requestQuery.Encode()

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	authHeader := fmt.Sprintf("Bearer %s", tempoOauthToken)
	req.Header.Add("Authorization", authHeader)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}

	var response TempoResponse
	if err = json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Cannot unmarshal JSON")
		return []string{}, err
	}

	var issueLinks []string
	for _, worklog := range response.Results {
		if jiraUserId == "" || jiraUserId == worklog.Author.AccountId {
			issueLinks = append(issueLinks, worklog.Issue.Self)
		}
	}
	return issueLinks, nil
}
