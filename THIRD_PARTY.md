# 第三方库说明

本文档详细列出了 Lark Bot 项目使用的所有第三方库及其用途。

## 📦 直接依赖

### AI服务集成

#### github.com/devinyf/dashscopego v0.1.1
- **用途**: 阿里云通义千问Go SDK
- **功能**: 提供与通义千问AI服务的完整集成
- **许可证**: MIT
- **项目地址**: https://github.com/devinyf/dashscopego
- **使用位置**: `ai/qwen/qwen_client.go`

#### github.com/larksuite/oapi-sdk-go/v3 v3.4.19
- **用途**: 飞书官方Go SDK
- **功能**: 提供飞书开放平台API的完整封装
- **许可证**: Apache-2.0
- **项目地址**: https://github.com/larksuite/oapi-sdk-go
- **使用位置**: `lark/lark.go`, `lark/safe_bot.go`

### 配置管理

#### github.com/spf13/viper v1.20.1
- **用途**: 配置管理库
- **功能**: 支持多种配置格式（YAML、JSON、TOML等）和环境变量
- **许可证**: MIT
- **项目地址**: https://github.com/spf13/viper
- **使用位置**: `conf/conf.go`

#### github.com/go-viper/mapstructure/v2 v2.2.1
- **用途**: 结构体映射工具
- **功能**: 将配置数据映射到Go结构体
- **许可证**: MIT
- **项目地址**: https://github.com/go-viper/mapstructure
- **使用位置**: `conf/conf.go`

### 系统工具

#### github.com/fsnotify/fsnotify v1.8.0
- **用途**: 文件系统监控
- **功能**: 监控文件系统变化，支持配置热重载
- **许可证**: BSD-3-Clause
- **项目地址**: https://github.com/fsnotify/fsnotify
- **使用位置**: `conf/conf.go`

#### github.com/google/uuid v1.6.0
- **用途**: UUID生成库
- **功能**: 生成符合RFC 4122标准的UUID
- **许可证**: Apache-2.0
- **项目地址**: https://github.com/google/uuid
- **使用位置**: 项目内部ID生成

#### github.com/gorilla/websocket v1.5.3
- **用途**: WebSocket实现
- **功能**: 提供WebSocket客户端和服务器实现
- **许可证**: BSD-2-Clause
- **项目地址**: https://github.com/gorilla/websocket
- **使用位置**: 飞书SDK内部使用

## 🔗 间接依赖

### 数据处理

#### github.com/gabriel-vasile/mimetype v1.4.9
- **用途**: MIME类型检测
- **功能**: 基于文件内容检测MIME类型
- **许可证**: MIT
- **项目地址**: https://github.com/gabriel-vasile/mimetype

#### github.com/pelletier/go-toml/v2 v2.2.3
- **用途**: TOML配置文件解析
- **功能**: 解析TOML格式的配置文件
- **许可证**: MIT
- **项目地址**: https://github.com/pelletier/go-toml

#### gopkg.in/yaml.v3 v3.0.1
- **用途**: YAML配置文件解析
- **功能**: 解析YAML格式的配置文件
- **许可证**: Apache-2.0
- **项目地址**: https://gopkg.in/yaml.v3

### 并发和工具

#### go.uber.org/atomic v1.9.0
- **用途**: 原子操作库
- **功能**: 提供原子操作和并发安全的数据类型
- **许可证**: MIT
- **项目地址**: https://github.com/uber-go/atomic

#### go.uber.org/mock v0.5.2
- **用途**: 测试模拟库
- **功能**: 生成测试用的模拟对象
- **许可证**: MIT
- **项目地址**: https://github.com/uber-go/mock

#### go.uber.org/multierr v1.9.0
- **用途**: 多错误处理
- **功能**: 组合多个错误为一个错误
- **许可证**: MIT
- **项目地址**: https://github.com/uber-go/multierr

### 网络和系统

#### golang.org/x/net v0.41.0
- **用途**: 网络扩展包
- **功能**: 提供额外的网络功能
- **许可证**: BSD-3-Clause
- **项目地址**: https://golang.org/x/net

#### golang.org/x/sys v0.33.0
- **用途**: 系统调用接口
- **功能**: 提供系统级API的跨平台接口
- **许可证**: BSD-3-Clause
- **项目地址**: https://golang.org/x/sys

#### golang.org/x/text v0.26.0
- **用途**: 文本处理工具
- **功能**: 提供文本处理、编码转换等功能
- **许可证**: BSD-3-Clause
- **项目地址**: https://golang.org/x/text

### 配置和文件系统

#### github.com/sagikazarmark/locafero v0.7.0
- **用途**: 本地文件系统适配器
- **功能**: 为Viper提供本地文件系统支持
- **许可证**: MIT
- **项目地址**: https://github.com/sagikazarmark/locafero

#### github.com/spf13/afero v1.12.0
- **用途**: 文件系统抽象层
- **功能**: 提供文件系统操作的抽象接口
- **许可证**: Apache-2.0
- **项目地址**: https://github.com/spf13/afero

#### github.com/spf13/cast v1.7.1
- **用途**: 类型转换工具
- **功能**: 安全地转换接口类型
- **许可证**: MIT
- **项目地址**: https://github.com/spf13/cast

#### github.com/spf13/pflag v1.0.6
- **用途**: 命令行标志解析
- **功能**: 解析命令行参数和标志
- **许可证**: BSD-3-Clause
- **项目地址**: https://github.com/spf13/pflag

#### github.com/subosito/gotenv v1.6.0
- **用途**: 环境变量加载
- **功能**: 从.env文件加载环境变量
- **许可证**: MIT
- **项目地址**: https://github.com/subosito/gotenv

### 错误处理和日志

#### github.com/sourcegraph/conc v0.3.0
- **用途**: 并发工具
- **功能**: 提供并发编程的工具和模式
- **许可证**: MIT
- **项目地址**: https://github.com/sourcegraph/conc

#### github.com/gogo/protobuf v1.3.2
- **用途**: Protocol Buffers实现
- **功能**: 高性能的Protocol Buffers实现
- **许可证**: BSD-3-Clause
- **项目地址**: https://github.com/gogo/protobuf

## 📋 许可证合规性

本项目使用的所有第三方库都采用开源许可证，主要包括：

- **MIT License**: 最宽松的许可证之一
- **Apache-2.0**: Apache软件基金会许可证
- **BSD-2-Clause**: 简化的BSD许可证
- **BSD-3-Clause**: 标准BSD许可证

## 🔍 安全审计

建议定期进行以下安全审计：

1. **依赖更新**: 定期更新依赖到最新版本
2. **漏洞扫描**: 使用工具扫描已知漏洞
3. **许可证检查**: 确保许可证兼容性

```bash
# 检查依赖更新
go list -u -m all

# 更新依赖
go get -u ./...

# 清理未使用的依赖
go mod tidy
```

## 📚 相关资源

- [Go Modules 官方文档](https://golang.org/ref/mod)
- [Go 依赖管理最佳实践](https://golang.org/doc/modules/managing-dependencies)
- [开源许可证选择指南](https://choosealicense.com/)

---

**注意**: 本文档会随着项目依赖的变化而更新，请定期检查最新版本。 