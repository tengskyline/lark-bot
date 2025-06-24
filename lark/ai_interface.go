package lark

import (
	"context"
)

// AI 客户端接口
type AIClient interface {
	// 获取AI名称
	Name() string

	// 流式对话
	Chat(prompt string ) ([]string, error)

	// 简单对话（一次性返回完整结果）
	SimpleChat(ctx context.Context, prompt string) (string, error)

	// 检查客户端是否可用
	IsAvailable() bool

	// 获取配置信息（不包含敏感信息）
	GetConfig() map[string]interface{}
}

// AI 响应处理函数类型
type AIResponseHandler func(ctx context.Context, chunk []byte) error

// AI 客户端管理器
type AIClientManager struct {
	clients       map[string]AIClient
	defaultClient string
}

func NewAIClientManager() *AIClientManager {
	return &AIClientManager{
		clients: make(map[string]AIClient),
	}
}

// 注册AI客户端
func (m *AIClientManager) RegisterClient(name string, client AIClient) {
	m.clients[name] = client

	// 如果是第一个注册的客户端，设为默认
	if m.defaultClient == "" {
		m.defaultClient = name
	}
}

// 获取默认AI客户端
func (m *AIClientManager) GetDefault() AIClient {
	if m.defaultClient == "" {
		return nil
	}
	return m.clients[m.defaultClient]
}

// 设置默认AI客户端
func (m *AIClientManager) SetDefault(name string) bool {
	if _, exists := m.clients[name]; exists {
		m.defaultClient = name
		return true
	}
	return false
}

// 获取指定AI客户端
func (m *AIClientManager) GetClient(name string) AIClient {
	return m.clients[name]
}

// 获取所有AI客户端
func (m *AIClientManager) GetAllClients() map[string]AIClient {
	result := make(map[string]AIClient)
	for name, client := range m.clients {
		result[name] = client
	}
	return result
}

// 列出所有可用的AI客户端名称
func (m *AIClientManager) ListAvailableClients() []string {
	var names []string
	for name, client := range m.clients {
		if client.IsAvailable() {
			names = append(names, name)
		}
	}
	return names
}
