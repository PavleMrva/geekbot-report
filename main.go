package main

import (
	"encoding/json"
	"fmt"
	"geekbot-report/external"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"
)

// if not precompiled, use ENV variables
var geekBotChannel string

type GeekBotQuestion struct {
	Question   string
	Answer     string
	IsAnswered bool
}

type SlackMessage struct {
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Text    string `json:"text"`
}

type GeekBotHistoryResponse struct {
	Messages []SlackMessage `json:"messages"`
}

var geekBotQuestions = []GeekBotQuestion{
	{
		Question:   "What did you do since yesterday?",
		Answer:     "",
		IsAnswered: false,
	},
	{
		Question:   "What will you do today?",
		Answer:     "sprint/support",
		IsAnswered: false,
	},
	{
		Question:   "Anything blocking your progress?",
		Answer:     "no",
		IsAnswered: false,
	},
	{
		Question:   "Do you need assistance from any other teams?",
		Answer:     "no",
		IsAnswered: false,
	},
	{
		Question:   "Anything miscellaneous?",
		Answer:     "no",
		IsAnswered: false,
	},
}

func sendGeekBotMessage(message string) {
	url := "https://slack.com/api/chat.postMessage"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
	  "channel": "%s",
	  "text": "%s"
	}`, geekBotChannel, message))

	external.SendSlackRequest(url, method, payload)
}

func getGeekBotLastMessage() GeekBotHistoryResponse {
	url := fmt.Sprintf("https://slack.com/api/conversations.history?channel=%s&limit=%d", geekBotChannel, 1)
	method := "GET"

	body := external.SendSlackRequest(url, method, nil)

	var response GeekBotHistoryResponse
	if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Cannot unmarshal JSON")
	}
	return response
}

func startGeekBotReport(report string) {
	lastMessage := getGeekBotLastMessage()

	// If the geekbot report has not yet been initialized
	if !strings.Contains(lastMessage.Messages[0].Text, geekBotQuestions[0].Question) {
		sendGeekBotMessage("report")
		time.Sleep(5 * time.Second)
	}

	for _, geekBotQuestion := range geekBotQuestions {
		fmt.Printf("\nGeekbot question: %s", geekBotQuestion.Question)
		for !geekBotQuestion.IsAnswered {
			messaage := getGeekBotLastMessage()
			fmt.Printf("\nRetrieved message: %s", messaage)

			if strings.Contains(messaage.Messages[0].Text, geekBotQuestion.Question) {
				if geekBotQuestion.Answer == "" {
					geekBotQuestion.Answer = report
				}

				sendGeekBotMessage(geekBotQuestion.Answer)
				geekBotQuestion.IsAnswered = true
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	_ = godotenv.Load()
	if geekBotChannel == "" {
		geekBotChannel = os.Getenv("GEEKBOT_CHANNEL_ID")
	}
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
	startGeekBotReport(report)
}
