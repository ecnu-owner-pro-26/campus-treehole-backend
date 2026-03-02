package main

import (
	"campus-memory/api/handler"
	"campus-memory/api/router"
	"campus-memory/application/assembler"
	"campus-memory/application/service"
	"campus-memory/infra"
	"campus-memory/infra/repo"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

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

	// 初始化 Repository 层
	userRepo := repo.NewUserRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	commentRepo := repo.NewCommentRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)

	// 初始化 Service 层
	authService := service.NewAuthService(userRepo)
	memoryService := service.NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)
	locationService := service.NewLocationService(locationRepo, campusRepo)
	campusService := service.NewCampusService(campusRepo, locationRepo)
	commentService := service.NewCommentService(*commentRepo, *userRepo, *likeRepo, assembler.NewCommentAssembler(), *memoryRepo)
	likeService := service.NewLikeService(likeRepo, memoryRepo, commentRepo)
	imageService := service.NewImageService(imageRepo, memoryRepo)
	quicknavService := service.NewQuicknavService(db)

	// 初始化 Handler 层
	authHandler := handler.NewAuthHandler(authService)
	memoryHandler := handler.NewMemoryHandler(memoryService)
	locationHandler := handler.NewLocationHandler(locationService)
	campusHandler := handler.NewCampusHandler(campusService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(likeService)
	imageHandler := handler.NewImageHandler(imageService)
	quicknavHandler := handler.NewQuicknavHandler(quicknavService)

	// 创建 Gin 引擎（不使用默认中间件，使用自定义中间件）
	r := gin.New()
	
	// 添加崩溃恢复中间件
	r.Use(gin.Recovery())

	// 设置路由
	router.SetupRoutes(r, authHandler, memoryHandler, commentHandler, likeHandler, locationHandler, campusHandler, imageHandler, quicknavHandler)

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
