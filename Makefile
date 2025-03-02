.DEFAULT_GOAL := gogen

# ANSI 颜色定义
COLOR_RESET := \033[0m
COLOR_CYAN := \033[36m
COLOR_GREEN := \033[32m
COLOR_RED := \033[31m

ECHO = @echo -e

# 获取 Git 提交哈希和时间，默认值处理非 Git 环境
GITCOMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GITDATE := $(shell git show -s --format='%ct' 2>/dev/null || echo "1970-01-01 00:00:00 +0000")

# 构造链接器标志
LDFLAGS := -ldflags "-X main.GitCommit=$(GITCOMMIT) -X main.GitDate=$(GITDATE)"
PROJECT_NAME := $(shell go list -m | awk -F/ '{print $$NF}' || echo "gogen")

# 整理 Go 模块依赖
tidy:
	$(ECHO) "$(COLOR_CYAN)Tidying Go modules...$(COLOR_RESET)"
	go mod tidy
	$(ECHO) "$(COLOR_GREEN)Tidy completed$(COLOR_RESET)"

# 编译 gogen 程序，嵌入 Git 提交信息
gogen: tidy
	go build -v $(LDFLAGS) -o $(PROJECT_NAME) ./main.go

# 清理生成的文件和 Go 缓存
clean:
	rm -f gogen
	go clean -cache -testcache

# 运行所有测试
test: tidy
	go test -v ./...

# 检查代码风格和潜在问题
lint: tidy
	golangci-lint run ./...

# 编译协议文件
proto:
	@test -f ./bin/compile.sh || (echo "compile.sh not found" && exit 1)
	sh ./bin/compile.sh

# 显示帮助信息
help:
	$(ECHO) "$(COLOR_CYAN)Available targets:$(COLOR_RESET)"
	$(ECHO) "  gogen    : Build the gogen binary with Git info"
	$(ECHO) "  clean    : Remove generated files and caches"
	$(ECHO) "  test     : Run all tests with coverage"
	$(ECHO) "  lint     : Check code style and issues"
	$(ECHO) "  proto    : Compile protocol files"
	$(ECHO) "  tidy     : Tidy Go module dependencies"
	$(ECHO) "  help     : Show this help message"

.PHONY: gogen clean test lint proto tidy
