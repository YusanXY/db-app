.PHONY: test test-cover test-verbose test-race

# 运行所有测试
test:
	go test ./... -v

# 运行测试并生成覆盖率报告
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 详细输出测试
test-verbose:
	go test ./... -v

# 竞态检测
test-race:
	go test ./... -race

# 运行特定包的测试
test-pkg:
	@read -p "请输入包路径 (例如: ./internal/service): " pkg; \
	go test $$pkg -v

# 运行特定测试函数
test-func:
	@read -p "请输入包路径: " pkg; \
	read -p "请输入测试函数名: " func; \
	go test $$pkg -run $$func -v

