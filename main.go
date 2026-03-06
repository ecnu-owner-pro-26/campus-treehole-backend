package main

import (
	"campus-memory/api/router"
	"campus-memory/infra"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// TODO: 添加 Swagger 支持时取消注释
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

// 应用程序入口
func main() {

	// 初始化数据库
	db, err := infra.InitDatabase("data/campus_memory.db")
	if err != nil {
		log.Fatalf("fail to init database: %v", err)
	}
	log.Println("Database initialized successfully")

	// 设置路由
	r := router.SetupRoutes(db.DB)

	// TODO: 集成 Swagger API 文档（需要先安装依赖）
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建HTTP服务器
	srv := &http.Server{Addr: ":8080", Handler: r}

	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 添加启动日志
	log.Println("Server started on :8080")
	// log.Println("Swagger docs: http://localhost:8080/swagger/index.html")

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// 关闭数据库
	if err := db.Close(); err != nil {
		log.Printf("Failed to close database: %v", err)
	}
	log.Println("Database closed successfully")

}
