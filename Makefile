.PHONY: build run clean deps init-db test lint fmt help dev

# 项目配置
APP_NAME := campus-memory
BUILD_DIR := build
MAIN_FILE := main.go

# Go 相关配置
GO := go
GOFLAGS := -v
LDFLAGS := -w -s

# 检测操作系统
ifeq ($(OS),Windows_NT)
    BINARY_EXT := .exe
    RM := del /Q
    MKDIR := mkdir
    RMDIR := rmdir /S /Q
else
    BINARY_EXT :=
    RM := rm -f
    MKDIR := mkdir -p
    RMDIR := rm -rf
endif

# 默认目标
all: deps build

# 显示帮助信息
help:
	@echo "Available commands:"
	@echo "  build     - 编译项目生成可执行文件"
	@echo "  run       - 直接运行项目"
	@echo "  dev       - 开发模式运行（热重载）"
	@echo "  clean     - 清理构建文件和临时文件"
	@echo "  deps      - 安装和更新项目依赖"
	@echo "  init-db   - 初始化数据库和测试数据"
	@echo "  test      - 运行测试"
	@echo "  lint      - 代码检查"
	@echo "  fmt       - 代码格式化"
	@echo "  help      - 显示此帮助信息"

# 安装和更新依赖
deps:
	@echo "正在安装项目依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "依赖安装完成"

# 编译项目
build: deps
	@echo "正在编译项目..."
	$(MKDIR) $(BUILD_DIR) 2>/dev/null || true
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT) $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)$(BINARY_EXT)"

# 运行项目
run: deps
	@echo "正在启动项目..."
	$(GO) run $(MAIN_FILE)

# 开发模式运行
dev: deps
	@echo "正在以开发模式启动项目..."
	@echo "注意: 需要安装 air 工具进行热重载"
	@echo "安装命令: go install github.com/cosmtrek/air@latest"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air 未安装，使用普通模式运行..."; \
		$(GO) run $(MAIN_FILE); \
	fi

# 清理构建文件
clean:
	@echo "正在清理构建文件..."
	$(RMDIR) $(BUILD_DIR) 2>/dev/null || true
	$(RM) *.log 2>/dev/null || true
	$(RM) *.db 2>/dev/null || true
	$(RM) *.sqlite 2>/dev/null || true
	$(RMDIR) data 2>/dev/null || true
	$(RMDIR) logs 2>/dev/null || true
	@echo "清理完成"

# 初始化数据库
init-db:
	@echo "正在初始化数据库..."
	@if [ -f scripts/init_db.sql ]; then \
		echo "执行数据库初始化脚本..."; \
		sqlite3 data/campus_memory.db < scripts/init_db.sql; \
	else \
		echo "数据库初始化脚本不存在，跳过..."; \
	fi
	@if [ -f scripts/init_campus_data.sql ]; then \
		echo "执行测试数据初始化脚本..."; \
		sqlite3 data/campus_memory.db < scripts/init_campus_data.sql; \
	else \
		echo "测试数据脚本不存在，跳过..."; \
	fi
	@echo "数据库初始化完成"

# 运行测试
test: deps
	@echo "正在运行测试..."
	$(GO) test -v ./...
	@echo "测试完成"

# 代码检查
lint:
	@echo "正在进行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，使用 go vet..."; \
		$(GO) vet ./...; \
	fi
	@echo "代码检查完成"

# 代码格式化
fmt:
	@echo "正在格式化代码..."
	$(GO) fmt ./...
	@echo "代码格式化完成"

# 生成 Swagger 文档
swagger:
	@echo "正在生成 API 文档..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g main.go -o docs; \
		echo "API 文档生成完成"; \
	else \
		echo "swag 未安装，请先安装: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 安装开发工具
install-tools:
	@echo "正在安装开发工具..."
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "开发工具安装完成"

# 检查环境
check-env:
	@echo "检查开发环境..."
	@echo "Go 版本: $$($(GO) version)"
	@echo "GOPATH: $$($(GO) env GOPATH)"
	@echo "GOROOT: $$($(GO) env GOROOT)"
	@if [ -f .env ]; then \
		echo "环境变量文件: 存在"; \
	else \
		echo "环境变量文件: 不存在 (请复制 .env.example 到 .env)"; \
	fi

# 完整的项目设置
setup: deps install-tools init-db
	@echo "项目设置完成！"
	@echo "运行 'make run' 启动项目"
	@echo "运行 'make dev' 以开发模式启动项目"