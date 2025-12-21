#!/bin/bash

# 运行后端服务脚本

set -e

echo "正在启动后端服务..."

# 检查go.mod是否存在
if [ ! -f "go.mod" ]; then
    echo "错误: 未找到go.mod文件，请先运行 'go mod init'"
    exit 1
fi

# 下载依赖
echo "正在下载依赖..."
go mod download

# 运行服务
echo "正在启动服务..."
go run cmd/api/main.go

