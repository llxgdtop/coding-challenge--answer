package router

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	// 创建 Gin 引擎
	r := gin.Default()

	// 健康检查接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// Todos 相关路由
		todos := api.Group("/todos")
		{
			todos.POST("", controllers.AddTodo)                    // 创建待办事项
			todos.GET("", controllers.GetTodos)                    // 获取待办事项列表（支持筛选和排序）
			todos.GET("/:id", controllers.GetTodoByID)             // 获取单个待办事项
			todos.PUT("/:id", controllers.UpdateTodo)              // 更新待办事项（编辑）
			todos.PUT("/:id/status", controllers.UpdateTodoStatus) // 更新待办事项状态
			todos.DELETE("/:id", controllers.DeleteTodo)           // 删除待办事项
		}
	}

	return r
}
