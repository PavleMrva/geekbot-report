package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"geekbot-report/external"
	"github.com/joho/godotenv"
	"strings"
)

type StandupQuestion struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type Standup struct {
	Id        int               `json:"id"`
	Channel   string            `json:"channel"`
	Questions []StandupQuestion `json:"questions"`
}

type Answer struct {
	Text string `json:"text"`
}

type Report struct {
	StandupId string            `json:"standup_id"`
	Answers   map[string]Answer `json:"answers"`
}

func getDailyStandup() Standup {
	url := "https://api.geekbot.com/v1/standups"
	method := "GET"

	body := external.SendGeekBotRequest(url, method, nil)

	var response []Standup
	if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
		fmt.Println("Cannot unmarshal JSON")
	}

	var dailyStandup Standup
	for _, standup := range response {
		if standup.Channel == "daily-standup" {
			dailyStandup = standup
		}
	}
	return dailyStandup
}

func sendGeekBotReport(report Report) {
	url := "https://api.geekbot.com/v1/reports"
	method := "POST"

	jsonData, _ := json.Marshal(report)
	payload := bytes.NewBuffer(jsonData)

	external.SendGeekBotRequest(url, method, payload)
}

func main() {
	_ = godotenv.Load()
	issues := external.MakeTempoRequest()

	report := ""
	for index, issue := range issues {
		issueStr := external.MakeJiraRequest(issue)
		report += issueStr

		if index != len(issues)-1 {
			report += "\n"
		}
	}
	fmt.Println(report)
	dailyStandup := getDailyStandup()

	dailyStandupReport := Report{
		StandupId: fmt.Sprintf("%d", dailyStandup.Id),
		Answers:   make(map[string]Answer),
	}

	for _, question := range dailyStandup.Questions {
		key := fmt.Sprintf("%d", question.Id)
		if strings.Contains(question.Text, "What did") {
			dailyStandupReport.Answers[key] = Answer{report}
		} else if strings.Contains(question.Text, "What will") {
			dailyStandupReport.Answers[key] = Answer{"sprint/support"}
		} else {
			dailyStandupReport.Answers[key] = Answer{"-"}
		}
	}

	sendGeekBotReport(dailyStandupReport)
}
