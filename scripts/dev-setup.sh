#!/bin/bash

# 开发环境设置脚本
# 用于新开发者快速设置本地开发环境

set -e

echo "🚀 设置 Lark-Bot 开发环境..."

# 检查 Go 版本
echo "📋 检查 Go 版本..."
go version

# 检查是否在正确的目录
if [ ! -f "go.mod" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

# 创建工作区（如果不存在）
if [ ! -f "go.work" ]; then
    echo "📝 创建 Go 工作区..."
    go work init .
    echo "✅ 工作区创建完成"
else
    echo "✅ 工作区已存在"
fi

# 设置配置文件
if [ ! -f "conf/config.yaml" ]; then
    echo "📋 创建启动配置文件..."
    cp conf/config.yaml.local conf/config.yaml
    echo "✅ 已创建 conf/config.yaml"
    echo "⚠️  请编辑 conf/config.yaml 并填入实际的配置值"
else
    echo "✅ 启动配置文件已存在"
fi

# 下载依赖
echo "📦 下载依赖..."
go mod download
go work sync

# 验证构建
echo "🔨 验证构建..."
go build .

# 运行测试
echo "🧪 运行测试..."
go test ./...

echo "🎉 开发环境设置完成！"
echo ""
echo "📚 常用命令："
echo "  go run .                              # 使用启动配置运行"
echo "  go run . -c conf/config.yaml          # 指定配置文件运行"
echo "  go build .                            # 构建程序" 
echo "  go test ./...                         # 运行测试"
echo "  go work sync                          # 同步工作区"
echo ""
echo "📝 配置相关："
echo "  vim conf/config.yaml                  # 编辑启动配置"
echo "  cat conf/config.yaml.local            # 查看示例配置"
echo ""
echo "💡 提示："
echo "  - go.work 文件仅用于本地开发，不会提交到版本控制"
echo "  - config.yaml 包含敏感信息，不会提交到版本控制"
echo "  - 请确保填写正确的飞书应用配置和AI API密钥" 