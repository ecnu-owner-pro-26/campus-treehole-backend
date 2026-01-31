.PHONY: build run test clean wire swagger

# 构建项目
build:
	go build -o bin/server main.go

# 运行项目
run:
	go run main.go

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/

# 生成依赖注入代码
wire:
	cd provider && wire

# 生成Swagger文档
swagger:
	swag init -g main.go -o docs

# 安装依赖
deps:
	go mod download
	go mod tidy

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run
