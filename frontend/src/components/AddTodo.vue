<template>
  <el-card class="add-todo-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <el-icon><CirclePlus /></el-icon>
        <span>添加新待办</span>
      </div>
    </template>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="80px"
      @submit.prevent="handleSubmit"
    >
      <el-form-item label="标题" prop="title">
        <el-input
          v-model="form.title"
          placeholder="请输入待办事项标题"
          maxlength="255"
          show-word-limit
          clearable
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          placeholder="请输入详细描述（可选）"
          :rows="3"
          maxlength="1000"
          show-word-limit
          clearable
        />
      </el-form-item>

      <el-form-item label="分类" prop="category">
        <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
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
        <el-rate
          v-model="form.priority"
          :max="5"
          show-score
          score-template="{value} 级"
          style="height: 32px; line-height: 32px"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          <el-icon><CirclePlus /></el-icon>
          <span>添加</span>
        </el-button>
        <el-button @click="handleReset">
          <el-icon><RefreshLeft /></el-icon>
          <span>重置</span>
        </el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { addTodo } from '../api/todo'

// 表单引用
const formRef = ref(null)
const loading = ref(false)

// 表单数据
const form = reactive({
  title: '',
  description: '',
  category: 'life', // 默认分类
  priority: 0, // 默认优先级
})

// 表单验证规则
const rules = {
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

// 定义事件
const emit = defineEmits(['success'])

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    // 验证表单
    await formRef.value.validate()

    loading.value = true

    // 调用 API
    const response = await addTodo({
      title: form.title,
      description: form.description || '', // 空字符串作为默认值
      category: form.category,
      priority: form.priority,
    })

    ElMessage.success('添加成功！')

    // 重置表单
    handleReset()

    // 通知父组件刷新列表
    emit('success', response.data)
  } catch (error) {
    console.error('添加失败:', error)
    // 错误提示已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

// 重置表单
const handleReset = () => {
  if (!formRef.value) return
  formRef.value.resetFields()
  form.priority = 0 // 手动重置优先级
}
</script>

<style scoped>
.add-todo-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.el-button {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>

