package services

import (
	"backend/config"
	"backend/models"
	"fmt"
	"testing"
)

var service *TodoService

// TestMain 在所有测试前初始化
func TestMain(m *testing.M) {
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	fmt.Println("Database connected for service testing")

	// 创建服务实例
	service = NewTodoService()

	// 运行所有测试
	m.Run()
}

// TestCreateTodo 测试创建待办事项
func TestCreateTodo(t *testing.T) {
	t.Run("创建带完整信息的待办事项", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title:       "服务层测试任务1",
			Description: "测试服务层的创建功能",
			Category:    "work",
			Priority:    5,
		}

		todo, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		if todo.ID == 0 {
			t.Error("ID 应该不为 0")
		}
		if todo.Title != "服务层测试任务1" {
			t.Errorf("标题不匹配: %s", todo.Title)
		}
		if todo.Category != "work" {
			t.Errorf("分类不匹配: %s", todo.Category)
		}
		if todo.Priority != 5 {
			t.Errorf("优先级不匹配: %d", todo.Priority)
		}

		t.Logf("✅ 创建成功，ID: %d", todo.ID)
	})

	t.Run("创建使用默认值的待办事项", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title: "服务层测试任务2（默认值）",
		}

		todo, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		if todo.Category != "life" {
			t.Errorf("默认分类应该是 life，实际: %s", todo.Category)
		}
		if todo.Priority != 0 {
			t.Errorf("默认优先级应该是 0，实际: %d", todo.Priority)
		}

		t.Logf("✅ 创建成功，默认分类: %s, 默认优先级: %d", todo.Category, todo.Priority)
	})

	t.Run("创建时自动去除首尾空格", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title:       "  空格测试任务  ",
			Description: "  空格描述  ",
			Category:    "study",
		}

		todo, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		if todo.Title != "空格测试任务" {
			t.Errorf("标题应该去除空格，实际: '%s'", todo.Title)
		}
		if todo.Description != "空格描述" {
			t.Errorf("描述应该去除空格，实际: '%s'", todo.Description)
		}

		t.Log("✅ 自动去除空格功能正常")
	})

	t.Run("验证：标题为空应该失败", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title: "",
		}

		_, err := service.CreateTodo(input)
		if err == nil {
			t.Error("标题为空应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截空标题: %v", err)
	})

	t.Run("验证：标题只有空格应该失败", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title: "   ",
		}

		_, err := service.CreateTodo(input)
		if err == nil {
			t.Error("标题只有空格应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截空白标题: %v", err)
	})

	t.Run("验证：无效分类应该失败", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title:    "测试任务",
			Category: "invalid_category",
		}

		_, err := service.CreateTodo(input)
		if err == nil {
			t.Error("无效分类应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效分类: %v", err)
	})

	t.Run("验证：优先级超出范围应该失败", func(t *testing.T) {
		input := &models.CreateTodoInput{
			Title:    "测试任务",
			Priority: 10,
		}

		_, err := service.CreateTodo(input)
		if err == nil {
			t.Error("优先级超出范围应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效优先级: %v", err)
	})
}

// TestGetAllTodos 测试获取所有待办事项
func TestGetAllTodos(t *testing.T) {
	t.Run("获取所有待办事项", func(t *testing.T) {
		todos, err := service.GetAllTodos("", "")
		if err != nil {
			t.Errorf("获取失败: %v", err)
			return
		}

		t.Logf("✅ 成功获取 %d 条待办事项", len(todos))
	})

	t.Run("按分类筛选", func(t *testing.T) {
		todos, err := service.GetAllTodos("work", "")
		if err != nil {
			t.Errorf("获取失败: %v", err)
			return
		}

		for _, todo := range todos {
			if todo.Category != "work" {
				t.Errorf("分类筛选失败，期望 work，实际 %s", todo.Category)
			}
		}

		t.Logf("✅ work 分类筛选成功，共 %d 条", len(todos))
	})

	t.Run("按优先级排序", func(t *testing.T) {
		todos, err := service.GetAllTodos("", "priority")
		if err != nil {
			t.Errorf("获取失败: %v", err)
			return
		}

		if len(todos) > 1 {
			for i := 0; i < len(todos)-1; i++ {
				if todos[i].Priority < todos[i+1].Priority {
					t.Error("优先级排序错误")
				}
			}
		}

		t.Logf("✅ 优先级排序成功")
	})

	t.Run("验证：无效分类应该失败", func(t *testing.T) {
		_, err := service.GetAllTodos("invalid", "")
		if err == nil {
			t.Error("无效分类应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效分类: %v", err)
	})

	t.Run("验证：无效排序参数应该失败", func(t *testing.T) {
		_, err := service.GetAllTodos("", "invalid_sort")
		if err == nil {
			t.Error("无效排序参数应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效排序参数: %v", err)
	})
}

// TestGetTodoByID 测试根据ID获取
func TestGetTodoByID(t *testing.T) {
	t.Run("获取存在的待办事项", func(t *testing.T) {
		// 先创建一个
		input := &models.CreateTodoInput{
			Title:    "用于ID查询的服务层测试",
			Category: "work",
			Priority: 3,
		}
		created, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		// 查询
		todo, err := service.GetTodoByID(created.ID)
		if err != nil {
			t.Errorf("查询失败: %v", err)
			return
		}

		if todo.ID != created.ID {
			t.Error("ID 不匹配")
		}

		t.Logf("✅ 成功查询 ID=%d 的待办事项", todo.ID)
	})

	t.Run("验证：ID为0应该失败", func(t *testing.T) {
		_, err := service.GetTodoByID(0)
		if err == nil {
			t.Error("ID为0应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效ID: %v", err)
	})

	t.Run("验证：不存在的ID应该失败", func(t *testing.T) {
		_, err := service.GetTodoByID(999999)
		if err == nil {
			t.Error("不存在的ID应该返回错误")
			return
		}

		t.Logf("✅ 正确处理不存在的ID: %v", err)
	})
}

// TestUpdateTodo 测试更新待办事项（编辑功能）
func TestUpdateTodo(t *testing.T) {
	t.Run("正常编辑待办事项", func(t *testing.T) {
		// 创建待办事项
		createInput := &models.CreateTodoInput{
			Title:       "原始标题",
			Description: "原始描述",
			Category:    "work",
			Priority:    3,
		}
		created, err := service.CreateTodo(createInput)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		originalVersion := created.Version
		t.Logf("创建后: Title=%s, Category=%s, Priority=%d, Version=%d",
			created.Title, created.Category, created.Priority, originalVersion)

		// 编辑待办事项
		updateInput := &models.UpdateTodoInput{
			Title:       "修改后的标题",
			Description: "修改后的描述",
			Category:    "study",
			Priority:    5,
			Version:     originalVersion,
		}
		updated, err := service.UpdateTodo(created.ID, updateInput)
		if err != nil {
			t.Errorf("更新失败: %v", err)
			return
		}

		// 验证更新结果
		if updated.Title != "修改后的标题" {
			t.Errorf("标题未更新，期望: 修改后的标题，实际: %s", updated.Title)
		}
		if updated.Description != "修改后的描述" {
			t.Errorf("描述未更新")
		}
		if updated.Category != "study" {
			t.Errorf("分类未更新，期望: study，实际: %s", updated.Category)
		}
		if updated.Priority != 5 {
			t.Errorf("优先级未更新，期望: 5，实际: %d", updated.Priority)
		}
		if updated.Version != originalVersion+1 {
			t.Errorf("版本号应该为 %d，实际为 %d", originalVersion+1, updated.Version)
		}

		t.Logf("✅ 编辑成功: Title=%s, Category=%s, Priority=%d, Version=%d",
			updated.Title, updated.Category, updated.Priority, updated.Version)

		// 清理
		service.DeleteTodo(created.ID)
	})

	t.Run("编辑时自动去除空格", func(t *testing.T) {
		// 创建待办事项
		createInput := &models.CreateTodoInput{
			Title:    "测试任务",
			Category: "work",
		}
		created, err := service.CreateTodo(createInput)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		// 编辑时带空格
		updateInput := &models.UpdateTodoInput{
			Title:       "  空格测试  ",
			Description: "  空格描述  ",
			Category:    "study",
			Priority:    3,
			Version:     created.Version,
		}
		updated, err := service.UpdateTodo(created.ID, updateInput)
		if err != nil {
			t.Errorf("更新失败: %v", err)
			return
		}

		if updated.Title != "空格测试" {
			t.Errorf("标题应该去除空格，实际: '%s'", updated.Title)
		}
		if updated.Description != "空格描述" {
			t.Errorf("描述应该去除空格，实际: '%s'", updated.Description)
		}

		t.Log("✅ 自动去除空格功能正常")

		// 清理
		service.DeleteTodo(created.ID)
	})

	t.Run("乐观锁：编辑时版本冲突", func(t *testing.T) {
		// 创建待办事项
		createInput := &models.CreateTodoInput{
			Title:    "用于乐观锁测试",
			Category: "work",
			Priority: 3,
		}
		created, err := service.CreateTodo(createInput)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		// 第一次编辑（模拟用户A）
		updateInput1 := &models.UpdateTodoInput{
			Title:    "用户A的修改",
			Category: "study",
			Priority: 4,
			Version:  0,
		}
		_, err = service.UpdateTodo(created.ID, updateInput1)
		if err != nil {
			t.Errorf("第一次编辑失败: %v", err)
			return
		}
		t.Log("用户A 编辑成功，版本号 0 -> 1")

		// 第二次编辑使用旧版本号（模拟用户B使用过期的版本号）
		updateInput2 := &models.UpdateTodoInput{
			Title:    "用户B的修改",
			Category: "life",
			Priority: 5,
			Version:  0, // 使用旧版本号
		}
		_, err = service.UpdateTodo(created.ID, updateInput2)
		if err == nil {
			t.Error("使用旧版本号编辑应该失败")
			return
		}

		// 检查是否是版本冲突错误
		if _, ok := err.(*VersionConflictError); !ok {
			t.Errorf("应该返回 VersionConflictError，实际: %T", err)
		}

		t.Logf("✅ 乐观锁正常工作（编辑场景）: %v", err)

		// 验证数据没有被覆盖
		final, _ := service.GetTodoByID(created.ID)
		if final.Title != "用户A的修改" {
			t.Error("数据被错误覆盖")
		}

		// 清理
		service.DeleteTodo(created.ID)
	})

	t.Run("验证：标题为空应该失败", func(t *testing.T) {
		// 创建待办事项
		createInput := &models.CreateTodoInput{
			Title:    "测试任务",
			Category: "work",
		}
		created, err := service.CreateTodo(createInput)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		// 尝试编辑为空标题
		updateInput := &models.UpdateTodoInput{
			Title:    "",
			Category: "work",
			Priority: 3,
			Version:  created.Version,
		}
		_, err = service.UpdateTodo(created.ID, updateInput)
		if err == nil {
			t.Error("标题为空应该返回错误")
		}

		t.Logf("✅ 正确拦截空标题: %v", err)

		// 清理
		service.DeleteTodo(created.ID)
	})

	t.Run("验证：无效分类应该失败", func(t *testing.T) {
		// 创建待办事项
		createInput := &models.CreateTodoInput{
			Title:    "测试任务",
			Category: "work",
		}
		created, err := service.CreateTodo(createInput)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		// 尝试编辑为无效分类
		updateInput := &models.UpdateTodoInput{
			Title:    "测试任务",
			Category: "invalid",
			Priority: 3,
			Version:  created.Version,
		}
		_, err = service.UpdateTodo(created.ID, updateInput)
		if err == nil {
			t.Error("无效分类应该返回错误")
		}

		t.Logf("✅ 正确拦截无效分类: %v", err)

		// 清理
		service.DeleteTodo(created.ID)
	})

	t.Run("验证：ID不存在应该失败", func(t *testing.T) {
		updateInput := &models.UpdateTodoInput{
			Title:    "测试任务",
			Category: "work",
			Priority: 3,
			Version:  0,
		}
		_, err := service.UpdateTodo(999999, updateInput)
		if err == nil {
			t.Error("不存在的ID应该返回错误")
		}

		t.Logf("✅ 正确处理不存在的ID: %v", err)
	})
}

// TestUpdateTodoStatus 测试更新状态
func TestUpdateTodoStatus(t *testing.T) {
	t.Run("正常更新状态", func(t *testing.T) {
		// 创建待办事项
		input := &models.CreateTodoInput{
			Title:    "用于状态更新的服务层测试",
			Category: "study",
			Priority: 4,
		}
		created, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		originalVersion := created.Version
		t.Logf("创建后版本号: %d", originalVersion)

		// 更新为已完成
		updateInput := &models.UpdateStatusInput{
			Completed: true,
			Version:   originalVersion,
		}
		updated, err := service.UpdateTodoStatus(created.ID, updateInput)
		if err != nil {
			t.Errorf("更新失败: %v", err)
			return
		}

		if !updated.Completed {
			t.Error("状态应该已更新为完成")
		}
		if updated.Version != originalVersion+1 {
			t.Errorf("版本号应该为 %d，实际为 %d", originalVersion+1, updated.Version)
		}

		t.Logf("✅ 成功更新状态，版本号 %d -> %d", originalVersion, updated.Version)
	})

	t.Run("乐观锁：版本冲突检测（前置检查）", func(t *testing.T) {
		// 创建待办事项
		input := &models.CreateTodoInput{
			Title:    "用于乐观锁前置检查测试",
			Category: "work",
			Priority: 5,
		}
		created, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		// 第一次更新
		updateInput1 := &models.UpdateStatusInput{
			Completed: true,
			Version:   0,
		}
		_, err = service.UpdateTodoStatus(created.ID, updateInput1)
		if err != nil {
			t.Errorf("第一次更新失败: %v", err)
			return
		}
		t.Log("第一次更新成功，版本号 0 -> 1")

		// 第二次更新使用旧版本号（应该失败）
		updateInput2 := &models.UpdateStatusInput{
			Completed: false,
			Version:   0, // 使用旧版本号
		}
		_, err = service.UpdateTodoStatus(created.ID, updateInput2)
		if err == nil {
			t.Error("使用旧版本号应该失败")
			return
		}

		// 检查是否是版本冲突错误
		if _, ok := err.(*VersionConflictError); !ok {
			t.Errorf("应该返回 VersionConflictError，实际: %T", err)
		}

		t.Logf("✅ 乐观锁前置检查正常工作: %v", err)
	})

	t.Run("验证：ID为0应该失败", func(t *testing.T) {
		updateInput := &models.UpdateStatusInput{
			Completed: true,
			Version:   0,
		}
		_, err := service.UpdateTodoStatus(0, updateInput)
		if err == nil {
			t.Error("ID为0应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效ID: %v", err)
	})

	t.Run("验证：负数版本号应该失败", func(t *testing.T) {
		updateInput := &models.UpdateStatusInput{
			Completed: true,
			Version:   -1,
		}
		_, err := service.UpdateTodoStatus(1, updateInput)
		if err == nil {
			t.Error("负数版本号应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效版本号: %v", err)
	})
}

// TestDeleteTodo 测试删除待办事项
func TestDeleteTodo(t *testing.T) {
	t.Run("删除存在的待办事项", func(t *testing.T) {
		// 创建待办事项
		input := &models.CreateTodoInput{
			Title:    "用于删除测试的服务层任务",
			Category: "life",
			Priority: 1,
		}
		created, err := service.CreateTodo(input)
		if err != nil {
			t.Errorf("创建失败: %v", err)
			return
		}

		todoID := created.ID
		t.Logf("创建了 ID=%d 的待办事项", todoID)

		// 删除
		err = service.DeleteTodo(todoID)
		if err != nil {
			t.Errorf("删除失败: %v", err)
			return
		}

		// 验证已删除
		_, err = service.GetTodoByID(todoID)
		if err == nil {
			t.Error("删除后查询应该失败")
			return
		}

		t.Logf("✅ 成功删除 ID=%d 的待办事项", todoID)
	})

	t.Run("验证：ID为0应该失败", func(t *testing.T) {
		err := service.DeleteTodo(0)
		if err == nil {
			t.Error("ID为0应该返回错误")
			return
		}

		t.Logf("✅ 正确拦截无效ID: %v", err)
	})

	t.Run("验证：删除不存在的待办事项应该失败", func(t *testing.T) {
		err := service.DeleteTodo(999999)
		if err == nil {
			t.Error("删除不存在的待办事项应该返回错误")
			return
		}

		t.Logf("✅ 正确处理不存在的记录: %v", err)
	})
}

// TestCompleteServiceWorkflow 测试完整服务层工作流
func TestCompleteServiceWorkflow(t *testing.T) {
	t.Run("完整的服务层CRUD+编辑工作流", func(t *testing.T) {
		// 1. 创建
		createInput := &models.CreateTodoInput{
			Title:       "完整工作流测试",
			Description: "测试服务层的完整流程",
			Category:    "work",
			Priority:    5,
		}
		created, err := service.CreateTodo(createInput)
		if err != nil {
			t.Fatalf("❌ 创建失败: %v", err)
		}
		t.Logf("✅ 1. 创建成功，ID=%d, Version=%d", created.ID, created.Version)

		// 2. 查询
		retrieved, err := service.GetTodoByID(created.ID)
		if err != nil {
			t.Fatalf("❌ 查询失败: %v", err)
		}
		t.Logf("✅ 2. 查询成功: %s", retrieved.Title)

		// 3. 编辑（新增）
		editInput := &models.UpdateTodoInput{
			Title:       "修改后的标题",
			Description: "修改后的描述",
			Category:    "study",
			Priority:    4,
			Version:     retrieved.Version,
		}
		edited, err := service.UpdateTodo(created.ID, editInput)
		if err != nil {
			t.Fatalf("❌ 编辑失败: %v", err)
		}
		if edited.Title != "修改后的标题" {
			t.Error("❌ 标题未更新")
		}
		t.Logf("✅ 3. 编辑成功，版本号: %d -> %d", retrieved.Version, edited.Version)

		// 4. 更新状态
		updateStatusInput := &models.UpdateStatusInput{
			Completed: true,
			Version:   edited.Version,
		}
		statusUpdated, err := service.UpdateTodoStatus(created.ID, updateStatusInput)
		if err != nil {
			t.Fatalf("❌ 更新状态失败: %v", err)
		}
		if !statusUpdated.Completed {
			t.Error("❌ 状态未更新")
		}
		t.Logf("✅ 4. 更新状态成功，版本号: %d -> %d", edited.Version, statusUpdated.Version)

		// 5. 删除
		err = service.DeleteTodo(created.ID)
		if err != nil {
			t.Fatalf("❌ 删除失败: %v", err)
		}
		t.Log("✅ 5. 删除成功")

		// 6. 验证删除
		_, err = service.GetTodoByID(created.ID)
		if err == nil {
			t.Error("❌ 删除后不应该能查询到")
		}
		t.Log("✅ 6. 验证删除成功")

		t.Log("🎉 完整服务层工作流测试全部通过（包含编辑功能）！")
	})
}
