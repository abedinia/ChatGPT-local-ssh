package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestData struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ResponseData struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func send_message_to_chatgpt(prompt, model, apiKey string) (*ResponseData, error) {
	url := "https://api.openai.com/v1/chat/completions"
	requestData := RequestData{
		Model:    model,
		Messages: []Message{{Role: "user", Content: prompt}},
	}

	jsonReq, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a prompt as a command-line argument.")
		os.Exit(1)
	}

	prompt := os.Args[1]
	apiKey := "API_KEY"
	model := "gpt-3.5-turbo"

	response, err := send_message_to_chatgpt(prompt, model, apiKey)
	if err != nil {
		fmt.Println("error has been gotten")
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		fmt.Println(response.Choices[0].Message.Content)
	} else {
		fmt.Println("No response or invalid response structure.")
	}
}
