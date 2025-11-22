package models

import (
	"backend/config"
	"fmt"
	"testing"
)

// TestMain 在所有测试前初始化数据库连接
func TestMain(m *testing.M) {
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	fmt.Println("Database connected for testing")

	// 运行所有测试
	m.Run()
}

// TestCreate 测试创建待办事项
func TestCreate(t *testing.T) {
	t.Run("创建带完整信息的待办事项", func(t *testing.T) {
		todo := &Todo{
			Title:       "测试任务1",
			Description: "这是一个测试任务的详细描述",
			Category:    "work",
			Priority:    5,
		}

		err := todo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		if todo.ID == 0 {
			t.Error("创建后 ID 应该不为 0")
		}
		if todo.Version != 0 {
			t.Error("新创建的待办事项版本号应该为 0")
		}
		if todo.Completed != false {
			t.Error("新创建的待办事项应该未完成")
		}

		t.Logf("✅ 成功创建待办事项，ID: %d", todo.ID)
	})

	t.Run("创建不带描述的待办事项", func(t *testing.T) {
		todo := &Todo{
			Title:    "测试任务2（无描述）",
			Category: "study",
			Priority: 3,
		}

		err := todo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		if todo.Description != "" {
			t.Logf("描述为空字符串: '%s'", todo.Description)
		}

		t.Logf("✅ 成功创建无描述待办事项，ID: %d", todo.ID)
	})

	t.Run("创建使用默认分类的待办事项", func(t *testing.T) {
		todo := &Todo{
			Title:    "测试任务3（默认分类）",
			Priority: 2,
		}

		err := todo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		if todo.Category != "life" {
			t.Errorf("默认分类应该是 'life'，实际是: %s", todo.Category)
		}

		t.Logf("✅ 成功创建待办事项，默认分类: %s", todo.Category)
	})
}

// TestGetAll 测试获取所有待办事项
func TestGetAll(t *testing.T) {
	t.Run("获取所有待办事项（无筛选）", func(t *testing.T) {
		todos, err := GetAll("", "")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		t.Logf("✅ 成功获取 %d 条待办事项", len(todos))

		if len(todos) > 0 {
			t.Logf("第一条: ID=%d, Title=%s, Category=%s, Priority=%d",
				todos[0].ID, todos[0].Title, todos[0].Category, todos[0].Priority)
		}
	})

	t.Run("按分类筛选 - work", func(t *testing.T) {
		todos, err := GetAll("work", "")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		for _, todo := range todos {
			if todo.Category != "work" {
				t.Errorf("筛选结果应该都是 work 分类，但发现: %s", todo.Category)
			}
		}

		t.Logf("✅ 成功获取 work 分类的 %d 条待办事项", len(todos))
	})

	t.Run("按分类筛选 - study", func(t *testing.T) {
		todos, err := GetAll("study", "")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		for _, todo := range todos {
			if todo.Category != "study" {
				t.Errorf("筛选结果应该都是 study 分类，但发现: %s", todo.Category)
			}
		}

		t.Logf("✅ 成功获取 study 分类的 %d 条待办事项", len(todos))
	})

	t.Run("按分类筛选 - life", func(t *testing.T) {
		todos, err := GetAll("life", "")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		for _, todo := range todos {
			if todo.Category != "life" {
				t.Errorf("筛选结果应该都是 life 分类，但发现: %s", todo.Category)
			}
		}

		t.Logf("✅ 成功获取 life 分类的 %d 条待办事项", len(todos))
	})
}

// TestGetAllWithSort 测试排序功能
func TestGetAllWithSort(t *testing.T) {
	t.Run("按优先级排序", func(t *testing.T) {
		todos, err := GetAll("", "priority")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		if len(todos) > 1 {
			// 验证降序排列
			for i := 0; i < len(todos)-1; i++ {
				if todos[i].Priority < todos[i+1].Priority {
					t.Errorf("优先级排序错误: todos[%d].Priority=%d < todos[%d].Priority=%d",
						i, todos[i].Priority, i+1, todos[i+1].Priority)
				}
			}
		}

		t.Logf("✅ 成功按优先级排序，共 %d 条", len(todos))
		if len(todos) > 0 {
			t.Logf("第一条优先级: %d, 最后一条优先级: %d",
				todos[0].Priority, todos[len(todos)-1].Priority)
		}
	})

	t.Run("按创建时间排序", func(t *testing.T) {
		todos, err := GetAll("", "created_at")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		if len(todos) > 1 {
			// 验证降序排列（最新的在前）
			for i := 0; i < len(todos)-1; i++ {
				if todos[i].CreatedAt.Before(todos[i+1].CreatedAt) {
					t.Errorf("时间排序错误: todos[%d] 早于 todos[%d]", i, i+1)
				}
			}
		}

		t.Logf("✅ 成功按创建时间排序，共 %d 条", len(todos))
		if len(todos) > 0 {
			t.Logf("第一条创建时间: %s", todos[0].CreatedAt.Format("2006-01-02 15:04:05"))
		}
	})

	t.Run("组合：按分类筛选并按优先级排序", func(t *testing.T) {
		todos, err := GetAll("work", "priority")
		if err != nil {
			t.Errorf("获取待办事项失败: %v", err)
			return
		}

		// 验证分类
		for _, todo := range todos {
			if todo.Category != "work" {
				t.Errorf("分类筛选失败: 期望 work，实际 %s", todo.Category)
			}
		}

		// 验证排序
		if len(todos) > 1 {
			for i := 0; i < len(todos)-1; i++ {
				if todos[i].Priority < todos[i+1].Priority {
					t.Errorf("优先级排序错误")
				}
			}
		}

		t.Logf("✅ 成功组合筛选和排序，work 分类共 %d 条", len(todos))
	})
}

// TestGetByID 测试根据ID查询
func TestGetByID(t *testing.T) {
	t.Run("查询存在的待办事项", func(t *testing.T) {
		// 先创建一个
		newTodo := &Todo{
			Title:       "用于ID查询的测试任务",
			Description: "测试 GetByID 方法",
			Category:    "work",
			Priority:    4,
		}
		err := newTodo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		// 查询
		todo, err := GetByID(newTodo.ID)
		if err != nil {
			t.Errorf("查询待办事项失败: %v", err)
			return
		}

		if todo.ID != newTodo.ID {
			t.Errorf("ID 不匹配: 期望 %d，实际 %d", newTodo.ID, todo.ID)
		}
		if todo.Title != newTodo.Title {
			t.Errorf("Title 不匹配")
		}

		t.Logf("✅ 成功查询 ID=%d 的待办事项: %s", todo.ID, todo.Title)
	})

	t.Run("查询不存在的待办事项", func(t *testing.T) {
		todo, err := GetByID(999999)
		if err == nil {
			t.Error("查询不存在的 ID 应该返回错误")
			return
		}
		if todo != nil {
			t.Error("查询不存在的 ID 应该返回 nil")
		}

		t.Logf("✅ 正确处理不存在的 ID，返回错误: %v", err)
	})
}

// TestUpdate 测试更新待办事项（编辑功能）
func TestUpdate(t *testing.T) {
	t.Run("正常更新待办事项", func(t *testing.T) {
		// 创建待办事项
		newTodo := &Todo{
			Title:       "原始标题",
			Description: "原始描述",
			Category:    "work",
			Priority:    3,
		}
		err := newTodo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		originalVersion := newTodo.Version
		t.Logf("创建后的版本号: %d", originalVersion)

		// 更新待办事项
		err = newTodo.Update(
			newTodo.ID,
			"修改后的标题",
			"修改后的描述",
			"study",
			5,
			originalVersion,
		)
		if err != nil {
			t.Errorf("更新失败: %v", err)
			return
		}

		// 查询验证
		updated, err := GetByID(newTodo.ID)
		if err != nil {
			t.Errorf("查询失败: %v", err)
			return
		}

		if updated.Title != "修改后的标题" {
			t.Errorf("标题应该已更新，实际: %s", updated.Title)
		}
		if updated.Description != "修改后的描述" {
			t.Errorf("描述应该已更新")
		}
		if updated.Category != "study" {
			t.Errorf("分类应该已更新，实际: %s", updated.Category)
		}
		if updated.Priority != 5 {
			t.Errorf("优先级应该已更新，实际: %d", updated.Priority)
		}
		if updated.Version != originalVersion+1 {
			t.Errorf("版本号应该为 %d，实际为 %d", originalVersion+1, updated.Version)
		}

		t.Logf("✅ 成功更新待办事项，版本号从 %d 变为 %d", originalVersion, updated.Version)
	})

	t.Run("版本冲突测试（编辑场景）", func(t *testing.T) {
		// 创建待办事项
		newTodo := &Todo{
			Title:    "用于乐观锁测试的任务",
			Category: "work",
			Priority: 3,
		}
		err := newTodo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		// 第一次更新（模拟用户A）
		err = newTodo.Update(newTodo.ID, "用户A的修改", "描述A", "study", 4, 0)
		if err != nil {
			t.Errorf("第一次更新失败: %v", err)
			return
		}
		t.Log("用户A 更新成功，版本号 0 -> 1")

		// 第二次更新使用旧版本号（模拟用户B使用过期的版本号）
		err = newTodo.Update(newTodo.ID, "用户B的修改", "描述B", "life", 5, 0)
		if err == nil {
			t.Error("使用过期版本号更新应该失败")
			return
		}

		if err.Error() != "version conflict: data has been modified by another user" {
			t.Errorf("错误信息不匹配: %v", err)
		}

		// 验证数据没有被覆盖
		final, _ := GetByID(newTodo.ID)
		if final.Title != "用户A的修改" {
			t.Error("数据被错误覆盖")
		}

		t.Logf("✅ 乐观锁正常工作（编辑场景），阻止了版本冲突: %v", err)
	})
}

// TestUpdateStatus 测试更新状态（乐观锁）
func TestUpdateStatus(t *testing.T) {
	t.Run("正常更新状态", func(t *testing.T) {
		// 创建待办事项
		newTodo := &Todo{
			Title:    "用于状态更新的测试任务",
			Category: "study",
			Priority: 3,
		}
		err := newTodo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		originalVersion := newTodo.Version
		t.Logf("创建后的版本号: %d", originalVersion)

		// 更新为已完成
		err = newTodo.UpdateStatus(newTodo.ID, true, originalVersion)
		if err != nil {
			t.Errorf("更新状态失败: %v", err)
			return
		}

		// 查询验证
		updated, err := GetByID(newTodo.ID)
		if err != nil {
			t.Errorf("查询失败: %v", err)
			return
		}

		if !updated.Completed {
			t.Error("状态应该已更新为完成")
		}
		if updated.Version != originalVersion+1 {
			t.Errorf("版本号应该为 %d，实际为 %d", originalVersion+1, updated.Version)
		}

		t.Logf("✅ 成功更新状态，版本号从 %d 变为 %d", originalVersion, updated.Version)
	})

	t.Run("版本冲突测试（乐观锁）", func(t *testing.T) {
		// 创建待办事项
		newTodo := &Todo{
			Title:    "用于乐观锁测试的任务",
			Category: "work",
			Priority: 5,
		}
		err := newTodo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		// 第一次更新（模拟用户A）
		err = newTodo.UpdateStatus(newTodo.ID, true, 0)
		if err != nil {
			t.Errorf("第一次更新失败: %v", err)
			return
		}
		t.Log("用户A 更新成功，版本号 0 -> 1")

		// 第二次更新使用旧版本号（模拟用户B使用过期的版本号）
		err = newTodo.UpdateStatus(newTodo.ID, false, 0)
		if err == nil {
			t.Error("使用过期版本号更新应该失败")
			return
		}

		if err.Error() != "version conflict: data has been modified by another user" {
			t.Errorf("错误信息不匹配: %v", err)
		}

		t.Logf("✅ 乐观锁正常工作，阻止了版本冲突: %v", err)
	})
}

// TestDelete 测试删除待办事项
func TestDelete(t *testing.T) {
	t.Run("删除存在的待办事项", func(t *testing.T) {
		// 创建待办事项
		newTodo := &Todo{
			Title:    "用于删除测试的任务",
			Category: "life",
			Priority: 1,
		}
		err := newTodo.Create()
		if err != nil {
			t.Errorf("创建待办事项失败: %v", err)
			return
		}

		todoID := newTodo.ID
		t.Logf("创建了 ID=%d 的待办事项", todoID)

		// 删除
		err = Delete(todoID)
		if err != nil {
			t.Errorf("删除失败: %v", err)
			return
		}

		// 验证已删除
		_, err = GetByID(todoID)
		if err == nil {
			t.Error("删除后查询应该失败")
			return
		}

		t.Logf("✅ 成功删除 ID=%d 的待办事项", todoID)
	})

	t.Run("删除不存在的待办事项", func(t *testing.T) {
		err := Delete(999999)
		if err == nil {
			t.Error("删除不存在的待办事项应该返回错误")
			return
		}

		if err.Error() != "todo not found" {
			t.Errorf("错误信息不匹配: %v", err)
		}

		t.Logf("✅ 正确处理删除不存在的记录: %v", err)
	})
}

// TestCompleteWorkflow 测试完整工作流
func TestCompleteWorkflow(t *testing.T) {
	t.Run("完整的CRUD+编辑工作流", func(t *testing.T) {
		// 1. 创建
		todo := &Todo{
			Title:       "完整工作流测试",
			Description: "测试创建->查询->编辑->更新状态->删除的完整流程",
			Category:    "work",
			Priority:    5,
		}
		err := todo.Create()
		if err != nil {
			t.Fatalf("❌ 创建失败: %v", err)
		}
		t.Logf("✅ 1. 创建成功，ID=%d, Version=%d", todo.ID, todo.Version)

		// 2. 查询
		retrieved, err := GetByID(todo.ID)
		if err != nil {
			t.Fatalf("❌ 查询失败: %v", err)
		}
		t.Logf("✅ 2. 查询成功: %s", retrieved.Title)

		// 3. 编辑
		err = todo.Update(todo.ID, "修改后的标题", "修改后的描述", "study", 4, retrieved.Version)
		if err != nil {
			t.Fatalf("❌ 编辑失败: %v", err)
		}
		t.Log("✅ 3. 编辑成功")

		// 4. 验证编辑
		edited, err := GetByID(todo.ID)
		if err != nil {
			t.Fatalf("❌ 查询编辑后的记录失败: %v", err)
		}
		if edited.Title != "修改后的标题" {
			t.Error("❌ 标题未更新")
		}
		if edited.Version != 1 {
			t.Errorf("❌ 版本号应该为 1，实际为 %d", edited.Version)
		}
		t.Log("✅ 4. 验证编辑成功")

		// 5. 更新状态
		err = todo.UpdateStatus(todo.ID, true, edited.Version)
		if err != nil {
			t.Fatalf("❌ 更新状态失败: %v", err)
		}
		t.Log("✅ 5. 更新状态成功")

		// 6. 验证状态更新
		statusUpdated, err := GetByID(todo.ID)
		if err != nil {
			t.Fatalf("❌ 查询状态更新后的记录失败: %v", err)
		}
		if !statusUpdated.Completed {
			t.Error("❌ 状态未更新")
		}
		if statusUpdated.Version != 2 {
			t.Errorf("❌ 版本号应该为 2，实际为 %d", statusUpdated.Version)
		}
		t.Log("✅ 6. 验证状态更新成功")

		// 7. 删除
		err = Delete(todo.ID)
		if err != nil {
			t.Fatalf("❌ 删除失败: %v", err)
		}
		t.Log("✅ 7. 删除成功")

		// 8. 验证删除
		_, err = GetByID(todo.ID)
		if err == nil {
			t.Error("❌ 删除后不应该能查询到")
		}
		t.Log("✅ 8. 验证删除成功")

		t.Log("🎉 完整工作流测试全部通过（包含编辑功能）！")
	})
}
