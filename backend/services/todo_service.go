package services

import (
	customerrors "backend/errors"
	"backend/models"
	"errors"
	"fmt"
	"strings"
)

// TodoService 待办事项业务逻辑服务，一切数据库查询放到models/todo.go中
type TodoService struct{}

func NewTodoService() *TodoService {
	return &TodoService{}
}

// validateCreateInput 验证创建输入
func (s *TodoService) validateCreateInput(input *models.CreateTodoInput) error {
	// 标题验证
	if strings.TrimSpace(input.Title) == "" {
		return customerrors.ErrTitleRequired
	}

	if len(input.Title) > 255 {
		return customerrors.ErrTitleTooLong
	}

	// 分类验证。空字符串是允许的，会在后续设置为默认值 "life"
	// 如果分类不合法，返回错误
	if input.Category != "" {
		validCategories := []string{"work", "study", "life"}
		if !contains(validCategories, input.Category) {
			return customerrors.ErrInvalidCategory(input.Category)
		}
	}

	// 优先级验证
	if input.Priority < 0 || input.Priority > 5 {
		return customerrors.ErrInvalidPriority
	}

	return nil
}

// contains 辅助函数：检查item是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// CreateTodo 创建待办事项
func (s *TodoService) CreateTodo(input *models.CreateTodoInput) (*models.Todo, error) {
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	todo := &models.Todo{
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Category:    input.Category,
		Priority:    input.Priority,
	}

	// Category 为空时，默认设置为 "life"
	if todo.Category == "" {
		todo.Category = "life"
	}
	// Priority 超出范围时，修正为合法范围
	if todo.Priority < 0 {
		todo.Priority = 0
	}
	if todo.Priority > 5 {
		todo.Priority = 5
	}

	// 数据库插入
	if err := todo.Create(); err != nil {
		return nil, customerrors.WrapCreateError(err)
	}

	return todo, nil
}

// GetAllTodos 获取所有待办事项
func (s *TodoService) GetAllTodos(category string, sortBy string) ([]models.Todo, error) {
	// 验证分类参数，避免调接口时故意传不正确的category
	if category != "" && category != "all" {
		validCategories := []string{"work", "study", "life"}
		if !contains(validCategories, category) {
			return nil, customerrors.ErrInvalidCategory(category)
		}
	}

	// 验证排序参数
	if sortBy != "" && sortBy != "priority" && sortBy != "created_at" {
		return nil, customerrors.ErrInvalidSort(sortBy)
	}

	todos, err := models.GetAll(category, sortBy)
	if err != nil {
		return nil, customerrors.WrapQueryError(err)
	}

	return todos, nil
}

// GetTodoByID 根据ID获取待办事项
func (s *TodoService) GetTodoByID(id uint) (*models.Todo, error) {
	if id == 0 {
		return nil, customerrors.ErrInvalidID
	}

	todo, err := models.GetByID(id)
	if err != nil {
		return nil, customerrors.ErrTodoNotFoundWithID(id)
	}

	return todo, nil
}

// validateUpdateInput 验证编辑输入
func (s *TodoService) validateUpdateInput(input *models.UpdateTodoInput) error {
	// 标题验证
	if strings.TrimSpace(input.Title) == "" {
		return customerrors.ErrTitleRequired
	}

	if len(input.Title) > 255 {
		return customerrors.ErrTitleTooLong
	}

	// 分类验证（编辑时分类是必填的）
	validCategories := []string{"work", "study", "life"}
	if !contains(validCategories, input.Category) {
		return customerrors.ErrInvalidCategory(input.Category)
	}

	// 优先级验证
	if input.Priority < 0 || input.Priority > 5 {
		return customerrors.ErrInvalidPriority
	}

	// 版本号验证
	if input.Version < 0 {
		return customerrors.ErrInvalidVersion
	}

	return nil
}

// VersionConflictError 版本冲突错误
type VersionConflictError struct {
	Message         string
	CurrentVersion  int
	ProvidedVersion int
	LatestData      *models.Todo
}

func (e *VersionConflictError) Error() string {
	return e.Message
}

// UpdateTodo 更新待办事项
// 可以更新标题、描述、分类、优先级，使用乐观锁保护
func (s *TodoService) UpdateTodo(id uint, input *models.UpdateTodoInput) (*models.Todo, error) {
	// 验证 ID
	if id == 0 {
		return nil, errors.New("invalid id: id must be greater than 0")
	}

	// 验证输入
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// 先查询当前记录是否存在
	existingTodo, err := models.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found with id %d", id)
	}

	// 乐观锁冲突检测
	if existingTodo.Version != input.Version {
		return nil, &VersionConflictError{
			Message:         "version conflict: data has been modified by another user",
			CurrentVersion:  existingTodo.Version,
			ProvidedVersion: input.Version,
			LatestData:      existingTodo,
		}
	}

	// 清理数据（去除首尾空格）
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)

	// 调用 Model 层更新
	todo := &models.Todo{}
	if err := todo.Update(id, title, description, input.Category, input.Priority, input.Version); err != nil {
		// 处理乐观锁冲突（双重检查）
		if strings.Contains(err.Error(), "version conflict") {
			// 获取最新数据返回给客户端
			latestTodo, _ := models.GetByID(id)
			return nil, &VersionConflictError{
				Message:         err.Error(),
				CurrentVersion:  latestTodo.Version,
				ProvidedVersion: input.Version,
				LatestData:      latestTodo,
			}
		}
		return nil, customerrors.WrapUpdateError(err)
	}

	// 返回更新后的数据
	updatedTodo, err := models.GetByID(id)
	if err != nil {
		return nil, customerrors.WrapGetError(err)
	}

	return updatedTodo, nil
}

// UpdateTodoStatus 更新待办事项状态
func (s *TodoService) UpdateTodoStatus(id uint, input *models.UpdateStatusInput) (*models.Todo, error) {
	if id == 0 {
		return nil, customerrors.ErrInvalidID
	}

	if input.Version < 0 {
		return nil, customerrors.ErrInvalidVersion
	}

	// 先查询当前记录是否存在
	existingTodo, err := models.GetByID(id)
	if err != nil {
		return nil, customerrors.ErrTodoNotFoundWithID(id)
	}

	// 乐观锁冲突检测
	if existingTodo.Version != input.Version {
		return nil, &VersionConflictError{
			Message:         "version conflict: data has been modified by another user",
			CurrentVersion:  existingTodo.Version,
			ProvidedVersion: input.Version,
			LatestData:      existingTodo,
		}
	}

	// 5. 调用 Model 层更新状态
	todo := &models.Todo{}
	if err := todo.UpdateStatus(id, input.Completed, input.Version); err != nil {
		// 处理乐观锁冲突（双重检查）
		if strings.Contains(err.Error(), "version conflict") {
			// 获取最新数据返回给客户端
			latestTodo, _ := models.GetByID(id)
			return nil, &VersionConflictError{
				Message:         err.Error(),
				CurrentVersion:  latestTodo.Version,
				ProvidedVersion: input.Version,
				LatestData:      latestTodo,
			}
		}
		return nil, fmt.Errorf("failed to update todo status: %w", err)
	}

	// 6. 返回更新后的数据
	updatedTodo, err := models.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated todo: %w", err)
	}

	return updatedTodo, nil
}

// DeleteTodo 删除待办事项
func (s *TodoService) DeleteTodo(id uint) error {
	// 验证 ID
	if id == 0 {
		return customerrors.ErrInvalidID
	}

	// 先检查是否存在
	_, err := models.GetByID(id)
	if err != nil {
		return customerrors.ErrTodoNotFoundWithID(id)
	}

	// 调用 Model 层删除
	if err := models.Delete(id); err != nil {
		return customerrors.WrapDeleteError(err)
	}

	return nil
}
