package models

import (
	"backend/config"
	"errors"
	"time"
)

// Todo 待办事项模型结构体
type Todo struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title" binding:"required,min=1,max=255"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:enum('work','study','life');default:'life'" json:"category" binding:"omitempty,oneof=work study life"`
	Priority    int       `gorm:"default:0" json:"priority" binding:"omitempty,min=0,max=5"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	Version     int       `gorm:"default:0" json:"version"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Todo) TableName() string {
	return "todos"
}

// CreateTodoInput 创建待办事项的输入结构
type CreateTodoInput struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description"` // 描述非必须
	Category    string `json:"category" binding:"omitempty,oneof=work study life"`
	Priority    int    `json:"priority" binding:"omitempty,min=0,max=5"`
}

// UpdateStatusInput 更新状态的输入结构
type UpdateStatusInput struct {
	Completed bool `json:"completed"`
	Version   int  `json:"version" binding:"required"`
}

// Create 创建待办事项
func (t *Todo) Create() error {
	// 设置默认值
	if t.Category == "" {
		t.Category = "life"
	}
	if t.Priority < 0 {
		t.Priority = 0
	}

	result := config.DB.Create(t)
	return result.Error
}

// GetAll 获取所有待办事项
// 支持按分类筛选和排序
func GetAll(category string, sortBy string) ([]Todo, error) {
	var todos []Todo
	query := config.DB.Model(&Todo{})

	// 分类筛选
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	// 排序
	switch sortBy {
	case "priority":
		query = query.Order("priority DESC, created_at DESC")
	default:
		// 默认按创建时间降序（包括 sortBy="created_at" 和空值的情况）
		query = query.Order("created_at DESC")
	}

	err := query.Find(&todos).Error
	return todos, err
}

// GetByID 根据ID获取待办事项
func GetByID(id uint) (*Todo, error) {
	var todo Todo
	result := config.DB.First(&todo, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

// UpdateStatus 更新完成状态
func (t *Todo) UpdateStatus(id uint, completed bool, version int) error {
	result := config.DB.Model(&Todo{}).
		Where("id = ? AND version = ?", id, version). // 假如用户同时多设备点击更新完成状态，那么只有一个设备会成功，另一个设备在where语句查不出来
		Updates(map[string]interface{}{
			"completed": completed,
			"version":   version + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	// 检查是否有行被更新（乐观锁检查），查不出来就会影响行数为0
	if result.RowsAffected == 0 {
		return errors.New("version conflict: data has been modified by another user")
	}

	return nil
}

// Delete 删除待办事项，硬删除，因为待办事项一般不需要找回
func Delete(id uint) error {
	result := config.DB.Delete(&Todo{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}
