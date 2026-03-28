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
    
    // 初始化控制器
    userController := controllers.NewUserController(userService)
    
    // 公开路由
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", userController.Register)
        auth.POST("/login", userController.Login)
    }
    
    // 需要认证的路由
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", userController.GetProfile)
        protected.PUT("/profile", userController.UpdateProfile)
        protected.POST("/change-password", userController.ChangePassword)
    }
    
    // 管理员路由
    admin := r.Group("/api/admin")
    admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        admin.GET("/users", userController.GetUserList)
        admin.PUT("/users/:id/status", userController.UpdateUserStatus)
    }
}