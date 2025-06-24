# 🤖 Lark Bot

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

一个基于 Go 语言开发的高性能飞书机器人，支持多 AI 平台、高并发处理、插件化架构。

## ✨ 功能特性

- 🤖 **多 AI 平台**: 集成 OpenAI、Qwen、Claude 等多种 AI 服务
- 🔄 **高并发处理**: 支持大量消息并发处理，性能优异
- 🏭 **AI 工厂模式**: 统一的 AI 客户端工厂，支持热插拔
- 🔌 **插件系统**: 可扩展的插件架构，支持自定义功能
- 📊 **性能监控**: 内置性能监控和指标收集
- 🔐 **安全配置**: 支持环境变量和加密配置
- 🐳 **容器化部署**: 提供 Docker 和 Docker Compose 支持

## 🚀 快速开始

### 环境要求

- **Go**: 1.24+
- **飞书应用**: 已创建的飞书机器人应用
- **AI API**: 至少一个AI提供商的API密钥

### 安装部署

#### 1. 克隆项目
```bash
git clone https://github.com/tengskyline/lark-bot.git
cd lark-bot
```

#### 2. 配置 AI 服务
```bash
# 设置环境变量
export LARK_BOT_OPENAI_API_KEY="your-openai-api-key-here"
export LARK_BOT_OPENAI_MODEL="gpt-4o-mini"
```

#### 3. 运行程序
```bash
go run main.go
```

## 📁 项目结构

```
lark-bot/
├── main.go                 # 程序入口
├── go.mod                  # 模块依赖
├── go.work                 # 工作区配置
├── ai/                     # AI 客户端工厂与适配器
│   ├── factory.go         # AI 客户端统一工厂
│   ├── openai/            # OpenAI 客户端实现
│   │   ├── openai_client.go
│   │   └── examples/      # OpenAI 使用示例
│   ├── qwen/              # Qwen 通义千问适配
│   └── claude/            # Claude 适配
├── lark/                   # Lark 机器人核心逻辑
│   ├── handler.go         # 高并发消息处理器
│   ├── ai_interface.go    # AI 客户端接口
│   ├── config_loader.go   # 配置加载器
│   └── plugin_*.go        # 插件系统
├── conf/                   # 配置结构体定义
├── config/                 # 配置文件目录
└── plugins/                # 插件实现
    └── jenkins/           # Jenkins 插件
```

## 🤖 AI 管理与适配

### 统一 AI 客户端工厂

通过 `ai/factory.go` 的 `ClientFactory`，可根据配置动态创建和管理多种 AI 客户端：

```go
factory := ai.NewClientFactory()
client, err := factory.CreateClient(conf.AIConfig{
    Provider: "openai",
    APIKey:   "sk-xxx",
    Model:    "gpt-4o-mini",
})
```

### 支持的 AI 提供商

- **OpenAI** - GPT 系列模型（gpt-4o, gpt-4o-mini 等）
- **Qwen** - 通义千问系列模型
- **Claude** - Anthropic Claude 系列模型

### OpenAI 适配

新版 OpenAI 客户端已迁移到 `ai/openai`，使用官方 go-openai SDK，支持：

- ✅ 流式对话（Streaming）
- ✅ 自动重试机制
- ✅ 超时控制
- ✅ 配置热更新
- ✅ 线程安全

创建方式：
```go
client := openai.NewOpenAIClient(&openai.OpenAIConfig{
    APIKey:      "sk-xxx",
    BaseURL:     "",           // 可选，自定义 API 端点
    Model:       "gpt-4o-mini",
    MaxTokens:   4000,
    Temperature: 0.7,
    Timeout:     60 * time.Second,
    MaxRetries:  3,
})
```

## ⚙️ 配置方法

### 1. 环境变量（推荐）

```bash
export LARK_BOT_OPENAI_API_KEY="your-openai-api-key-here"
export LARK_BOT_OPENAI_MODEL="gpt-4o-mini"
export LARK_BOT_OPENAI_MAX_TOKENS="4000"
export LARK_BOT_OPENAI_TEMPERATURE="0.7"
export LARK_BOT_OPENAI_TIMEOUT="60s"
export LARK_BOT_OPENAI_MAX_RETRIES="3"
```

### 2. 主配置文件

在 `conf/config.yaml` 中配置：

```yaml
AI:
  Provider: "openai"
  APIKey: "your-openai-api-key-here"
  BaseURL: ""                    # 可选，用于代理或自定义服务
  Model: "gpt-4o-mini"
  Options:
    max_tokens: "4000"
    temperature: "0.7"
    timeout: "60s"
    max_retries: "3"
```

## 🏃‍♂️ 运行与示例

### 运行主程序

```bash
go run main.go
```

### 运行 OpenAI 示例

```bash
cd ai/openai/examples
go run example.go
```

### 测试 AI 连接

```go
ctx := context.Background()
err := client.TestConnection(ctx)
if err != nil {
    log.Printf("连接失败: %v", err)
}
```

## 🔧 扩展性

### 添加新的 AI 提供商

1. 在 `ai/` 目录下创建新的包（如 `ai/gemini/`）
2. 实现 `lark.AIClient` 接口
3. 在 `ai/factory.go` 中注册新的提供商

```go
// 实现 AIClient 接口
type GeminiClient struct {
    // ...
}

func (c *GeminiClient) Name() string { return "Gemini" }
func (c *GeminiClient) Chat(prompt string) ([]string, error) { /* ... */ }
func (c *GeminiClient) SimpleChat(ctx context.Context, prompt string) (string, error) { /* ... */ }
func (c *GeminiClient) IsAvailable() bool { /* ... */ }
func (c *GeminiClient) GetConfig() map[string]interface{} { /* ... */ }

// 在工厂中注册
case "gemini":
    return gemini.New(config.APIKey, config.Model), nil
```

### 插件开发

支持自定义消息处理插件，实现 `lark.JobPlugin` 接口即可。

## 📊 性能特性

- **高并发处理**: 工作协程池 + 消息队列
- **AI 限流**: 信号量控制 AI 并发数
- **HTTP 连接池**: 复用 HTTP 连接
- **流式响应**: 实时显示 AI 回复
- **优雅关闭**: 支持优雅停机

## 🔍 监控与日志

### 性能监控

```go
// 启动监控报告器
go handler.StartMetricsReporter()

// 输出示例
// Handler Stats - Processed: 100, Errors: 5, Queued: 10, Dropped: 0
```

### 日志输出

```
使用AI客户端: OpenAI
AI对话失败: context deadline exceeded
Handler Stats - Processed: 100, Errors: 5, Queued: 10, Dropped: 0
```

## 📚 相关文档

- [OpenAI 集成指南](./OPENAI_INTEGRATION.md)
- [AI 工厂与管理器](./ai/factory.go)
- [Lark Bot 主流程](./lark/handler.go)
- [配置加载与校验](./lark/config_loader.go)
- [OpenAI 示例](./ai/openai/examples/)

## ❓ 常见问题

### 如何切换 AI 平台？

修改配置文件中的 `provider` 字段：
```yaml
AI:
  Provider: "openai"  # 或 "qwen", "claude"
```

### 如何自定义 AI 参数？

在配置文件中补充对应字段：
```yaml
openai:
  temperature: 0.5
  max_tokens: 3000
  timeout: "120s"
```

### 如何添加新的 AI 提供商？

1. 实现 `lark.AIClient` 接口
2. 在 `ai/factory.go` 中注册
3. 更新支持的提供商列表

### 如何处理 API 限流？

- 调整 `MaxConcurrentAI` 参数
- 启用自动重试机制
- 使用指数退避策略

## 📄 License

MIT License - 详见 [LICENSE](LICENSE) 文件



