package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GPT struct {
	apiKey string
}

const (
	model   = "gpt-3.5-turbo"
	apiUrl  = "https://api.openai.com/v1/chat/completions"
	timeOut = 60
)

func New(apiKey string) GPT {
	return GPT{apiKey: apiKey}
}

type gptReq struct {
	Model    string   `json:"model"`
	Messages []gptMsg `json:"messages"`
}

type gptMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type gptResp struct {
	ID      string       `json:"id"`
	Choices []gptChoices `json:"choices"`
}

type gptChoices struct {
	Message gptMsg `json:"message"`
}

func (g *GPT) Handle(message string) (string, error) {
	reqBody := gptReq{
		Model: model,
		Messages: []gptMsg{{
			Role:    "user",
			Content: message,
		}},
	}
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	client := &http.Client{Timeout: timeOut * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gpt api error: %d %s", resp.StatusCode, respBody)
	}

	respData := &gptResp{}
	err = json.Unmarshal(respBody, respData)
	if err != nil {
		return "", err
	}

	if len(respData.Choices) > 0 {
		return respData.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("gpt resp no data: %s", respBody)
}
