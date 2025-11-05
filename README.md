个人博客系统后端
基于 Go + Gin + GORM 开发的个人博客系统后端，提供完整的文章管理、用户认证和评论功能。

📋 项目特性
✅ 用户注册与登录（JWT认证）

✅ 博客文章CRUD操作

✅ 文章评论功能

✅ Swagger API文档

✅ 密码加密存储

✅ 权限控制（用户只能操作自己的资源）

✅ 完整的错误处理和日志记录

🛠 技术栈
编程语言: Go 1.21+

Web框架: Gin

ORM: GORM

数据库: MySQL 8.0+

认证: JWT

文档: Swagger/OpenAPI 3.0

密码加密: bcrypt

📦 项目结构
text
GoBlogService/
├── main.go                # 应用入口
├── config/
│   ├── config.go          # 配置管理
│   └── config.yaml        # 配置文件
├── models/
│   ├── user.go            # 用户模型
│   ├── post.go            # 文章模型
│   └── comment.go         # 评论模型
├── controllers/
│   ├── auth.go            # 认证控制器
│   ├── post.go            # 文章控制器
│   └── comment.go         # 评论控制器
├── middleware/
│   ├── auth.go            # JWT认证中间件
│   └── logger.go          # 日志中间件
├── utils/
│   ├── jwt.go             # JWT工具
│   └── password.go        # 密码工具
├── database/
│   ├── database.go        # 数据库连接
│   └── seeds.go           # 样本数据
└── docs/                  # Swagger文档
🚀 快速开始
环境要求
Go 1.21 或更高版本

MySQL 8.0 或更高版本

Git

安装步骤
克隆项目

bash
git clone <repository-url>
cd blog-backend
安装依赖

bash
go mod tidy
安装 Swag 工具

bash
go install github.com/swaggo/swag/cmd/swag@latest
生成 Swagger 文档

bash
swag init
数据库配置

创建 MySQL 数据库：

sql
CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
环境变量配置

创建 .env 文件（可选）：

env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=blog
JWT_SECRET=your-super-secret-jwt-key
SERVER_PORT=9090
或者直接修改 config/config.go 中的默认值。

启动应用
开发模式

bash
go run main.go
构建并运行

bash
go build -o blog-backend
./blog-backend
使用 Air 热重载（开发推荐）

bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 运行
air
应用启动后，访问以下地址：

应用API: http://localhost:9090

Swagger文档: http://localhost:9090/swagger/index.html

健康检查: http://localhost:9090/health

📚 API 文档
认证接口
1. 用户注册
URL: POST /auth/register

Body:

json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
响应:

json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
2. 用户登录
URL: POST /auth/login

Body:

json
{
  "username": "testuser",
  "password": "password123"
}
响应:

json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
文章接口
3. 获取文章列表
URL: GET /posts

认证: 不需要

响应:

json
{
  "posts": [
    {
      "id": 1,
      "title": "第一篇博客",
      "content": "这是我的第一篇博客内容...",
      "user_id": 1,
      "user": {
        "id": 1,
        "username": "admin"
      },
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ]
}
4. 创建文章
URL: POST /posts

认证: Bearer Token 需要

Headers: Authorization: Bearer <jwt-token>

Body:

json
{
  "title": "新文章标题",
  "content": "文章内容..."
}
5. 更新文章
URL: PUT /posts/:id

认证: Bearer Token 需要（仅作者可操作）

6. 删除文章
URL: DELETE /posts/:id

认证: Bearer Token 需要（仅作者可操作）

评论接口
7. 创建评论
URL: POST /comments

认证: Bearer Token 需要

Body:

json
{
  "content": "这是一条评论",
  "post_id": 1
}
8. 获取文章评论
URL: GET /posts/:id/comments

认证: 不需要

🔧 Postman 测试指南
导入 Postman 集合
下载 Postman Collection JSON 文件

打开 Postman → Import → 选择文件

导入后设置环境变量：

base_url: http://localhost:9090

token: (登录后自动设置)

测试流程
步骤 1: 用户注册
执行 "用户注册" 请求

检查响应状态码是否为 201

步骤 2: 用户登录
执行 "用户登录" 请求

复制响应中的 token 值

在环境变量中设置 token 变量

步骤 3: 创建文章
执行 "创建文章" 请求

检查响应状态码是否为 201

记录返回的文章 ID

步骤 4: 获取文章列表
执行 "获取文章列表" 请求

验证返回的文章数据

步骤 5: 添加评论
执行 "创建评论" 请求

使用步骤3中获取的文章ID

步骤 6: 测试权限控制
注册第二个用户并获取token

尝试修改/删除第一个用户的文章

应该返回 403 状态码

测试用例示例
bash
# 1. 注册用户
curl -X POST http://localhost:9090/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123",
    "email": "test@example.com"
  }'

# 2. 用户登录
curl -X POST http://localhost:9090/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'

# 3. 创建文章 (使用上一步获取的token)
curl -X POST http://localhost:9090/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "测试文章",
    "content": "这是测试文章的内容"
  }'

# 4. 获取文章列表
curl -X GET http://localhost:9090/posts

# 5. 添加评论
curl -X POST http://localhost:9090/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "content": "这是一条测试评论",
    "post_id": 1
  }'
🗄 数据库配置
表结构
users 表

sql
CREATE TABLE users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
posts 表

sql
CREATE TABLE posts (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(200) NOT NULL,
  content TEXT NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
comments 表

sql
CREATE TABLE comments (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  content TEXT NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  post_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
🔒 安全特性
JWT Token 认证

密码 bcrypt 加密

SQL 注入防护（GORM 参数化查询）

CORS 中间件支持

请求频率限制（可扩展）

🐛 故障排除
常见问题
数据库连接失败

检查 MySQL 服务是否运行

验证数据库配置信息

确认数据库用户权限

Swagger 文档无法访问

运行 swag init 重新生成文档

检查 docs 文件夹是否存在

JWT 认证失败

检查 Token 是否过期（默认24小时）

验证 JWT Secret 配置

端口占用

修改 SERVER_PORT 环境变量

杀死占用端口的进程

日志查看
应用运行日志会输出到控制台，包含：

请求方法、路径、状态码、响应时间

错误信息和堆栈跟踪

数据库查询日志（开发环境）

📝 开发指南
添加新的 API 端点
在对应的控制器中添加处理方法

添加 Swagger 注释

在 main.go 中注册路由

运行 swag init 更新文档

数据库迁移
bash
# 自动迁移（开发环境）
go run main.go --migrate

# 或直接在代码中调用
database.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
📄 许可证
MIT License

🤝 贡献
欢迎提交 Issue 和 Pull Request！