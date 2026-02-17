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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 应用程序入口
func main() {

	// 初始化数据库
	db, err := infra.InitDatabase("data/campus_memory.db")
	if err != nil {
		log.Fatalf("fail to init database: %v", err)
	}
	log.Println(" Database initialized successfully")

	// 设置路由
	r := router.SetupRoutes(db)

	// 集成 Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建HTTP服务器
	srv := &http.Server{Addr: ":8080", Handler: r}

	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

	// 关闭数据库
	db.Close()

}
