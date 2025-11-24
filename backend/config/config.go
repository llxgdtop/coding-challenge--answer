package config

import (
	customerrors "backend/errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDefaultConfig 获取默认数据库配置
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     "4000",
		User:     "root",
		Password: "",
		DBName:   "todo_app",
	}
}

// InitDB 初始化数据库连接
func InitDB() error {
	config := GetDefaultConfig()

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启 SQL 日志
	})

	if err != nil {
		return fmt.Errorf("%w: %v", customerrors.ErrDatabaseConnection, err)
	}

	// 获取底层的 sql.DB 对象，用于配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("%w: %v", customerrors.ErrDatabaseInit, err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)      // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)     // 最大打开连接数
	sqlDB.SetConnMaxLifetime(3600) // 连接最大生命周期（秒）

	log.Println("Database connected successfully!")
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
