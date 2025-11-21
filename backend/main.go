package main

import (
	"backend/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 健康检查接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
