package ai

import (
	"fmt"

	"github.com/tengskyline/lark-bot/ai/claude"
	"github.com/tengskyline/lark-bot/ai/openai"
	"github.com/tengskyline/lark-bot/ai/qwen"
	"github.com/tengskyline/lark-bot/conf"
	"github.com/tengskyline/lark-bot/lark"
)

// AI客户端工厂
type ClientFactory struct{}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

// 根据配置创建AI客户端
func (f *ClientFactory) CreateClient(config conf.AIConfig) (lark.AIClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	switch config.Provider {
	case "qwen":
		return qwen.New(config.APIKey, config.Model), nil
	case "openai":
		return openai.NewOpenAIClient(&openai.OpenAIConfig{
			APIKey:  config.APIKey,
			BaseURL: config.BaseURL,
			Model:   config.Model,
		}), nil
	case "claude":
		return claude.New(config.APIKey, config.BaseURL, config.Model), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", config.Provider)
	}
}

// 获取支持的AI提供商列表
func (f *ClientFactory) GetSupportedProviders() []string {
	return []string{"qwen", "openai", "claude"}
}

// 创建多个AI客户端
func (f *ClientFactory) CreateMultipleClients(configs []conf.AIConfig) (map[string]lark.AIClient, error) {
	clients := make(map[string]lark.AIClient)

	for _, config := range configs {
		client, err := f.CreateClient(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for provider %s: %v", config.Provider, err)
		}
		clients[config.Provider] = client
	}

	return clients, nil
}

// 验证配置
func (f *ClientFactory) ValidateConfig(config conf.AIConfig) error {
	if config.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	supportedProviders := f.GetSupportedProviders()
	for _, provider := range supportedProviders {
		if provider == config.Provider {
			return nil
		}
	}

	return fmt.Errorf("unsupported provider: %s, supported providers: %v", config.Provider, supportedProviders)
}
