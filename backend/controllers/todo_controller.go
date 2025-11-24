package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var todoService = services.NewTodoService()

// AddTodo 添加待办事项
// POST /api/todos
func AddTodo(c *gin.Context) {
	var input models.CreateTodoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid input: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用 Service 层创建
	todo, err := todoService.CreateTodo(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to create todo: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    todo,
	})
}

// GetTodos 获取待办事项列表
// GET /api/todos?category=work&sort=priority
func GetTodos(c *gin.Context) {
	// 获取查询参数
	category := c.DefaultQuery("category", "")
	sortBy := c.DefaultQuery("sort", "")

	// 调用 Service 层获取列表
	todos, err := todoService.GetAllTodos(category, sortBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    todos,
	})
}

// GetTodoByID 根据 ID 获取待办事项
// GET /api/todos/:id
func GetTodoByID(c *gin.Context) {
	// 获取 ID 参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid ID format",
			"data":    nil,
		})
		return
	}

	// 调用 Service 层查询
	todo, err := todoService.GetTodoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Todo not found: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    todo,
	})
}

// UpdateTodo 更新待办事项（编辑）
// PUT /api/todos/:id
func UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid ID format",
			"data":    nil,
		})
		return
	}

	var input models.UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid input: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用 Service 层更新
	todo, err := todoService.UpdateTodo(uint(id), &input)
	if err != nil {
		// 判断是否是版本冲突错误
		if conflictErr, ok := err.(*services.VersionConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"code":             http.StatusConflict,
				"message":          conflictErr.Message,
				"current_version":  conflictErr.CurrentVersion,
				"provided_version": conflictErr.ProvidedVersion,
				"latest_data":      conflictErr.LatestData,
			})
			return
		}

		// 其他错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update todo: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    todo,
	})
}

// UpdateTodoStatus 更新待办事项状态（完成/未完成）
// PUT /api/todos/:id/status
func UpdateTodoStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid ID format",
			"data":    nil,
		})
		return
	}

	var input models.UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid input: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用 Service 层更新状态
	todo, err := todoService.UpdateTodoStatus(uint(id), &input)
	if err != nil {
		// 判断是否是版本冲突错误
		if conflictErr, ok := err.(*services.VersionConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"code":             http.StatusConflict,
				"message":          conflictErr.Message,
				"current_version":  conflictErr.CurrentVersion,
				"provided_version": conflictErr.ProvidedVersion,
				"latest_data":      conflictErr.LatestData,
			})
			return
		}

		// 其他错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update status: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    todo,
	})
}

// DeleteTodo 删除待办事项
// DELETE /api/todos/:id
func DeleteTodo(c *gin.Context) {
	// 获取 ID 参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid ID format",
			"data":    nil,
		})
		return
	}

	// 调用 Service 层删除
	err = todoService.DeleteTodo(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Failed to delete todo: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    nil,
	})
}
