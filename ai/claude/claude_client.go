package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Claude 客户端适配器
type ClaudeClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func New(apiKey, baseURL, model string) *ClaudeClient {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	if model == "" {
		model = "claude-3-sonnet-20240229"
	}

	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 0, // 流式请求不设置超时
		},
	}
}

func (c *ClaudeClient) Name() string {
	return fmt.Sprintf("Claude (%s)", c.model)
}

func (c *ClaudeClient) Chat(prompt string) ([]string, error) {
	// 构建请求体
	requestBody := map[string]interface{}{
		"model":      c.model,
		"max_tokens": 4000,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": true,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", c.baseURL+"/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// 处理流式响应
	scanner := bufio.NewScanner(resp.Body)
	chunks := make([]string, 0, 100)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			// 解析JSON响应
			var response map[string]interface{}
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				continue
			}

			// 检查事件类型
			if eventType, ok := response["type"].(string); ok {
				if eventType == "content_block_delta" {
					if delta, ok := response["delta"].(map[string]interface{}); ok {
						if text, ok := delta["text"].(string); ok {
							chunks = append(chunks, text)
						}
					}
				}
			}
		}
	}

	return chunks, scanner.Err()
}

func (c *ClaudeClient) SimpleChat(ctx context.Context, prompt string) (string, error) {
	var result strings.Builder

	chunks, err := c.Chat(prompt)
	for _, chunk := range chunks {
		result.WriteString(chunk)
	}

	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (c *ClaudeClient) IsAvailable() bool {
	return c.apiKey != ""
}

func (c *ClaudeClient) GetConfig() map[string]interface{} {
	config := make(map[string]interface{})
	config["provider"] = "claude"
	config["model"] = c.model
	config["base_url"] = c.baseURL
	config["api_key_configured"] = c.apiKey != ""
	return config
}
