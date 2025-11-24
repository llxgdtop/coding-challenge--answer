package main

import (
	"backend/config"
	"backend/router"
	"log"
)

func main() {
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 配置路由
	r := router.SetupRouter()

	// 启动服务器
	log.Println("Server starting on :8080")
	log.Println("API available at: http://localhost:8080/api/todos")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
