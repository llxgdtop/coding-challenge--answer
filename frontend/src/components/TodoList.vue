<template>
  <div class="todo-list-container">
    <!-- 筛选和排序栏 -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <!-- 分类筛选 -->
          <div class="filter-group">
            <span class="filter-label">分类：</span>
            <el-radio-group v-model="filters.category" size="small" @change="handleFilterChange">
              <el-radio-button label="">全部</el-radio-button>
              <el-radio-button label="work">
                <el-icon><Briefcase /></el-icon>
                <span>工作</span>
              </el-radio-button>
              <el-radio-button label="study">
                <el-icon><Reading /></el-icon>
                <span>学习</span>
              </el-radio-button>
              <el-radio-button label="life">
                <el-icon><Coffee /></el-icon>
                <span>生活</span>
              </el-radio-button>
            </el-radio-group>
          </div>

          <!-- 排序方式 -->
          <div class="filter-group">
            <span class="filter-label">排序：</span>
            <el-select v-model="filters.sort" size="small" style="width: 140px" @change="handleFilterChange">
              <el-option label="创建时间" value="created_at" />
              <el-option label="优先级" value="priority" />
            </el-select>
          </div>
        </div>

        <!-- 刷新按钮 -->
        <div class="filter-right">
          <el-button
            :icon="Refresh"
            circle
            size="small"
            :loading="loading"
            @click="fetchTodos"
          />
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="statistics">
        <el-statistic title="总任务" :value="statistics.total">
          <template #prefix>
            <el-icon><List /></el-icon>
          </template>
        </el-statistic>
        <el-divider direction="vertical" />
        <el-statistic title="未完成" :value="statistics.pending">
          <template #prefix>
            <el-icon><Clock /></el-icon>
          </template>
        </el-statistic>
        <el-divider direction="vertical" />
        <el-statistic title="已完成" :value="statistics.completed">
          <template #prefix>
            <el-icon><CircleCheck /></el-icon>
          </template>
        </el-statistic>
      </div>
    </el-card>

    <!-- 待办列表 -->
    <div v-loading="loading" class="todos-wrapper">
      <!-- 空状态 -->
      <el-empty
        v-if="!loading && todos.length === 0"
        description="暂无待办事项，添加一个吧！"
        :image-size="120"
      >
        <template #image>
          <el-icon :size="120" color="#909399"><Document /></el-icon>
        </template>
      </el-empty>

      <!-- 待办列表 -->
      <div v-else class="todos-list">
        <TransitionGroup name="list">
          <TodoItem
            v-for="todo in todos"
            :key="todo.id"
            :todo="todo"
            @update="fetchTodos"
            @delete="handleTodoDelete"
          />
        </TransitionGroup>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Briefcase,
  Reading,
  Coffee,
  Refresh,
  List,
  Clock,
  CircleCheck,
  Document,
} from '@element-plus/icons-vue'
import TodoItem from './TodoItem.vue'
import { getTodos } from '../api/todo'

// 状态管理
const loading = ref(false)
const todos = ref([])
const filters = reactive({
  category: '', // 分类筛选
  sort: 'created_at', // 排序方式
})

// 自动刷新定时器
let refreshTimer = null

// 计算统计信息
const statistics = computed(() => {
  const total = todos.value.length
  const completed = todos.value.filter((t) => t.completed).length
  const pending = total - completed
  const completionRate = total > 0 ? Math.round((completed / total) * 100) : 0

  return {
    total,
    completed,
    pending,
    completionRate,
  }
})

// 获取待办列表
const fetchTodos = async () => {
  try {
    loading.value = true

    const params = {}
    if (filters.category) {
      params.category = filters.category
    }
    if (filters.sort) {
      params.sort = filters.sort
    }

    const response = await getTodos(params)
    todos.value = response.data || []
  } catch (error) {
    console.error('获取待办列表失败:', error)
    ElMessage.error('获取待办列表失败')
  } finally {
    loading.value = false
  }
}

// 筛选或排序变化
const handleFilterChange = () => {
  fetchTodos()
}

// 删除待办
const handleTodoDelete = (id) => {
  todos.value = todos.value.filter((t) => t.id !== id)
}

// 启动自动刷新（每 30 秒）
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    fetchTodos()
  }, 30000) // 30 秒
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchTodos()
  startAutoRefresh()
})

// 组件卸载时清理定时器
onUnmounted(() => {
  stopAutoRefresh()
})

// 暴露方法供父组件调用
defineExpose({
  fetchTodos,
})
</script>

<style scoped>
.todo-list-container {
  width: 100%;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
}

.filter-left {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 24px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.filter-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.statistics {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 20px 16px;
  background-color: #f5f7fa;
  border-radius: 4px;
  min-height: 80px;
  overflow: visible;
}

.statistics .el-statistic {
  flex-shrink: 0;
  min-width: 100px;
}

.statistics .el-statistic :deep(.el-statistic__head) {
  font-size: 14px;
  margin-bottom: 4px;
}

.statistics .el-statistic :deep(.el-statistic__content) {
  font-size: 24px;
  line-height: 1.5;
}

.statistics .el-divider {
  height: 50px;
  margin: 0;
}

.todos-wrapper {
  min-height: 200px;
}

.todos-list {
  width: 100%;
}

/* 列表过渡动画 */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-move {
  transition: transform 0.3s ease;
}

/* 响应式 */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: flex-start;
  }

  .filter-left {
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
  }

  .filter-right {
    align-self: flex-end;
  }

  .statistics {
    flex-wrap: wrap;
    justify-content: center;
  }

}
</style>

