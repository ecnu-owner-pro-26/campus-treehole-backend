@echo off
REM ===========================================
REM 校园记忆树洞后端 - Windows 开发环境设置脚本
REM ===========================================

setlocal enabledelayedexpansion

echo ===========================================
echo 校园记忆树洞后端 - 开发环境设置
echo ===========================================
echo.

REM 检查 Go 是否安装
echo [INFO] 检查 Go 环境...
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Go 未安装，请先安装 Go 1.21 或更高版本
    echo [INFO] 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

REM 显示 Go 版本
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo [SUCCESS] 检测到 Go 版本: %GO_VERSION%

REM 检查 Git 是否安装
echo [INFO] 检查 Git 环境...
where git >nul 2>nul
if %errorlevel% neq 0 (
    echo [WARNING] Git 未安装，建议安装 Git
    echo [INFO] 下载地址: https://git-scm.com/download/win
) else (
    echo [SUCCESS] Git 检查通过
)

REM 创建必要的目录
echo [INFO] 创建项目目录结构...
if not exist "data" mkdir data
if not exist "logs" mkdir logs
if not exist "storage\uploads" mkdir storage\uploads
if not exist "build" mkdir build
if not exist "scripts\output" mkdir scripts\output
if not exist "scripts\logs" mkdir scripts\logs

REM 创建 .gitkeep 文件
type nul > storage\.gitkeep
type nul > logs\.gitkeep

echo [SUCCESS] 目录结构创建完成

REM 复制环境变量配置文件
echo [INFO] 设置环境变量配置文件...
if not exist ".env" (
    if exist ".env.example" (
        copy .env.example .env >nul
        echo [SUCCESS] 已复制 .env.example 到 .env
        echo [WARNING] 请根据实际情况修改 .env 文件中的配置
    ) else (
        echo [ERROR] .env.example 文件不存在
        pause
        exit /b 1
    )
) else (
    echo [INFO] .env 文件已存在，跳过复制
)

REM 设置 Go 代理（中国用户）
echo [INFO] 设置 Go 模块代理...
set GOPROXY=https://goproxy.cn,direct
echo [INFO] GOPROXY 设置为: %GOPROXY%

REM 安装项目依赖
echo [INFO] 安装项目依赖...
go mod download
if %errorlevel% neq 0 (
    echo [ERROR] 依赖下载失败
    pause
    exit /b 1
)

go mod tidy
if %errorlevel% neq 0 (
    echo [ERROR] 依赖整理失败
    pause
    exit /b 1
)

go mod verify
if %errorlevel% neq 0 (
    echo [ERROR] 依赖验证失败
    pause
    exit /b 1
)

echo [SUCCESS] 依赖安装完成

REM 询问是否安装开发工具
echo.
set /p INSTALL_TOOLS="是否安装开发工具 (air, golangci-lint, swag)? [y/N]: "
if /i "%INSTALL_TOOLS%"=="y" (
    echo [INFO] 安装开发工具...
    
    echo [INFO] 安装 air (热重载工具)...
    go install github.com/cosmtrek/air@latest
    
    echo [INFO] 安装 golangci-lint (代码检查工具)...
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    echo [INFO] 安装 swag (API文档生成工具)...
    go install github.com/swaggo/swag/cmd/swag@latest
    
    echo [SUCCESS] 开发工具安装完成
) else (
    echo [INFO] 跳过开发工具安装
)

REM 初始化数据库
echo [INFO] 初始化数据库...
if exist "scripts\init_db.sql" (
    echo [INFO] 执行数据库初始化脚本...
    REM Windows 下可能需要手动执行 SQL 脚本
    echo [WARNING] 请手动执行: sqlite3 data\campus_memory.db ^< scripts\init_db.sql
) else (
    echo [WARNING] 数据库初始化脚本不存在: scripts\init_db.sql
    REM 创建空数据库文件
    type nul > data\campus_memory.db
    echo [INFO] 创建空数据库文件
)

if exist "scripts\init_campus_data.sql" (
    echo [INFO] 执行测试数据初始化脚本...
    echo [WARNING] 请手动执行: sqlite3 data\campus_memory.db ^< scripts\init_campus_data.sql
) else (
    echo [WARNING] 测试数据脚本不存在: scripts\init_campus_data.sql
)

REM 验证项目设置
echo [INFO] 验证项目设置...

REM 检查 Go 模块
go list -m all >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Go 模块验证失败
    pause
    exit /b 1
)
echo [SUCCESS] Go 模块验证通过

REM 测试编译
echo [INFO] 测试项目编译...
go build -o build\test_build.exe main.go 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] 项目编译失败，请检查代码
    pause
    exit /b 1
)
echo [SUCCESS] 项目编译测试通过
del build\test_build.exe 2>nul

REM 检查数据库文件
if exist "data\campus_memory.db" (
    echo [SUCCESS] 数据库文件存在
) else (
    echo [ERROR] 数据库文件不存在
)

REM 检查环境变量文件
if exist ".env" (
    echo [SUCCESS] 环境变量文件存在
) else (
    echo [ERROR] 环境变量文件不存在
)

REM 显示后续步骤
echo.
echo [SUCCESS] 开发环境设置完成！
echo.
echo ===========================================
echo 后续步骤:
echo 1. 检查并修改 .env 文件中的配置
echo 2. 运行项目:
echo    go run main.go        # 普通模式运行
echo    air                    # 开发模式运行（需要安装 air）
echo 3. 测试API:
echo    使用 scripts\test_api.http 文件测试接口
echo 4. 查看更多命令:
echo    查看 Makefile 文件了解所有可用命令
echo ===========================================
echo.
echo [INFO] 项目地址: http://localhost:8080
echo [INFO] API文档: http://localhost:8080/swagger/index.html (如果启用)
echo.

pause