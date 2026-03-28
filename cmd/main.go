package main

import (
	"flag"
    "log"
    
    "LunoFuns-Server/configs"
    "LunoFuns-Server/database"
    "LunoFuns-Server/internal/models"
    "LunoFuns-Server/internal/routes"
    
    "github.com/gin-gonic/gin"
)

func main() {
	// 命令行参数指定配置文件
    var configPath string
    flag.StringVar(&configPath, "config", "../configs/config.yaml", "config file path")
    flag.Parse()
    
    // 加载配置
    cfg, err := configs.LoadConfig(configPath)
    if err != nil {
        log.Fatal("Failed to load config: ", err)
    }
    
    // 设置Gin模式
    if cfg.IsProduction() {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // 初始化数据库
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to init database: ", err)
    }
    
    // 自动迁移
    if err := database.DB.AutoMigrate(&models.User{}); err != nil {
        log.Fatal("Failed to migrate database: ", err)
    }
    
    // 启动服务器
    r := gin.Default()

    // 注册路由
    routes.SetupRoutes(r, database.DB)

    addr := cfg.GetServerAddr()
    log.Printf("Server starting on %s (env: %s)", addr, cfg.App.Env)
    
    if err := r.Run(addr); err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}