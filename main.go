package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"geekbot-report/external"
	"github.com/joho/godotenv"
	"log"
	"os"
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

func getDailyStandup() (Standup, error) {
	url := "https://api.geekbot.com/v1/standups"
	method := "GET"

	body, err := external.SendGeekBotRequest(url, method, nil)

	if err != nil {
		return Standup{}, err
	}

	var response []Standup
	if err = json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
		fmt.Println("Cannot unmarshal JSON")
		return Standup{}, err
	}

	var dailyStandup Standup
	for _, standup := range response {
		if standup.Channel == "daily-standup" {
			dailyStandup = standup
		}
	}
	return dailyStandup, nil
}

func sendGeekBotReport(report Report) ([]byte, error) {
	url := "https://api.geekbot.com/v1/reports"
	method := "POST"

	jsonData, _ := json.Marshal(report)
	payload := bytes.NewBuffer(jsonData)

	return external.SendGeekBotRequest(url, method, payload)
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	_ = godotenv.Load()
	issues, err := external.MakeTempoRequest()
	if err != nil {
		panic("Error while sending Tempo request")
	}

	report := ""
	for index, issue := range issues {
		issueStr, err := external.MakeJiraRequest(issue)
		if err != nil {
			panic("Error while sending Tempo request")
		}
		report += issueStr

		if index != len(issues)-1 {
			report += "\n"
		}
	}
	dailyStandup, err := getDailyStandup()

	if err != nil {
		panic("Error while fetching daily standup")
	}

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

	confirmationQuestion := fmt.Sprintf("Your report for today:\n%s\nAre you sure that you want to send this report?", report)
	answer := askForConfirmation(confirmationQuestion)

	if answer {
		_, err = sendGeekBotReport(dailyStandupReport)

		if err != nil {
			panic("Error while sending daily standup")
		}
		fmt.Println("Report sent successfully!")
	} else {
		fmt.Println("Report cancelled. Exiting...")
	}
	os.Exit(0)
}
