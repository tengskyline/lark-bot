package qwen

import (
	"context"
	"fmt"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
)

// 千问AI客户端适配器
type QwenClient struct {
	client *dashscopego.TongyiClient
	apiKey string
	model  string
}

func New(apiKey, model string) *QwenClient {
	if model == "" {
		model = qwen.QwenTurbo
	}

	return &QwenClient{
		client: dashscopego.NewTongyiClient(model, apiKey),
		apiKey: apiKey,
		model:  model,
	}
}

func (q *QwenClient) Name() string {
	return fmt.Sprintf("千问AI (%s)", q.model)
}

func (q *QwenClient) Chat(prompt string) ([]string, error) {
	content := qwen.TextContent{Text: prompt}

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: qwen.RoleUser, Content: &content},
		},
	}
	chunks := make([]string, 0, 100)
	// 创建流式回调函数
	streamCallback := func(ctx context.Context, chunk []byte) error {
		if len(chunk) > 0 {
			chunks = append(chunks, string(chunk))
		}
		return nil
	}

	req := &dashscopego.TextRequest{
		Input:       input,
		StreamingFn: streamCallback,
	}

	// 发送请求
	ctx := context.TODO()
	_, err := q.client.CreateCompletion(ctx, req)
	return chunks, err
}

func (q *QwenClient) SimpleChat(ctx context.Context, prompt string) (string, error) {
	content := qwen.TextContent{Text: prompt}

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: qwen.RoleUser, Content: &content},
		},
	}

	req := &dashscopego.TextRequest{
		Input: input,
		// 不设置StreamingFn，获取完整响应
	}

	resp, err := q.client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Output.Choices) > 0 {
		return resp.Output.Choices[0].Message.Content.ToString(), nil
	}

	return "", fmt.Errorf("no response from qwen")
}

func (q *QwenClient) IsAvailable() bool {
	return q.apiKey != "" && q.client != nil
}

func (q *QwenClient) GetConfig() map[string]interface{} {
	config := make(map[string]interface{})
	config["provider"] = "qwen"
	config["model"] = q.model
	config["api_key_configured"] = q.apiKey != ""
	return config
}
