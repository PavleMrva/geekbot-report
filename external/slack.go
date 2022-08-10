package external

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var slackOauthToken string

func SendSlackRequest(url, method string, payload io.Reader) []byte {
	if slackOauthToken == "" {
		slackOauthToken = os.Getenv("SLACK_OAUTH_TOKEN")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	authToken := fmt.Sprintf("Bearer %s", slackOauthToken)
	req.Header.Add("Authorization", authToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return body
}
