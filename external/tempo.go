package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Worklog struct {
	Issue struct {
		Self string `json:"self"`
	} `json:"issue"`
}

type TempoResponse struct {
	Results []Worklog `json:"results"`
}

var tempoOauthToken string

func MakeTempoRequest() []string {
	if tempoOauthToken == "" {
		tempoOauthToken = os.Getenv("TEMPO_OAUTH_TOKEN")
	}

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
