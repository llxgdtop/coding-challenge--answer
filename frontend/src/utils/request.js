import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api', // 基础 URL，会自动添加到所有请求前
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 可以在这里添加 token 等认证信息
    // const token = localStorage.getItem('token')
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`
    // }
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 将英文错误信息转换为中文
const translateErrorMessage = (message) => {
  const errorMap = {
    'title is required and cannot be empty': '标题不能为空',
    'title cannot exceed 255 characters': '标题长度不能超过 255 个字符',
    'invalid category': '分类无效',
    'priority must be between 0 and 5': '优先级必须在 0 到 5 之间',
    'invalid id': 'ID 无效',
    'todo not found': '待办事项不存在',
    'version conflict': '数据已被其他设备修改',
    'Invalid input': '输入内容有误',
    'required': '必填项未填写',
  }

  // 查找匹配的错误信息
  for (const [en, zh] of Object.entries(errorMap)) {
    if (message.toLowerCase().includes(en.toLowerCase())) {
      return zh
    }
  }

  return message // 如果没有匹配，返回原始信息
}

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data

    // 根据后端约定的 code 判断请求是否成功
    // code 为 0 表示成功
    if (res.code === 0) {
      return res
    } else {
      // 业务错误 - 转换为中文提示
      const errorMsg = translateErrorMessage(res.message || '请求失败')
      ElMessage.error(errorMsg)
      return Promise.reject(new Error(errorMsg))
    }
  },
  (error) => {
    console.error('响应错误:', error)

    // 处理 HTTP 错误状态码
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 400:
          ElMessage.error(translateErrorMessage(data.message) || '请求参数错误')
          break
        case 404:
          ElMessage.error(translateErrorMessage(data.message) || '请求的资源不存在')
          break
        case 409:
          // 版本冲突（乐观锁）- 不在这里提示，交给业务层处理
          return Promise.reject(error.response)
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(translateErrorMessage(data.message) || `请求失败 (${status})`)
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      // 其他错误
      ElMessage.error(error.message || '请求失败')
    }

    return Promise.reject(error)
  }
)

export default request

