package router

import (
	"github.com/gin-gonic/gin"
	"campus-memory/middleware"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// TODO: 配置CORS中间件
	// r.Use(middleware.CORS())
	// r.Use(middleware.Logger())

	// TODO: 创建handler实例
	// authHandler := handler.NewAuthHandler(authService)
	// memoryHandler := handler.NewMemoryHandler(memoryService)
	// commentHandler := handler.NewCommentHandler(commentService)
	// likeHandler := handler.NewLikeHandler(likeService)
	// locationHandler := handler.NewLocationHandler(locationService)
	// campusHandler := handler.NewCampusHandler(campusService)

	// API路由组
	api := r.Group("/api")

	// ==================== 公开路由（无需认证） ====================

	// 1. 微信登录
	auth := api.Group("/auth")
	{
		// auth.POST("/wechat/login", authHandler.WechatLogin) // 微信登录
	}

	// 2. 快速导航 - 校区和地点列表
	campuses := api.Group("/campuses")
	{
		// campuses.GET("", campusHandler.ListCampuses)                    // 获取校区列表
		// campuses.GET("/:id/locations", campusHandler.GetCampusLocations) // 获取校区地点列表
	}

	// ==================== 需要认证的路由 ====================

	authenticated := api.Group("")
	authenticated.Use(middleware.JWTAuth())
	{
		// 用户信息管理
		authRoutes := authenticated.Group("/auth")
		{
			// authRoutes.GET("/profile", authHandler.GetProfile)       // 获取个人信息
			// authRoutes.PUT("/profile", authHandler.UpdateProfile)    // 更新个人信息
		}

		// 记忆相关路由
		memories := authenticated.Group("/memories")
		{
			// memories.POST("", memoryHandler.CreateMemory)           // 创建记忆
			// memories.GET("/:id", memoryHandler.GetMemory)           // 获取记忆详情
			// memories.GET("", memoryHandler.ListMemories)            // 获取记忆列表
			// memories.PUT("/:id", memoryHandler.UpdateMemory)        // 更新记忆
			// memories.DELETE("/:id", memoryHandler.DeleteMemory)     // 删除记忆

			// 记忆点赞（统一方法）
			// memories.POST("/:id/like", likeHandler.ToggleLike) // 切换记忆点赞状态

			// 记忆评论
			// memories.POST("/:id/comments", commentHandler.CreateComment)   // 创建评论
			// memories.GET("/:id/comments", commentHandler.ListComments)     // 获取评论列表
		}

		// 评论相关路由
		comments := authenticated.Group("/comments")
		{
			// comments.DELETE("/:id", commentHandler.DeleteComment)     // 删除评论

			// 评论点赞（统一方法）
			// comments.POST("/:id/like", likeHandler.ToggleLike) // 切换评论点赞状态
		}

		// 地点相关路由
		locations := authenticated.Group("/locations")
		{
			// locations.GET("/:id/memories", locationHandler.GetLocationMemories) // 获取地点记忆列表
		}
	}
}
