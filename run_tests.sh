#!/bin/bash

# 数据库应用 - 后端测试脚本
# 使用临时 Go 容器运行测试

echo "=============================================="
echo "       数据库应用 - 后端单元测试"
echo "=============================================="
echo ""
echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"

if [ ! -d "$BACKEND_DIR" ]; then
    echo "❌ 错误: 找不到 backend 目录"
    exit 1
fi

echo "📁 项目目录: $BACKEND_DIR"
echo "🐳 启动临时 Go 测试容器..."
echo ""

# 运行测试
echo "=============================================="
echo "                 测试结果"
echo "=============================================="
echo ""

# 使用临时 Go 容器运行测试（使用 Debian 镜像支持 CGO/SQLite）
docker run --rm \
    -v "$BACKEND_DIR":/app \
    -w /app \
    -e GOPROXY=https://goproxy.cn,direct \
    -e CGO_ENABLED=1 \
    m.daocloud.io/docker.io/library/golang:1.21 \
    go test ./... -v -cover 2>&1 | tee /tmp/test_output.txt

TEST_EXIT_CODE=${PIPESTATUS[0]}

echo ""
echo "=============================================="
echo "                 测试摘要"
echo "=============================================="
echo ""

# 统计测试结果
PASS_COUNT=$(grep -c "^--- PASS" /tmp/test_output.txt 2>/dev/null || true)
FAIL_COUNT=$(grep -c "^--- FAIL" /tmp/test_output.txt 2>/dev/null || true)
SKIP_COUNT=$(grep -c "^--- SKIP" /tmp/test_output.txt 2>/dev/null || true)

# 确保是数字
PASS_COUNT=${PASS_COUNT:-0}
FAIL_COUNT=${FAIL_COUNT:-0}
SKIP_COUNT=${SKIP_COUNT:-0}

echo "✅ 通过: $PASS_COUNT"
echo "❌ 失败: $FAIL_COUNT"
echo "⏭️  跳过: $SKIP_COUNT"
echo ""

# 显示覆盖率信息
echo "📊 测试覆盖率:"
grep "coverage:" /tmp/test_output.txt 2>/dev/null || echo "   (无覆盖率数据)"
echo ""

echo "结束时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# 根据测试结果返回状态码
if [ "$TEST_EXIT_CODE" -ne 0 ]; then
    echo "⚠️  有测试失败，请检查上面的输出"
    exit 1
else
    echo "🎉 所有测试通过！"
    exit 0
fi
