package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"LunoFuns-Server/internal/controllers"
	"LunoFuns-Server/internal/middleware"
	"LunoFuns-Server/internal/services"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化服务
	userService := services.NewUserService(db)
	uploadService := services.NewUploadService()
	videoService := services.NewVideoService(db)

	// 初始化控制器
	userController := controllers.NewUserController(userService)
	uploadController := controllers.NewUploadController(uploadService)
	videoController := controllers.NewVideoController(videoService)

	// 公开路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// 视频公开路由
	videos := r.Group("/api/videos")
	{
		videos.GET("", videoController.GetVideoList)
		videos.GET("/:id", videoController.GetVideoDetail)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/auth/logout", userController.Logout)
		protected.GET("/profile", userController.GetProfile)
		protected.PUT("/profile", userController.UpdateProfile)
		protected.POST("/change-password", userController.ChangePassword)

		// 视频发布路由
		protected.POST("/videos", videoController.UploadVideo)

		// 上传相关
		upload := protected.Group("/upload")
		{
			upload.POST("/token", uploadController.GetUploadToken)                       // 封面等简单上传
			upload.POST("/multipart/init", uploadController.InitMultipartUpload)         // 初始化分片上传
			upload.POST("/multipart/complete", uploadController.CompleteMultipartUpload) // 完成分片上传合并
		}
	}

	// 管理员路由
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/users", userController.GetUserList)
		admin.PUT("/users/:id/status", userController.UpdateUserStatus)
	}
}
