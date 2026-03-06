#!/bin/bash

# ===========================================
# 校园记忆树洞后端 - 开发环境一键设置脚本
# ===========================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if command -v "$1" >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# 检查Go版本
check_go_version() {
    local required_version="1.21"
    local current_version
    
    if ! check_command go; then
        log_error "Go 未安装，请先安装 Go $required_version 或更高版本"
        log_info "下载地址: https://golang.org/dl/"
        exit 1
    fi
    
    current_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    log_info "检测到 Go 版本: $current_version"
    
    # 简单的版本比较
    if [[ "$(printf '%s\n' "$required_version" "$current_version" | sort -V | head -n1)" != "$required_version" ]]; then
        log_error "Go 版本过低，需要 $required_version 或更高版本"
        exit 1
    fi
    
    log_success "Go 版本检查通过"
}

# 检查Git
check_git() {
    if ! check_command git; then
        log_error "Git 未安装，请先安装 Git"
        exit 1
    fi
    log_success "Git 检查通过"
}

# 检查SQLite
check_sqlite() {
    if ! check_command sqlite3; then
        log_warning "SQLite3 未安装，将尝试自动安装..."
        
        # 根据操作系统安装SQLite
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            if check_command apt-get; then
                sudo apt-get update && sudo apt-get install -y sqlite3
            elif check_command yum; then
                sudo yum install -y sqlite
            elif check_command dnf; then
                sudo dnf install -y sqlite
            else
                log_error "无法自动安装 SQLite3，请手动安装"
                exit 1
            fi
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            if check_command brew; then
                brew install sqlite
            else
                log_error "请先安装 Homebrew，然后运行: brew install sqlite"
                exit 1
            fi
        else
            log_error "请手动安装 SQLite3"
            exit 1
        fi
    fi
    log_success "SQLite3 检查通过"
}

# 创建必要的目录
create_directories() {
    log_info "创建项目目录结构..."
    
    local dirs=(
        "data"
        "logs" 
        "storage/uploads"
        "build"
        "scripts/output"
        "scripts/logs"
    )
    
    for dir in "${dirs[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            log_info "创建目录: $dir"
        fi
    done
    
    # 创建 .gitkeep 文件保持空目录
    touch storage/.gitkeep
    touch logs/.gitkeep
    
    log_success "目录结构创建完成"
}

# 复制环境变量配置文件
setup_env_file() {
    log_info "设置环境变量配置文件..."
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            log_success "已复制 .env.example 到 .env"
            log_warning "请根据实际情况修改 .env 文件中的配置"
        else
            log_error ".env.example 文件不存在"
            exit 1
        fi
    else
        log_info ".env 文件已存在，跳过复制"
    fi
}

# 安装Go依赖
install_dependencies() {
    log_info "安装项目依赖..."
    
    # 设置Go模块代理（中国用户）
    if [[ "${GOPROXY:-}" == "" ]]; then
        export GOPROXY=https://goproxy.cn,direct
        log_info "设置 GOPROXY 为 https://goproxy.cn,direct"
    fi
    
    # 下载依赖
    go mod download
    
    # 整理依赖
    go mod tidy
    
    # 验证依赖
    go mod verify
    
    log_success "依赖安装完成"
}

# 安装开发工具
install_dev_tools() {
    log_info "安装开发工具..."
    
    local tools=(
        "github.com/cosmtrek/air@latest"
        "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        "github.com/swaggo/swag/cmd/swag@latest"
    )
    
    for tool in "${tools[@]}"; do
        log_info "安装 $tool"
        if go install "$tool"; then
            log_success "安装成功: $tool"
        else
            log_warning "安装失败: $tool (可选工具，不影响项目运行)"
        fi
    done
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    
    # 确保数据目录存在
    mkdir -p data
    
    # 检查数据库初始化脚本
    if [ -f "scripts/init_db.sql" ]; then
        log_info "执行数据库初始化脚本..."
        sqlite3 data/campus_memory.db < scripts/init_db.sql
        log_success "数据库结构初始化完成"
    else
        log_warning "数据库初始化脚本不存在: scripts/init_db.sql"
        # 创建一个空的数据库文件
        touch data/campus_memory.db
        log_info "创建空数据库文件"
    fi
    
    # 检查测试数据脚本
    if [ -f "scripts/init_campus_data.sql" ]; then
        log_info "执行测试数据初始化脚本..."
        sqlite3 data/campus_memory.db < scripts/init_campus_data.sql
        log_success "测试数据初始化完成"
    else
        log_warning "测试数据脚本不存在: scripts/init_campus_data.sql"
    fi
}

# 验证项目设置
verify_setup() {
    log_info "验证项目设置..."
    
    # 检查Go模块
    if go list -m all >/dev/null 2>&1; then
        log_success "Go 模块验证通过"
    else
        log_error "Go 模块验证失败"
        exit 1
    fi
    
    # 检查编译
    log_info "测试项目编译..."
    if go build -o build/test_build main.go 2>/dev/null; then
        log_success "项目编译测试通过"
        rm -f build/test_build
    else
        log_error "项目编译失败，请检查代码"
        exit 1
    fi
    
    # 检查数据库文件
    if [ -f "data/campus_memory.db" ]; then
        log_success "数据库文件存在"
    else
        log_error "数据库文件不存在"
        exit 1
    fi
    
    # 检查环境变量文件
    if [ -f ".env" ]; then
        log_success "环境变量文件存在"
    else
        log_error "环境变量文件不存在"
        exit 1
    fi
}

# 显示后续步骤
show_next_steps() {
    log_success "开发环境设置完成！"
    echo
    echo "===========================================
    echo "后续步骤:"
    echo "1. 检查并修改 .env 文件中的配置"
    echo "2. 运行项目:"
    echo "   make run          # 普通模式运行"
    echo "   make dev          # 开发模式运行（热重载）"
    echo "3. 测试API:"
    echo "   使用 scripts/test_api.http 文件测试接口"
    echo "4. 查看更多命令:"
    echo "   make help         # 显示所有可用命令"
    echo "==========================================="
    echo
    log_info "项目地址: http://localhost:8080"
    log_info "API文档: http://localhost:8080/swagger/index.html (如果启用)"
}

# 主函数
main() {
    echo "==========================================="
    echo "校园记忆树洞后端 - 开发环境设置"
    echo "==========================================="
    echo
    
    # 检查运行环境
    log_info "检查运行环境..."
    check_go_version
    check_git
    check_sqlite
    
    # 设置项目
    log_info "设置项目环境..."
    create_directories
    setup_env_file
    install_dependencies
    
    # 安装开发工具（可选）
    read -p "是否安装开发工具 (air, golangci-lint, swag)? [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        install_dev_tools
    else
        log_info "跳过开发工具安装"
    fi
    
    # 初始化数据库
    init_database
    
    # 验证设置
    verify_setup
    
    # 显示后续步骤
    show_next_steps
}

# 错误处理
trap 'log_error "脚本执行失败，请检查错误信息"; exit 1' ERR

# 运行主函数
main "$@"