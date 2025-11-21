## TODO List 项目说明文档

### 1. 技术选型

- 编程语言：Go+html+css+JavaScript，理由：对Go较为熟悉，开发较为快速。

- 框架：后端用Gin框架，理由：基于前缀树进行路由匹配，相较于别的框架如Beego采用正则表达式的要快的多。前端用Vue3，理由：生态丰富，更加常见。
- 数据库：TiDB，理由：随着业务增长，只需增加节点即可平滑扩展成分布式集群、兼容MySQL 5.7协议，现有的 MySQL 工具链基本都能直接使用，避免后期更换麻烦
- 替代方案对比：在当前这种非常简单的场景下确实应该用mysql，但如果不考虑服务器的性能的话，采用TiDB可以避免后期重构分库分表的成本。相当于一步到位。



### 2.项目结构设计

采用前后端分离的架构，初步设定目录结构与职责如下：

```
coding-challenge--answer/
├── backend/                 # 后端服务
│   ├── main.go             # 入口文件
│   ├── config/             # 配置文件
│   │   └── config.go       # 数据库配置
│   ├── models/             # 数据模型
│   │   └── todo.go         # TODO模型定义
│   ├── controllers/        # 控制器
│   │   └── todo_controller.go
│   ├── services/           # 业务逻辑
│   │   └── todo_service.go
│   ├── middleware/         # 中间件
│   │   └── cors.go         # CORS处理
│   ├── router/             # 路由
│   │   └── router.go
│   └── go.mod              # Go依赖管理
├── frontend/               # 前端应用
│   ├── src/
│   │   ├── main.js         # 入口文件
│   │   ├── App.vue         # 根组件
│   │   ├── components/     # 组件目录
│   │   │   ├── TodoList.vue    # TODO列表组件
│   │   │   ├── TodoItem.vue    # TODO项组件
│   │   │   └── AddTodo.vue     # 添加TODO组件
│   │   ├── api/            # API请求
│   │   │   └── todo.js
│   │   └── utils/          # 工具函数
│   │       └── request.js  # Axios封装
│   ├── package.json
│   └── vite.config.js
```

sql建表语句

```sql
CREATE TABLE todos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category ENUM('work', 'study', 'life') DEFAULT 'life',
    priority INT DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    version INT DEFAULT 0,  -- 用于乐观锁
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 索引优化
CREATE INDEX idx_category ON todos(category);
CREATE INDEX idx_priority ON todos(priority);
CREATE INDEX idx_completed ON todos(completed);
```





