package external

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var geekBotAPIKey string

func SendGeekBotRequest(url, method string, payload io.Reader) ([]byte, error) {
	if geekBotAPIKey == "" {
		geekBotAPIKey = os.Getenv("GEEKBOT_API_KEY")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Authorization", geekBotAPIKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}
