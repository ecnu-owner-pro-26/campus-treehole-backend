package router

// SetupRoutes 设置所有路由
func SetupRoutes() {
	// TODO: 初始化Gin路由器
	// TODO: 配置CORS中间件
	
	// TODO: 创建handler实例
	// authHandler := handler.NewAuthHandler(db)
	// memoryHandler := handler.NewMemoryHandler(db)
	// commentHandler := handler.NewCommentHandler(db)
	// locationHandler := handler.NewLocationHandler(db)
	// campusHandler := handler.NewCampusHandler(db)

	// API路由组
	// api := router.Group("/api")
	
	// 1. 快速导航路由（无需认证）- 支持选择校区和地点
	// campuses := api.Group("/campuses")
	// {
	//     campuses.GET("", campusHandler.ListCampuses)                    // 获取校区列表
	//     campuses.GET("/:id/locations", campusHandler.GetCampusLocations) // 获取校区地点列表
	// }

	// 2. 认证相关路由（无需认证）
	// auth := api.Group("/auth")
	// {
	//     auth.POST("/register", authHandler.Register) // 用户注册
	//     auth.POST("/login", authHandler.Login)       // 用户登录
	// }

	// 3. 需要认证的路由
	// authenticated := api.Group("")
	// authenticated.Use(middleware.JWTAuth())
	// {
	//     // 用户信息
	//     authenticated.GET("/auth/profile", authHandler.GetProfile) // 获取个人信息

	//     // 记忆相关路由 - 支持发表和浏览记忆
	//     memories := authenticated.Group("/memories")
	//     {
	//         memories.POST("", memoryHandler.CreateMemory)           // 创建记忆
	//         memories.GET("/:id", memoryHandler.GetMemory)           // 获取记忆详情
	//         memories.GET("", memoryHandler.ListMemories)            // 获取记忆列表
	//         memories.PUT("/:id", memoryHandler.UpdateMemory)        // 更新记忆
	//         memories.DELETE("/:id", memoryHandler.DeleteMemory)     // 删除记忆

	//         // 记忆点赞 - 支持给记忆点赞
	//         memories.POST("/:id/like", memoryHandler.LikeMemory)     // 点赞记忆
	//         memories.DELETE("/:id/like", memoryHandler.UnlikeMemory) // 取消点赞记忆

	//         // 记忆的评论 - 支持给记忆评论
	//         memories.POST("/:id/comments", commentHandler.CreateComment)   // 创建评论
	//         memories.GET("/:id/comments", commentHandler.ListComments)     // 获取评论列表
	//     }

	//     // 评论相关路由 - 支持评论管理和点赞
	//     comments := authenticated.Group("/comments")
	//     {
	//         comments.DELETE("/:id", commentHandler.DeleteComment)     // 删除评论
	//         comments.POST("/:id/like", commentHandler.LikeComment)     // 点赞评论
	//         comments.DELETE("/:id/like", commentHandler.UnlikeComment) // 取消点赞评论
	//     }

	//     // 地点相关路由 - 支持浏览地点记忆
	//     locations := authenticated.Group("/locations")
	//     {
	//         locations.GET("/:id/memories", locationHandler.GetLocationMemories) // 获取地点记忆列表
	//     }
	// }

	// TODO: 返回配置好的路由器
}
