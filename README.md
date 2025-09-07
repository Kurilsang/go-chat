# Go Chat - 实时聊天应用

一个基于 Go 语言和 WebSocket 技术的实时聊天应用，支持用户注册登录、私聊功能和在线用户管理。

## 功能特性

### ✅ 已实现功能

- **用户认证系统**
  - 用户注册与登录
  - 密码加密存储（bcrypt）
  - 用户信息管理

- **实时通信**
  - WebSocket 连接管理
  - 私聊消息发送与接收
  - 在线用户实时更新
  - 连接状态监控
  - 心跳检测机制

- **数据存储**
  - MySQL 数据库（用户数据）
  - MongoDB 支持（消息存储）
  - Redis 缓存（会话管理）

- **Web 界面**
  - 测试聊天页面
  - 消息发送与显示
  - 在线用户列表

## 技术栈

- **后端**: Go 1.23+
- **Web 框架**: Gin
- **WebSocket**: Gorilla WebSocket
- **数据库**: MySQL + MongoDB + Redis
- **ORM**: GORM
- **前端**: HTML + JavaScript + CSS

## 快速开始

### 环境要求

- Go 1.23+
- MySQL 8.0+
- MongoDB 5.0+（可选）
- Redis 6.0+（可选）

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd go_chat
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置数据库**
   
   创建 MySQL 数据库：
   ```sql
   CREATE DATABASE go_chat;
   ```

4. **配置文件**
   
   编辑 `config/config.ini`：
   ```ini
   [service]
   AppMode = debug
   HttpPort = 8080

   [mysql]
   DbHost = localhost
   DbPort = 3306
   DbUser = root
   DbPassWord = your_password
   DbName = go_chat

   [MongoDB]
   MongoDBName = go_chat_mongo
   MongoDBAddr = localhost
   MongoDBPort = 27017

   [redis]
   RedisAddr = localhost:6379
   RedisDbName = 0
   ```

5. **运行应用**
   ```bash
   go run main.go
   ```

### 访问应用

- **API 服务**: http://localhost:8080
- **测试页面**: http://localhost:8080
- **WebSocket**: ws://localhost:8080/ws

## API 接口

### 认证接口

- **用户注册**
  ```
  POST /api/v1/auth/register
  Content-Type: application/json
  
  {
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456",
    "phone": "13800138000"
  }
  ```

- **用户登录**
  ```
  POST /api/v1/auth/login
  Content-Type: application/json
  
  {
    "username": "testuser",
    "password": "123456"
  }
  ```

- **获取用户信息**
  ```
  GET /api/v1/user/profile?id=1
  ```

### WebSocket 接口

- **连接地址**: `/ws?user_id=1&username=testuser`
- **消息格式**:
  ```json
  {
    "type": "private",
    "to_user_id": 2,
    "content": "Hello World",
    "timestamp": "2024-01-01T00:00:00Z"
  }
  ```

## 使用说明

### 测试聊天功能

1. 访问 http://localhost:8080
2. 输入用户 ID 和用户名
3. 点击"连接"建立 WebSocket 连接
4. 在消息发送区域输入目标用户 ID 和消息内容
5. 点击"发送"即可发送私聊消息

### 多用户测试

1. 打开多个浏览器标签页
2. 使用不同的用户 ID 连接
3. 在用户之间发送消息进行测试

## 项目结构

```
go_chat/
├── config/          # 配置管理
│   ├── config.go    # 配置加载
│   ├── config.ini   # 配置文件
│   ├── mysql.go     # MySQL配置
│   ├── mongodb.go   # MongoDB配置
│   └── redis.go     # Redis配置
├── controller/      # 控制器层
│   └── auth_controller.go
├── global/          # 全局变量
│   └── global.go
├── model/           # 数据模型
│   ├── init.go      # 数据库初始化
│   └── user.go      # 用户模型
├── router/          # 路由管理
│   └── router.go
├── service/         # 业务逻辑层
│   └── auth_service.go
├── websocket/       # WebSocket相关
│   ├── client.go    # 客户端连接
│   ├── handler.go   # WebSocket处理器
│   ├── hub.go       # 连接池管理
│   └── message.go   # 消息结构
├── static/          # 静态文件
│   ├── index.html   # 测试页面
│   ├── css/         # 样式文件
│   └── js/          # JavaScript文件
└── main.go          # 程序入口
```

## 开发说明

### 核心组件

- **Hub**: WebSocket 连接池管理器，负责客户端注册、消息分发
- **Client**: 客户端连接封装，处理消息收发和连接状态
- **Message**: 消息结构定义，支持多种消息类型
- **Handler**: WebSocket 请求处理器，负责连接升级和参数验证

### 消息类型

- `private`: 私聊消息
- `heartbeat`: 心跳检测
- `typing`: 正在输入
- `read`: 消息已读
- `join`: 用户加入
- `leave`: 用户离开
- `user_list`: 在线用户列表

## 注意事项

- 确保 MySQL 数据库已创建且配置正确
- WebSocket 连接需要提供有效的 user_id 和 username 参数
- 建议在生产环境中添加 JWT 认证中间件
- MongoDB 和 Redis 为可选组件，项目可在仅 MySQL 环境下运行
