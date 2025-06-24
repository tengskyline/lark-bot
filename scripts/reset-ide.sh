#!/bin/bash

# 重置 IDE 状态脚本

echo "🔄 重置 IDE 状态..."

# 清理 Go 缓存
echo "📦 清理 Go 缓存..."
go clean -cache
go clean -modcache

# 重新下载依赖
echo "⬇️  重新下载依赖..."
go mod download

# 同步工作区
echo "🔄 同步工作区..."
go work sync

# 重新构建所有包
echo "🔨 重新构建所有包..."
go build -v ./...

# 重新安装 gopls
echo "🛠️  更新 gopls..."
go install golang.org/x/tools/gopls@latest

# 检查 gopls 版本
echo "📋 gopls 版本信息:"
gopls version

echo "✅ IDE 状态重置完成！"
echo ""
echo "请在 VS Code 中执行以下操作："
echo "1. 按 Ctrl+Shift+P 打开命令面板"
echo "2. 运行 'Go: Restart Language Server'"
echo "3. 运行 'Developer: Reload Window'"
echo ""
echo "如果仍有问题，请重启 VS Code。" 