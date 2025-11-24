package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var todoService = services.NewTodoService()

// AddTodo 添加待办事项
// POST /api/todos
func AddTodo(c *gin.Context) {
	var input models.CreateTodoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	// 调用 Service 层创建
	todo, err := todoService.CreateTodo(&input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, todo)
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
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, todos)
}

// GetTodoByID 根据 ID 获取待办事项
// GET /api/todos/:id
func GetTodoByID(c *gin.Context) {
	// 获取 ID 参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID format")
		return
	}

	// 调用 Service 层查询
	todo, err := todoService.GetTodoByID(uint(id))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, todo)
}

// UpdateTodo 更新待办事项（编辑）
// PUT /api/todos/:id
func UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID format")
		return
	}

	var input models.UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	// 调用 Service 层更新
	todo, err := todoService.UpdateTodo(uint(id), &input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, todo)
}

// UpdateTodoStatus 更新待办事项状态（完成/未完成）
// PUT /api/todos/:id/status
func UpdateTodoStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID format")
		return
	}

	var input models.UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	// 调用 Service 层更新状态
	todo, err := todoService.UpdateTodoStatus(uint(id), &input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, todo)
}

// DeleteTodo 删除待办事项
// DELETE /api/todos/:id
func DeleteTodo(c *gin.Context) {
	// 获取 ID 参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID format")
		return
	}

	// 调用 Service 层删除
	err = todoService.DeleteTodo(uint(id))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "Todo deleted successfully", nil)
}
