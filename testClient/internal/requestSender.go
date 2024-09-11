package internal

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendPostRequest(url string) {
	fullUrl := fmt.Sprintf("%s%s", PostEndpoint, url)
	jsonStr := fmt.Sprintf(`{"url":"%s"}`, url)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)
}
