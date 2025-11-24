package errors

import (
	"errors"
	"fmt"
)

// 验证错误
var (
	ErrInvalidID       = errors.New("invalid id: id must be greater than 0")
	ErrTitleRequired   = errors.New("title is required and cannot be empty")
	ErrTitleTooLong    = errors.New("title cannot exceed 255 characters")
	ErrInvalidPriority = errors.New("priority must be between 0 and 5")
	ErrInvalidVersion  = errors.New("invalid version: version must be non-negative")
)

// 业务错误
var (
	ErrTodoNotFound    = errors.New("todo not found")
	ErrVersionConflict = errors.New("version conflict: data has been modified by another user")
)

// 数据库错误
var (
	ErrDatabaseConnection = errors.New("failed to connect to database")
	ErrDatabaseInit       = errors.New("failed to initialize database")
)

// ErrInvalidCategory 无效分类错误
func ErrInvalidCategory(category string) error {
	return fmt.Errorf("invalid category: %s, must be one of: work, study, life", category)
}

// ErrInvalidSort 无效排序参数错误
func ErrInvalidSort(sortBy string) error {
	return fmt.Errorf("invalid sort parameter: %s, must be: priority or created_at", sortBy)
}

// ErrTodoNotFoundWithID 待办事项未找到（带ID）
func ErrTodoNotFoundWithID(id uint) error {
	return fmt.Errorf("%w: id=%d", ErrTodoNotFound, id)
}

// 包装错误函数（保留原始错误链）
func WrapCreateError(err error) error {
	return fmt.Errorf("failed to create todo: %w", err)
}

func WrapUpdateError(err error) error {
	return fmt.Errorf("failed to update todo: %w", err)
}

func WrapDeleteError(err error) error {
	return fmt.Errorf("failed to delete todo: %w", err)
}

func WrapQueryError(err error) error {
	return fmt.Errorf("failed to query todos: %w", err)
}

func WrapGetError(err error) error {
	return fmt.Errorf("failed to get todo: %w", err)
}
