package external

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var geekBotApiKey string

func SendGeekBotRequest(url, method string, payload io.Reader) []byte {
	if geekBotApiKey == "" {
		geekBotApiKey = os.Getenv("GEEKBOT_API_KEY")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Authorization", geekBotApiKey)
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
