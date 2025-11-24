<template>
  <div id="app">
    <!-- 头部 -->
    <header class="app-header">
      <div class="container">
        <div class="header-content">
          <div class="logo">
            <el-icon :size="32" color="#409eff"><Calendar /></el-icon>
            <h1>Todo List</h1>
          </div>
          <div class="header-info">
            <el-tag type="info" size="large">
              <el-icon><TrendCharts /></el-icon>
              <span>基于 Vue3 + Element-Plus + Gin + TiDB</span>
            </el-tag>
          </div>
        </div>
      </div>
    </header>

    <!-- 主体内容 -->
    <main class="app-main">
      <div class="container">
        <el-row :gutter="20">
          <!-- 左侧：添加待办 -->
          <el-col :xs="24" :sm="24" :md="10" :lg="8">
            <AddTodo @success="handleAddSuccess" />
          </el-col>

          <!-- 右侧：待办列表 -->
          <el-col :xs="24" :sm="24" :md="14" :lg="16">
            <TodoList ref="todoListRef" />
          </el-col>
        </el-row>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="app-footer">
      <div class="container">
        <p>© 2025 Todo List Application | Powered by Vue3 + Gin + TiDB</p>
        <p class="footer-tips">
          <el-icon><InfoFilled /></el-icon>
          <span>支持多设备协作</span>
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Calendar, TrendCharts, InfoFilled } from '@element-plus/icons-vue'
import AddTodo from './components/AddTodo.vue'
import TodoList from './components/TodoList.vue'

// TodoList 组件引用
const todoListRef = ref(null)

// 添加成功后刷新列表
const handleAddSuccess = () => {
  if (todoListRef.value) {
    todoListRef.value.fetchTodos()
  }
}
</script>

<style scoped>
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 头部 */
.app-header {
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-info .el-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  font-size: 13px;
}

/* 主体 */
.app-main {
  flex: 1;
  padding: 40px 0;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 页脚 */
.app-footer {
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 24px 0;
  text-align: center;
  color: #606266;
  font-size: 14px;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.1);
}

.app-footer p {
  margin: 8px 0;
}

.footer-tips {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 13px;
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .logo h1 {
    font-size: 24px;
  }

  .app-main {
    padding: 20px 0;
  }

  .el-col {
    margin-bottom: 20px;
  }
}
</style>
