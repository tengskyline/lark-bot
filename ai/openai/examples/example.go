package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tengskyline/lark-bot/ai/openai"
	"github.com/tengskyline/lark-bot/lark"
)

func main() {
	// 示例 1: 基本使用
	basicExample()

	// 示例 2: 流式对话
	streamingExample()

	// 示例 3: 配置管理
	configExample()
}

// 基本使用示例
func basicExample() {
	fmt.Println("=== 基本使用示例 ===")

	config := &openai.OpenAIConfig{
		APIKey:      os.Getenv("OPENAI_API_KEY"),
		Model:       "gpt-4o-mini",
		MaxTokens:   1000,
		Temperature: 0.7,
		Timeout:     30,
		MaxRetries:  3,
	}

	client := openai.NewOpenAIClient(config)

	// 测试连接
	ctx := context.Background()
	if err := client.TestConnection(ctx); err != nil {
		log.Printf("连接测试失败: %v", err)
		return
	}

	// 简单对话
	response, err := client.SimpleChat(ctx, "你好，请介绍一下你自己")
	if err != nil {
		log.Printf("对话失败: %v", err)
		return
	}

	fmt.Printf("AI 回复: %s\n\n", response)
}

// 流式对话示例
func streamingExample() {
	fmt.Println("=== 流式对话示例 ===")

	config := &openai.OpenAIConfig{
		APIKey:      os.Getenv("OPENAI_API_KEY"),
		Model:       "gpt-4o-mini",
		MaxTokens:   2000,
		Temperature: 0.7,
		Timeout:     60,
		MaxRetries:  3,
	}

	client := openai.NewOpenAIClient(config)

	// 流式对话
	chunks, err := client.Chat("请详细解释什么是机器学习，并给出一些实际应用例子")
	if err != nil {
		log.Printf("流式对话失败: %v", err)
		return
	}

	fmt.Println("流式回复:")
	for i, chunk := range chunks {
		fmt.Printf("[%d] %s", i+1, chunk)
	}
	fmt.Println("\n")
}

// 配置管理示例
func configExample() {
	fmt.Println("=== 配置管理示例 ===")

	// 从环境变量加载配置
	config, err := lark.LoadConfigFromEnv()
	if err != nil {
		fmt.Printf("从环境变量加载配置失败: %v\n", err)
		return
	}

	fmt.Printf("配置信息: %+v\n", config.OpenAI)

	// 创建客户端 - 转换配置类型
	openaiConfig := &openai.OpenAIConfig{
		APIKey:      config.OpenAI.APIKey,
		BaseURL:     config.OpenAI.BaseURL,
		Model:       config.OpenAI.Model,
		MaxTokens:   config.OpenAI.MaxTokens,
		Temperature: config.OpenAI.Temperature,
		Timeout:     config.OpenAI.Timeout,
		MaxRetries:  config.OpenAI.MaxRetries,
	}
	client := openai.NewOpenAIClient(openaiConfig)

	// 获取配置信息
	configInfo := client.GetConfig()
	fmt.Printf("客户端配置: %+v\n", configInfo)

	// 更新配置
	newConfig := &openai.OpenAIConfig{
		Model:       "gpt-4o",
		Temperature: 0.5,
		MaxTokens:   3000,
	}

	if err := client.UpdateConfig(newConfig); err != nil {
		fmt.Printf("更新配置失败: %v\n", err)
		return
	}

	fmt.Println("配置更新成功")
}

// 在 Lark Bot 中集成 OpenAI 的示例
func larkBotIntegrationExample() {
	fmt.Println("=== Lark Bot 集成示例 ===")

	// 从环境变量加载配置
	config, err := lark.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建 OpenAI 客户端 - 转换配置类型
	openaiConfig := &openai.OpenAIConfig{
		APIKey:      config.OpenAI.APIKey,
		BaseURL:     config.OpenAI.BaseURL,
		Model:       config.OpenAI.Model,
		MaxTokens:   config.OpenAI.MaxTokens,
		Temperature: config.OpenAI.Temperature,
		Timeout:     config.OpenAI.Timeout,
		MaxRetries:  config.OpenAI.MaxRetries,
	}
	openaiClient := openai.NewOpenAIClient(openaiConfig)

	// 创建处理器配置
	handlerConfig := &lark.HandlerConfig{
		WorkerPoolSize:    4,
		QueueSize:         100,
		MaxConcurrentAI:   5,
		MaxConcurrentHTTP: 10,
		Timeout:           30,
	}

	// 创建处理器
	handler := lark.NewLarkHandler(handlerConfig)

	// 注册 OpenAI 客户端
	handler.RegisterAIClient("openai", openaiClient)
	handler.SetDefaultAI("openai")

	// 列出可用的 AI 客户端
	availableClients := handler.GetAIClient("").Name()
	fmt.Printf("默认 AI 客户端: %s\n", availableClients)

	// 获取默认客户端
	defaultClient := handler.GetAIClient("")
	if defaultClient != nil {
		fmt.Printf("默认 AI 客户端: %s\n", defaultClient.Name())
	}

	fmt.Println("OpenAI 已成功集成到 Lark Bot 中")
}
