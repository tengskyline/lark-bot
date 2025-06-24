# 配置文件说明

本目录包含 Lark Bot 的配置文件，采用分层配置管理策略。

## 📁 文件结构

```
conf/
├── config.yaml.local    # 示例配置文件（提交到Git）
├── config.yaml          # 启动配置文件（不提交到Git）
├── config.yaml.prod     # 生产环境配置（不提交到Git）
├── config.yaml.staging  # 测试环境配置（不提交到Git）
└── README.md            # 本说明文件
```

## 🔧 配置文件说明

### config.yaml.local
- **用途**: 示例配置文件，包含所有配置项的说明和示例值
- **版本控制**: ✅ 提交到Git
- **内容**: 脱敏的示例值和详细注释
- **作用**: 作为配置模板，供开发者参考

### config.yaml
- **用途**: 启动配置文件，包含实际的配置值
- **版本控制**: ❌ 不提交到Git（已在.gitignore中）
- **内容**: 实际的配置值，包含敏感信息
- **作用**: 程序启动时使用的实际配置

### config.yaml.prod
- **用途**: 生产环境配置
- **版本控制**: ❌ 不提交到Git（已在.gitignore中）
- **内容**: 生产环境的实际配置值
- **作用**: 生产部署时使用

## 🚀 快速开始

### 1. 创建启动配置

```bash
# 复制示例配置
cp config.yaml.local config.yaml

# 编辑配置文件，填入实际值
vim config.yaml
```

### 2. 必填配置项

确保以下配置项已正确填写：

```yaml
# 飞书应用配置（必填）
AppId: "cli_your_actual_app_id"
AppSecret: "your_actual_app_secret"
VerificationToken: "your_verification_token"
EncryptKey: "your_encrypt_key"

# AI配置（必填）
AI:
  Provider: "qwen"  # 或 openai, claude
  APIKey: "your_actual_api_key"
```

### 3. 运行程序

```bash
# 使用启动配置运行（默认）
go run .

# 或指定配置文件
go run . -c conf/config.yaml
```

## 🔐 安全注意事项

1. **敏感信息保护**
   - `config.yaml` 包含实际配置值，已添加到`.gitignore`
   - 不要将包含敏感信息的配置文件提交到版本控制

2. **生产环境部署**
   - 使用环境变量或安全的配置管理工具
   - 定期轮换API密钥和Token

3. **配置验证**
   - 程序启动时会验证必要的配置项
   - 缺少必填配置会导致启动失败

## 🌍 环境变量支持

所有配置项都支持环境变量覆盖，格式为 `LARK_BOT_` + 大写配置项名：

```bash
# 示例
export LARK_BOT_APPID="cli_your_app_id"
export LARK_BOT_APP_SECRET="your_app_secret"
export LARK_BOT_AI_APIKEY="your_ai_api_key"
export LARK_BOT_AI_PROVIDER="qwen"
```

环境变量的优先级高于配置文件，适合在容器化部署中使用。

## 📝 配置项说明

详细的配置项说明请参考 `config.yaml.local` 文件中的注释。

## 🆘 故障排除

### 配置文件找不到
```bash
Error: config file not found: conf/config.yaml
```
**解决方案**: 复制示例配置文件 `cp config.yaml.local config.yaml`

### 配置验证失败
```bash
Error: required config item AppId is not set
```
**解决方案**: 检查配置文件中的必填项是否已正确填写

### 权限问题
```bash
Error: permission denied reading config file
```
**解决方案**: 检查配置文件的读取权限 `chmod 644 config.yaml` 