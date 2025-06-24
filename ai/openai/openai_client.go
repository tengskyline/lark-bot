package openai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAI 客户端配置
type OpenAIConfig struct {
	APIKey      string        `json:"api_key" yaml:"api_key"`
	BaseURL     string        `json:"base_url" yaml:"base_url"`       // 可选，用于自定义 API 端点
	Model       string        `json:"model" yaml:"model"`             // 使用的模型
	MaxTokens   int           `json:"max_tokens" yaml:"max_tokens"`   // 最大 token 数
	Temperature float32       `json:"temperature" yaml:"temperature"` // 温度参数
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`         // 请求超时时间
	MaxRetries  int           `json:"max_retries" yaml:"max_retries"` // 最大重试次数
}

// OpenAI 客户端实现
type OpenAIClient struct {
	config    *OpenAIConfig
	client    *openai.Client
	mutex     sync.RWMutex
	available bool
}

// 创建新的 OpenAI 客户端
func NewOpenAIClient(config *OpenAIConfig) *OpenAIClient {
	if config == nil {
		config = &OpenAIConfig{}
	}

	// 设置默认值
	if config.Model == "" {
		config.Model = openai.GPT4oMini
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	// 创建客户端
	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIClient{
		config:    config,
		client:    client,
		available: config.APIKey != "",
	}
}

// 实现 AIClient 接口
func (c *OpenAIClient) Name() string {
	return "OpenAI"
}

// 流式对话实现
func (c *OpenAIClient) Chat(prompt string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	var chunks []string
	var mu sync.Mutex

	// 创建流式请求
	req := openai.ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Stream:      true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("创建流式请求失败: %w", err)
	}
	defer stream.Close()

	// 处理流式响应
	for {
		response, err := stream.Recv()
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return chunks, fmt.Errorf("接收流式响应失败: %w", err)
		}

		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			content := response.Choices[0].Delta.Content
			mu.Lock()
			chunks = append(chunks, content)
			mu.Unlock()
		}
	}

	return chunks, nil
}

// 简单对话实现（一次性返回完整结果）
func (c *OpenAIClient) SimpleChat(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	// 创建请求
	req := openai.ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	// 带重试的请求
	var lastErr error
	for i := 0; i <= c.config.MaxRetries; i++ {
		resp, err := c.client.CreateChatCompletion(ctx, req)
		if err != nil {
			lastErr = err
			if i < c.config.MaxRetries {
				time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
				continue
			}
			return "", fmt.Errorf("OpenAI 请求失败 (重试 %d 次): %w", i+1, err)
		}

		if len(resp.Choices) > 0 {
			return resp.Choices[0].Message.Content, nil
		}

		return "", fmt.Errorf("OpenAI 返回空响应")
	}

	return "", lastErr
}

// 检查客户端是否可用
func (c *OpenAIClient) IsAvailable() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.available
}

// 获取配置信息（不包含敏感信息）
func (c *OpenAIClient) GetConfig() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return map[string]interface{}{
		"model":       c.config.Model,
		"max_tokens":  c.config.MaxTokens,
		"temperature": c.config.Temperature,
		"timeout":     c.config.Timeout.String(),
		"max_retries": c.config.MaxRetries,
		"base_url":    c.config.BaseURL,
		"available":   c.available,
	}
}

// 更新配置
func (c *OpenAIClient) UpdateConfig(config *OpenAIConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if config.APIKey != "" {
		c.config.APIKey = config.APIKey
		// 重新创建客户端
		clientConfig := openai.DefaultConfig(config.APIKey)
		if config.BaseURL != "" {
			clientConfig.BaseURL = config.BaseURL
		}
		c.client = openai.NewClientWithConfig(clientConfig)
		c.available = true
	}

	if config.Model != "" {
		c.config.Model = config.Model
	}
	if config.MaxTokens > 0 {
		c.config.MaxTokens = config.MaxTokens
	}
	if config.Temperature > 0 {
		c.config.Temperature = config.Temperature
	}
	if config.Timeout > 0 {
		c.config.Timeout = config.Timeout
	}
	if config.MaxRetries > 0 {
		c.config.MaxRetries = config.MaxRetries
	}
	if config.BaseURL != "" {
		c.config.BaseURL = config.BaseURL
	}

	return nil
}

// 测试连接
func (c *OpenAIClient) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model:     c.config.Model,
		Messages:  []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: "Hello"}},
		MaxTokens: 10,
	}

	_, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		c.mutex.Lock()
		c.available = false
		c.mutex.Unlock()
		return fmt.Errorf("连接测试失败: %w", err)
	}

	c.mutex.Lock()
	c.available = true
	c.mutex.Unlock()

	return nil
}
