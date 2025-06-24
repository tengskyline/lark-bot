# OpenAI 集成指南

本项目已集成 [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) 库，支持与 OpenAI API 进行交互。

## 功能特性

- ✅ 支持流式对话（Streaming）
- ✅ 支持简单对话（一次性返回）
- ✅ 自动重试机制
- ✅ 连接池和并发控制
- ✅ 配置热更新
- ✅ 多模型支持
- ✅ 自定义 API 端点

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置 OpenAI

#### 方法一：配置文件

创建 `config/openai.yaml`：

```yaml
openai:
  api_key: "your-openai-api-key-here"
  model: "gpt-4o-mini"
  max_tokens: 4000
  temperature: 0.7
  timeout: "60s"
  max_retries: 3
```

#### 方法二：环境变量

```bash
export LARK_BOT_OPENAI_API_KEY="your-openai-api-key-here"
export LARK_BOT_OPENAI_MODEL="gpt-4o-mini"
export LARK_BOT_OPENAI_MAX_TOKENS="4000"
export LARK_BOT_OPENAI_TEMPERATURE="0.7"
```

### 3. 基本使用

```go
package main

import (
    "context"
    "fmt"
    openai "github.com/tengskyline/lark-bot/ai/openai"
)

func main() {
    // 创建 OpenAI 客户端
    config := &openai.OpenAIConfig{
        APIKey:      "your-api-key",
        Model:       "gpt-4o-mini",
        MaxTokens:   1000,
        Temperature: 0.7,
    }
    
    client := openai.NewOpenAIClient(config)
    
    // 简单对话
    ctx := context.Background()
    response, err := client.SimpleChat(ctx, "你好")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(response)
}
```

### 4. 流式对话

```go
// 流式对话
chunks, err := client.Chat("请解释什么是人工智能")
if err != nil {
    panic(err)
}

for _, chunk := range chunks {
    fmt.Print(chunk)
}
```

## 在 Lark Bot 中集成

### 1. 注册 OpenAI 客户端

```go
// 创建处理器
handler := lark.NewLarkHandler(&lark.HandlerConfig{
    WorkerPoolSize:    4,
    QueueSize:         100,
    MaxConcurrentAI:   5,
    MaxConcurrentHTTP: 10,
})

// 创建 OpenAI 客户端
openaiClient := openai.NewOpenAIClient(&config.OpenAI)

// 注册为默认 AI 客户端
handler.RegisterAIClient("openai", openaiClient)
handler.SetDefaultAI("openai")
```

### 2. 自动处理消息

当用户发送消息时，系统会自动：

1. 创建交互式卡片
2. 调用 OpenAI API
3. 流式更新卡片内容
4. 处理错误和重试

## 配置选项

### OpenAI 配置

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `api_key` | string | - | OpenAI API 密钥（必需） |
| `base_url` | string | "" | 自定义 API 端点 |
| `model` | string | "gpt-4o-mini" | 使用的模型 |
| `max_tokens` | int | 4000 | 最大 token 数 |
| `temperature` | float32 | 0.7 | 温度参数 (0.0-2.0) |
| `timeout` | duration | 60s | 请求超时时间 |
| `max_retries` | int | 3 | 最大重试次数 |

### 支持的模型

- `gpt-4o` - GPT-4 Omni
- `gpt-4o-mini` - GPT-4 Omni Mini
- `gpt-4-turbo` - GPT-4 Turbo
- `gpt-3.5-turbo` - GPT-3.5 Turbo

## 高级功能

### 1. 连接测试

```go
err := client.TestConnection(ctx)
if err != nil {
    log.Printf("连接失败: %v", err)
}
```

### 2. 配置热更新

```go
newConfig := &openai.OpenAIConfig{
    Model:       "gpt-4o",
    Temperature: 0.5,
    MaxTokens:   3000,
}

err := client.UpdateConfig(newConfig)
```

### 3. 获取配置信息

```go
configInfo := client.GetConfig()
fmt.Printf("配置: %+v\n", configInfo)
```

### 4. 多客户端支持

```go
// 注册多个 AI 客户端
handler.RegisterAIClient("openai", openaiClient)
handler.RegisterAIClient("openai-gpt4", gpt4Client)

// 切换默认客户端
handler.SetDefaultAI("openai-gpt4")
```

## 错误处理

### 常见错误

1. **API 密钥无效**
   ```
   OpenAI 请求失败: openai: invalid api key
   ```

2. **网络超时**
   ```
   OpenAI 请求失败: context deadline exceeded
   ```

3. **模型不存在**
   ```
   OpenAI 请求失败: openai: model not found
   ```

### 重试机制

系统会自动重试失败的请求，重试间隔为指数退避：

- 第1次重试：1秒后
- 第2次重试：2秒后
- 第3次重试：3秒后

## 性能优化

### 1. 并发控制

```go
handlerConfig := &lark.HandlerConfig{
    MaxConcurrentAI: 5,  // 限制 AI 并发数
}
```

### 2. 连接池

系统使用 HTTP 客户端池来复用连接，提高性能。

### 3. 流式处理

使用流式 API 可以实时显示 AI 回复，提升用户体验。

## 监控和日志

### 1. 性能监控

```go
// 启动监控报告器
go handler.StartMetricsReporter()
```

### 2. 日志输出

系统会输出详细的日志信息：

```
使用AI客户端: OpenAI
AI对话失败: context deadline exceeded
Handler Stats - Processed: 100, Errors: 5, Queued: 10, Dropped: 0
```

## 示例代码

完整示例请参考 `examples/openai_example.go`。

## 故障排除

### 1. 依赖问题

```bash
# 清理并重新下载依赖
go clean -modcache
go mod tidy
```

### 2. 配置问题

```bash
# 验证配置文件
go run examples/openai_example.go
```

### 3. 网络问题

```bash
# 测试网络连接
curl -H "Authorization: Bearer YOUR_API_KEY" \
     https://api.openai.com/v1/models
```

## 相关链接

- [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) - Go OpenAI 库
- [OpenAI API 文档](https://platform.openai.com/docs) - 官方 API 文档
- [Lark Bot 文档](./README.md) - 项目主文档 