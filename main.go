package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"geekbot-report/external"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type StandupQuestion struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Standup struct {
	ID        int               `json:"id"`
	Channel   string            `json:"channel"`
	Questions []StandupQuestion `json:"questions"`
}

type Answer struct {
	Text string `json:"text"`
}

type Report struct {
	StandupID string            `json:"standup_id"`
	Answers   map[string]Answer `json:"answers"`
}

func getUniqueIssues(issues []string) []string {
	var uniqueIssues []string

	for _, issue := range issues {
		issueExists := false

		for _, uniqueIssue := range uniqueIssues {
			if uniqueIssue == issue {
				issueExists = true
			}
		}

		if !issueExists {
			uniqueIssues = append(uniqueIssues, issue)
		}
	}

	return uniqueIssues
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
		log.Println("Cannot unmarshal JSON", err.Error())
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

func isDate(stringDate string) bool {
	_, err := time.Parse("2006-01-02", stringDate)
	return err == nil
}

func askForDate(s string) string {
	reader := bufio.NewReader(os.Stdin)
	todaysDate := time.Now().Weekday()

	var defaultDate string

	if todaysDate == time.Monday {
		defaultDate = time.Now().Add(-72 * time.Hour).Format("2006-01-02")
	} else {
		defaultDate = time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	}

	for {
		fmt.Printf("%s\n(format: YYYY-MM-DD, default: last work-day): ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.TrimSpace(response)
		if response == "" {
			return defaultDate
		}

		isResponseDate := isDate(response)
		if isResponseDate {
			return response
		}
	}
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

	selectedDate := askForDate("Choose the date form which you want to generate stand-up report")

	fmt.Printf("\nGenerating report since %s\n", selectedDate)

	issues, err := external.MakeTempoRequest(selectedDate)

	if err != nil {
		log.Panic("Error while sending Tempo request", err)
	}

	if len(issues) == 0 {
		log.Panic("No issues for report")
	}

	uniqueIssues := getUniqueIssues(issues)

	wg := sync.WaitGroup{}
	wg.Add(len(uniqueIssues) + 1)

	var reportIssues sort.StringSlice

	for index, issue := range uniqueIssues {
		go func(index int, issue string) {
			defer wg.Done()

			issueStr, err := external.MakeJiraRequest(issue)
			if err != nil {
				log.Fatal("Error while sending JIRA request")
			}

			reportIssues = append(reportIssues, issueStr)
		}(index, issue)
	}

	var dailyStandup Standup

	go func() {
		defer wg.Done()

		dailyStandup, err = getDailyStandup()

		if err != nil {
			log.Fatal("Error while fetching daily standup", err)
		}
	}()

	wg.Wait()

	reportIssues.Sort()
	report := strings.Join(reportIssues, "\n")

	dailyStandupReport := Report{
		StandupID: fmt.Sprintf("%d", dailyStandup.ID),
		Answers:   make(map[string]Answer),
	}

	for _, question := range dailyStandup.Questions {
		key := fmt.Sprintf("%d", question.ID)

		switch {
		case strings.Contains(question.Text, "What did"):
			dailyStandupReport.Answers[key] = Answer{report}
		case strings.Contains(question.Text, "What will"):
			dailyStandupReport.Answers[key] = Answer{"sprint/support"}
		default:
			dailyStandupReport.Answers[key] = Answer{"-"}
		}
	}

	confirmationQuestion := fmt.Sprintf("\nYour report for today:\n%s\nAre you sure that you want to send this report?", report)
	answer := askForConfirmation(confirmationQuestion)

	if answer {
		_, err = sendGeekBotReport(dailyStandupReport)

		if err != nil {
			log.Fatal("Error while sending daily standup")
		}

		fmt.Println("Report sent successfully!")
	} else {
		fmt.Println("Report cancelled. Exiting...")
	}

	os.Exit(0)
}
