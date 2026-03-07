package router

import (
	"campus-memory/api/handler"
	"campus-memory/application/assembler"
	"campus-memory/application/service"
	"campus-memory/infra/repo"
	"campus-memory/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置所有路由
func SetupRoutes(db *gorm.DB) *gin.Engine {
	// 创建 Gin 引擎
	r := gin.Default()

	// 配置静态文件服务（用于图片访问）
	r.Static("/uploads", "./uploads")

	// TODO: 配置CORS中间件
	// r.Use(middleware.CORS())
	// r.Use(middleware.Logger())

	// ==================== 初始化仓储层 ====================
	userRepo := repo.NewUserRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	commentRepo := repo.NewCommentRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)

	// ==================== 初始化组装器 ====================
	commentAssembler := assembler.NewCommentAssembler()

	// ==================== 初始化服务层 ====================
	authService := service.NewAuthService(userRepo)
	campusService := service.NewCampusService(campusRepo, locationRepo)
	locationService := service.NewLocationService(locationRepo, campusRepo)
	memoryService := service.NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)
	commentService := service.NewCommentService(commentRepo, userRepo, likeRepo, commentAssembler, memoryRepo)
	likeService := service.NewLikeService(likeRepo, memoryRepo, commentRepo)
	imageService := service.NewImageService(imageRepo, memoryRepo)
	quicknavService := service.NewQuickNavService(db)

	// ==================== 初始化处理器层 ====================
	authHandler := handler.NewAuthHandler(authService)
	campusHandler := handler.NewCampusHandler(campusService)
	locationHandler := handler.NewLocationHandler(locationService)
	memoryHandler := handler.NewMemoryHandler(memoryService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(likeService)
	imageHandler := handler.NewImageHandler(imageService)
	quicknavHandler := handler.NewQuickNavHandler(*quicknavService)

	// API路由组
	api := r.Group("/api")

	// ==================== 公开路由（无需认证） ====================

	// 1. 微信登录
	auth := api.Group("/auth")
	{
		auth.POST("/wechat/login", authHandler.WechatLogin)   // 微信登录（标准路由）
		auth.POST("/wechat-login", authHandler.WechatLogin)   // 微信登录（兼容连字符）
		auth.POST("/login", authHandler.WechatLogin)          // 微信登录（兼容简化路由）
	}

	// 2. 校区和地点列表
	campuses := api.Group("/campuses")
	{
		campuses.GET("", campusHandler.ListCampuses)                       // 获取校区列表
		campuses.GET("/:id", campusHandler.GetCampus)                      // 获取校区详情
		campuses.GET("/:id/locations", campusHandler.GetCampusWithLocations) // 获取校区地点列表
	}

	// 3. 快速导航（公开访问）
	quicknav := api.Group("/quicknav")
	{
		quicknav.GET("/tree", quicknavHandler.GetNavTree)                       // 获取导航树
		quicknav.GET("/search", quicknavHandler.SearchLocations)                // 搜索地点
		quicknav.GET("/popular", quicknavHandler.GetPopularLocations)           // 获取热门地点
		quicknav.GET("/category", quicknavHandler.GetLocationsByCategory)       // 根据类别获取地点
	}

	// 4. 地点列表（公开访问）
	locations := api.Group("/locations")
	{
		locations.GET("", locationHandler.ListLocations)       // 获取地点列表
		locations.GET("/:id", locationHandler.GetLocation)     // 获取地点详情
		locations.GET("/search", locationHandler.SearchLocations) // 搜索地点
	}

	// 5. 记忆列表（公开访问，使用可选认证）
	r.GET("/api/memories", middleware.OptionalAuth(), memoryHandler.ListMemories)       // 获取记忆列表
	r.GET("/api/memories/:id", middleware.OptionalAuth(), memoryHandler.GetMemory)      // 获取记忆详情

	// 6. 评论列表（公开访问，使用可选认证）
	r.GET("/api/comments", middleware.OptionalAuth(), commentHandler.ListComments)      // 获取评论列表
	r.GET("/api/comments/:parentId/replies", middleware.OptionalAuth(), commentHandler.ListReplies) // 获取回复列表

	// ==================== 需要认证的路由 ====================

	authenticated := api.Group("")
	authenticated.Use(middleware.JWTAuth())
	{
		// 图片上传（兼容简化路由）
		authenticated.POST("/upload", imageHandler.UploadImage)
		// 用户信息管理
		authRoutes := authenticated.Group("/auth")
		{
			authRoutes.GET("/profile", authHandler.GetProfile)       // 获取个人信息
			authRoutes.PUT("/profile", authHandler.UpdateProfile)    // 更新个人信息
		}

		// 记忆相关路由
		memories := authenticated.Group("/memories")
		{
			memories.POST("", memoryHandler.CreateMemory)           // 创建记忆
			memories.PUT("/:id", memoryHandler.UpdateMemory)        // 更新记忆
			memories.DELETE("/:id", memoryHandler.DeleteMemory)     // 删除记忆

			// 记忆点赞
			memories.POST("/:id/like", likeHandler.ToggleLike)      // 切换记忆点赞状态
		}

		// 评论相关路由
		comments := authenticated.Group("/comments")
		{
			comments.POST("", commentHandler.CreateComment)         // 创建评论
			comments.DELETE("/:id", commentHandler.DeleteComment)   // 删除评论

			// 评论点赞
			comments.POST("/:id/like", likeHandler.ToggleLike)      // 切换评论点赞状态
		}

		// 地点管理（需要认证）
		locationsAuth := authenticated.Group("/locations")
		{
			locationsAuth.POST("", locationHandler.CreateLocation)       // 创建地点
			locationsAuth.PUT("/:id", locationHandler.UpdateLocation)    // 更新地点
			locationsAuth.DELETE("/:id", locationHandler.DeleteLocation) // 删除地点
		}

		// 图片上传
		images := authenticated.Group("/images")
		{
			images.POST("/upload", imageHandler.UploadImage)                    // 上传图片
			images.DELETE("/:id", imageHandler.DeleteImage)                     // 删除图片
			images.GET("/memory/:memory_id", imageHandler.GetImagesByMemoryID)  // 获取记忆的所有图片
		}
	}

	return r
}
