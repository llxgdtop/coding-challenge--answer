import request from '../utils/request'

/**
 * 获取待办事项列表
 * @param {Object} params - 查询参数
 * @param {string} params.category - 分类筛选 (work/study/life/all)
 * @param {string} params.sort - 排序方式 (priority/created_at)
 */
export function getTodos(params) {
  return request({
    url: '/todos',
    method: 'get',
    params,
  })
}

/**
 * 根据 ID 获取单个待办事项
 * @param {number} id - 待办事项 ID
 */
export function getTodoById(id) {
  return request({
    url: `/todos/${id}`,
    method: 'get',
  })
}

/**
 * 添加待办事项
 * @param {Object} data - 待办事项数据
 * @param {string} data.title - 标题（必填）
 * @param {string} data.description - 描述（可选）
 * @param {string} data.category - 分类（work/study/life，可选，默认 life）
 * @param {number} data.priority - 优先级（0-5，可选，默认 0）
 */
export function addTodo(data) {
  return request({
    url: '/todos',
    method: 'post',
    data,
  })
}

/**
 * 更新待办事项（编辑）
 * @param {number} id - 待办事项 ID
 * @param {Object} data - 更新数据
 * @param {string} data.title - 标题
 * @param {string} data.description - 描述
 * @param {string} data.category - 分类
 * @param {number} data.priority - 优先级
 * @param {number} data.version - 版本号（乐观锁）
 */
export function updateTodo(id, data) {
  return request({
    url: `/todos/${id}`,
    method: 'put',
    data,
  })
}

/**
 * 更新待办事项状态（完成/未完成）
 * @param {number} id - 待办事项 ID
 * @param {Object} data - 更新数据
 * @param {boolean} data.completed - 是否完成
 * @param {number} data.version - 版本号（乐观锁）
 */
export function updateTodoStatus(id, data) {
  return request({
    url: `/todos/${id}/status`,
    method: 'put',
    data,
  })
}

/**
 * 删除待办事项
 * @param {number} id - 待办事项 ID
 */
export function deleteTodo(id) {
  return request({
    url: `/todos/${id}`,
    method: 'delete',
  })
}

