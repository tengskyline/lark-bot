package lark

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/tengskyline/lark-bot/ai/openai"
)

// 配置结构体
type Config struct {
	OpenAI       openai.OpenAIConfig `mapstructure:"openai"`
	SystemPrompt string              `mapstructure:"system_prompt"`
	Conversation ConversationConfig  `mapstructure:"conversation"`
}

// 对话配置
type ConversationConfig struct {
	MaxHistory    int  `mapstructure:"max_history"`
	EnableContext bool `mapstructure:"enable_context"`
}

// 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// 验证配置
func validateConfig(config *Config) error {
	if config.OpenAI.APIKey == "" || config.OpenAI.APIKey == "your-openai-api-key-here" {
		return fmt.Errorf("OpenAI API 密钥未设置")
	}

	if config.OpenAI.Model == "" {
		config.OpenAI.Model = "gpt-4o-mini"
	}

	if config.OpenAI.MaxTokens <= 0 {
		config.OpenAI.MaxTokens = 4000
	}

	if config.OpenAI.Temperature < 0 || config.OpenAI.Temperature > 2 {
		return fmt.Errorf("温度参数必须在 0-2 之间")
	}

	if config.OpenAI.Timeout <= 0 {
		config.OpenAI.Timeout = 60 * time.Second
	}

	if config.OpenAI.MaxRetries < 0 {
		config.OpenAI.MaxRetries = 3
	}

	if config.Conversation.MaxHistory <= 0 {
		config.Conversation.MaxHistory = 10
	}

	return nil
}

// 从环境变量加载配置
func LoadConfigFromEnv() (*Config, error) {
	viper.SetEnvPrefix("LARK_BOT")
	viper.AutomaticEnv()

	config := &Config{
		OpenAI: openai.OpenAIConfig{
			APIKey:      viper.GetString("OPENAI_API_KEY"),
			BaseURL:     viper.GetString("OPENAI_BASE_URL"),
			Model:       viper.GetString("OPENAI_MODEL"),
			MaxTokens:   viper.GetInt("OPENAI_MAX_TOKENS"),
			Temperature: float32(viper.GetFloat64("OPENAI_TEMPERATURE")),
			Timeout:     viper.GetDuration("OPENAI_TIMEOUT"),
			MaxRetries:  viper.GetInt("OPENAI_MAX_RETRIES"),
		},
		SystemPrompt: viper.GetString("SYSTEM_PROMPT"),
		Conversation: ConversationConfig{
			MaxHistory:    viper.GetInt("CONVERSATION_MAX_HISTORY"),
			EnableContext: viper.GetBool("CONVERSATION_ENABLE_CONTEXT"),
		},
	}

	// 设置默认值
	if config.OpenAI.Model == "" {
		config.OpenAI.Model = "gpt-4o-mini"
	}
	if config.OpenAI.MaxTokens == 0 {
		config.OpenAI.MaxTokens = 4000
	}
	if config.OpenAI.Temperature == 0 {
		config.OpenAI.Temperature = 0.7
	}
	if config.OpenAI.Timeout == 0 {
		config.OpenAI.Timeout = 60 * time.Second
	}
	if config.OpenAI.MaxRetries == 0 {
		config.OpenAI.MaxRetries = 3
	}
	if config.Conversation.MaxHistory == 0 {
		config.Conversation.MaxHistory = 10
	}

	return config, validateConfig(config)
}
