# OpenAI 示例

本目录包含 OpenAI 客户端的使用示例。

## 运行示例

### 1. 设置环境变量

```bash
export OPENAI_API_KEY="your-openai-api-key-here"
```

### 2. 运行示例

```bash
cd ai/openai/examples
go run example.go
```

## 示例内容

- **基本使用示例**: 简单的对话测试
- **流式对话示例**: 流式响应处理
- **配置管理示例**: 配置加载和更新
- **Lark Bot 集成示例**: 在 Lark Bot 中集成 OpenAI

## 注意事项

- 确保已设置正确的 OpenAI API 密钥
- 示例中的配置路径可能需要根据实际情况调整
- 某些示例需要 Lark Bot 的相关依赖 