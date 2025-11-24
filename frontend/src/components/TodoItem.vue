<template>
  <div class="todo-item-wrapper">
    <el-card class="todo-item" :class="{ completed: todo.completed }" shadow="hover">
    <div class="todo-content">
      <!-- 左侧：完成状态复选框 -->
      <div class="todo-checkbox">
        <el-checkbox
          :model-value="todo.completed"
          size="large"
          :disabled="statusLoading"
          @change="handleStatusChange"
        />
      </div>

      <!-- 中间：待办信息 -->
      <div class="todo-info">
        <div class="todo-header">
          <h3 class="todo-title">{{ todo.title }}</h3>
          <div class="todo-badges">
            <!-- 分类标签 -->
            <el-tag :type="getCategoryType(todo.category)" size="small" effect="plain">
              <el-icon><component :is="getCategoryIcon(todo.category)" /></el-icon>
              <span>{{ getCategoryLabel(todo.category) }}</span>
            </el-tag>
            <!-- 优先级标签 -->
            <el-tag v-if="todo.priority > 0" :type="getPriorityType(todo.priority)" size="small">
              优先级 {{ todo.priority }}
            </el-tag>
          </div>
        </div>

        <p v-if="todo.description" class="todo-description">
          {{ todo.description }}
        </p>

        <div class="todo-meta">
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            创建于 {{ formatDate(todo.created_at) }}
          </span>
          <span v-if="todo.completed" class="meta-item completed-text">
            <el-icon><CircleCheck /></el-icon>
            已完成
          </span>
        </div>
      </div>

      <!-- 右侧：操作按钮 -->
      <div class="todo-actions">
        <el-button
          type="primary"
          :icon="Edit"
          circle
          size="small"
          title="编辑"
          @click="handleEdit"
        />
        <el-button
          type="danger"
          :icon="Delete"
          circle
          size="small"
          title="删除"
          @click="handleDelete"
        />
      </div>
    </div>
  </el-card>

  <!-- 编辑对话框 -->
  <el-dialog
    v-model="editDialogVisible"
    title="编辑待办事项"
    width="500px"
    :close-on-click-modal="false"
  >
    <el-form
      ref="editFormRef"
      :model="editForm"
      :rules="editRules"
      label-width="80px"
    >
      <el-form-item label="标题" prop="title">
        <el-input
          v-model="editForm.title"
          placeholder="请输入标题"
          maxlength="255"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="editForm.description"
          type="textarea"
          placeholder="请输入描述（可选）"
          :rows="3"
          maxlength="1000"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="分类" prop="category">
        <el-select v-model="editForm.category" style="width: 100%">
          <el-option label="工作" value="work">
            <el-icon><Briefcase /></el-icon>
            <span style="margin-left: 8px">工作</span>
          </el-option>
          <el-option label="学习" value="study">
            <el-icon><Reading /></el-icon>
            <span style="margin-left: 8px">学习</span>
          </el-option>
          <el-option label="生活" value="life">
            <el-icon><Coffee /></el-icon>
            <span style="margin-left: 8px">生活</span>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="优先级" prop="priority">
        <el-rate v-model="editForm.priority" :max="5" show-score score-template="{value} 级" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="editDialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="editLoading" @click="handleEditSubmit">
        保存
      </el-button>
    </template>
  </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete, Clock, CircleCheck, Briefcase, Reading, Coffee } from '@element-plus/icons-vue'
import { updateTodoStatus, updateTodo, deleteTodo } from '../api/todo'

// 定义 props
const props = defineProps({
  todo: {
    type: Object,
    required: true,
  },
})

// 定义事件
const emit = defineEmits(['update', 'delete'])

// 状态管理
const statusLoading = ref(false)
const editDialogVisible = ref(false)
const editLoading = ref(false)
const editFormRef = ref(null)

// 编辑表单
const editForm = reactive({
  title: '',
  description: '',
  category: '',
  priority: 0,
})

// 编辑表单验证规则
const editRules = {
  title: [
    { required: true, message: '标题不能为空', trigger: 'blur' },
    { min: 1, max: 255, message: '标题长度需在 1 到 255 个字符之间', trigger: 'blur' },
  ],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  priority: [
    { required: true, message: '请选择优先级', trigger: 'change' },
    { type: 'number', min: 0, max: 5, message: '优先级范围为 0-5', trigger: 'change' },
  ],
}

// 切换完成状态
const handleStatusChange = async (value) => {
  try {
    statusLoading.value = true

    // 确保 version 字段存在，默认为 0
    const version = props.todo.version !== undefined ? props.todo.version : 0

    const response = await updateTodoStatus(props.todo.id, {
      completed: value,
      version: version,
    })

    ElMessage.success(value ? '已标记为完成' : '已标记为未完成')
    
    // 刷新列表获取最新数据（包含更新后的 version）
    emit('update')
  } catch (error) {
    // 处理版本冲突
    if (error.status === 409) {
      const conflictData = error.data
      ElMessageBox.confirm(
        `该待办事项已被其他设备修改（当前版本：${conflictData.current_version}）。是否刷新最新数据？`,
        '数据冲突',
        {
          confirmButtonText: '刷新数据',
          cancelButtonText: '取消',
          type: 'warning',
        }
      )
        .then(() => {
          emit('update')
        })
        .catch(() => {
          // 用户取消
        })
    } else {
      // 其他错误，刷新列表获取最新数据
      emit('update')
    }
  } finally {
    statusLoading.value = false
  }
}

// 打开编辑对话框
const handleEdit = () => {
  editForm.title = props.todo.title
  editForm.description = props.todo.description || ''
  editForm.category = props.todo.category
  editForm.priority = props.todo.priority
  editDialogVisible.value = true
}

// 提交编辑
const handleEditSubmit = async () => {
  if (!editFormRef.value) return

  try {
    await editFormRef.value.validate()
    editLoading.value = true

    // 确保 version 字段存在，默认为 0
    const version = props.todo.version !== undefined ? props.todo.version : 0

    await updateTodo(props.todo.id, {
      title: editForm.title,
      description: editForm.description,
      category: editForm.category,
      priority: editForm.priority,
      version: version,
    })

    ElMessage.success('修改成功')
    editDialogVisible.value = false
    emit('update')
  } catch (error) {
    // 处理版本冲突
    if (error.status === 409) {
      const conflictData = error.data
      ElMessageBox.confirm(
        `该待办事项已被其他设备修改（当前版本：${conflictData.current_version}）。是否刷新最新数据？`,
        '数据冲突',
        {
          confirmButtonText: '刷新数据',
          cancelButtonText: '取消',
          type: 'warning',
        }
      )
        .then(() => {
          editDialogVisible.value = false
          emit('update')
        })
        .catch(() => {
          // 用户取消
        })
    } else {
      // 其他错误，刷新列表
      emit('update')
    }
  } finally {
    editLoading.value = false
  }
}

// 删除待办
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除这个待办事项吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await deleteTodo(props.todo.id)
    ElMessage.success('删除成功')
    emit('delete', props.todo.id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 辅助函数：获取分类类型
const getCategoryType = (category) => {
  const typeMap = {
    work: 'danger',
    study: 'warning',
    life: 'success',
  }
  return typeMap[category] || ''
}

// 辅助函数：获取分类图标
const getCategoryIcon = (category) => {
  const iconMap = {
    work: 'Briefcase',
    study: 'Reading',
    life: 'Coffee',
  }
  return iconMap[category] || 'Document'
}

// 辅助函数：获取分类标签
const getCategoryLabel = (category) => {
  const labelMap = {
    work: '工作',
    study: '学习',
    life: '生活',
  }
  return labelMap[category] || category
}

// 辅助函数：获取优先级类型
const getPriorityType = (priority) => {
  if (priority >= 4) return 'danger'
  if (priority >= 2) return 'warning'
  return 'info'
}

// 辅助函数：格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date

  // 小于 1 分钟
  if (diff < 60 * 1000) {
    return '刚刚'
  }
  // 小于 1 小时
  if (diff < 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 1000))} 分钟前`
  }
  // 小于 1 天
  if (diff < 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 60 * 1000))} 小时前`
  }
  // 小于 7 天
  if (diff < 7 * 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (24 * 60 * 60 * 1000))} 天前`
  }

  // 显示完整日期
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
.todo-item {
  margin-bottom: 12px;
  transition: all 0.3s;
}

.todo-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.todo-item.completed {
  opacity: 0.7;
  background-color: #f5f7fa;
}

.todo-content {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.todo-checkbox {
  flex-shrink: 0;
  padding-top: 4px;
}

.todo-info {
  flex: 1;
  min-width: 0;
}

.todo-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.todo-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  word-break: break-word;
}

.todo-item.completed .todo-title {
  text-decoration: line-through;
  color: #909399;
}

.todo-badges {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.todo-badges .el-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.todo-description {
  margin: 0 0 8px 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
}

.todo-item.completed .todo-description {
  color: #909399;
}

.todo-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 12px;
  color: #909399;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.completed-text {
  color: #67c23a;
  font-weight: 500;
}

.todo-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .todo-content {
    flex-direction: column;
  }

  .todo-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .todo-actions {
    align-self: flex-end;
  }
}
</style>

