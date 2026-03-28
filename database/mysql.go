package database

import (
	"fmt"
    "LunoFuns-Server/configs"
    "log"
    "time"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)
var DB *gorm.DB

func InitDB() error {
	cfg := configs.GlobalConfig

	dsn := cfg.GetDSN()

	var logLevel logger.LogLevel
	switch cfg.Database.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
        PrepareStmt:            true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)	
	log.Println("Database connection established successfully")
	return nil
}