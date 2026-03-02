package router

import (
	"campus-memory/api/handler"
	"campus-memory/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine,
	authHandler *handler.AuthHandler,
	memoryHandler *handler.MemoryHandler,
	commentHandler *handler.CommentHandler,
	likeHandler *handler.LikeHandler,
	locationHandler *handler.LocationHandler,
	campusHandler *handler.CampusHandler,
	imageHandler *handler.ImageHandler,
	quicknavHandler *handler.QuicknavHandler,
) {
	// 配置CORS中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// API路由组
	api := r.Group("/api")

	// ==================== 公开路由（无需认证） ====================

	// 1. 微信登录
	auth := api.Group("/auth")
	{
		auth.POST("/wechat/login", authHandler.WechatLogin) // 微信登录
	}

	// 2. 快速导航 - 校区和地点列表（公开访问）
	campuses := api.Group("/campuses")
	{
		campuses.GET("", campusHandler.ListCampuses)                       // 获取所有校区列表
		campuses.GET("/:id", campusHandler.GetCampus)                      // 获取校区详情
		campuses.GET("/:id/locations", campusHandler.GetCampusWithLocations) // 获取校区及其地点列表
	}

	// 3. 快速导航树（公开访问）
	api.GET("/quicknav", quicknavHandler.GetNavTree) // 获取完整导航树

	// ==================== 需要认证的路由 ====================

	authenticated := api.Group("")
	authenticated.Use(middleware.JWTAuth())
	{
		// 用户信息管理
		authRoutes := authenticated.Group("/auth")
		{
			authRoutes.GET("/profile", authHandler.GetProfile)    // 获取个人信息
			authRoutes.PUT("/profile", authHandler.UpdateProfile) // 更新个人信息
		}

		// 记忆相关路由
		memories := authenticated.Group("/memories")
		{
			memories.POST("", memoryHandler.CreateMemory)       // 创建记忆
			memories.GET("/:id", memoryHandler.GetMemory)       // 获取记忆详情
			memories.GET("", memoryHandler.ListMemories)        // 获取记忆列表
			memories.PUT("/:id", memoryHandler.UpdateMemory)    // 更新记忆
			memories.DELETE("/:id", memoryHandler.DeleteMemory) // 删除记忆

			// 记忆点赞（统一方法）
			memories.POST("/:id/like", likeHandler.ToggleLike) // 切换记忆点赞状态

			// 记忆的图片
			memories.GET("/:memory_id/images", imageHandler.GetImagesByMemoryID) // 获取记忆的所有图片
		}

		// 评论相关路由
		comments := authenticated.Group("/comments")
		{
			comments.POST("", commentHandler.CreateComment)           // 创建评论
			comments.GET("", commentHandler.ListComments)             // 获取评论列表
			comments.GET("/:parent_id/replies", commentHandler.ListReplies) // 获取回复列表
			comments.DELETE("/:id", commentHandler.DeleteComment)     // 删除评论

			// 评论点赞（统一方法）
			comments.POST("/:id/like", likeHandler.ToggleLike) // 切换评论点赞状态
		}

		// 地点相关路由
		locations := authenticated.Group("/locations")
		{
			locations.GET("/:id", locationHandler.GetLocation)       // 获取地点详情
			locations.GET("", locationHandler.ListLocations)         // 获取地点列表
			locations.GET("/search", locationHandler.SearchLocations) // 搜索地点
			locations.POST("", locationHandler.CreateLocation)       // 创建地点（管理员）
			locations.PUT("/:id", locationHandler.UpdateLocation)    // 更新地点（管理员）
			locations.DELETE("/:id", locationHandler.DeleteLocation) // 删除地点（管理员）
		}

		// 图片相关路由
		images := authenticated.Group("/images")
		{
			images.POST("/upload", imageHandler.UploadImage)  // 上传图片
			images.DELETE("/:id", imageHandler.DeleteImage)   // 删除图片
		}
	}
}
