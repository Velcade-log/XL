# 问卷小系统

一个简单的问卷管理系统，支持发布问卷、填写问卷、查看答卷等功能。

## 技术栈

- **后端**: Go + Gin 框架
- **数据库**: PostgreSQL
- **前端**: 原生 HTML + CSS + JavaScript（无框架）
- **鉴权**: JWT (JSON Web Token)
- **密码加密**: bcrypt

---

## 一、项目目录结构

```
survey-system/
├── backend/                     # Go 后端项目
│   ├── main.go                  # 主程序入口
│   ├── go.mod                   # Go 模块依赖
│   ├── config/
│   │   └── config.go            # 配置管理（环境变量）
│   ├── database/
│   │   └── database.go          # 数据库连接初始化
│   ├── middleware/
│   │   ├── auth.go              # JWT 鉴权中间件
│   │   └── cors.go              # CORS 跨域中间件
│   ├── models/
│   │   ├── user.go              # 用户模型
│   │   ├── survey.go            # 问卷模型
│   │   └── answer.go            # 答卷模型
│   ├── handlers/
│   │   ├── auth_handler.go      # 登录接口处理器
│   │   ├── survey_handler.go    # 问卷接口处理器
│   │   └── answer_handler.go    # 答卷接口处理器
│   └── utils/
│       └── jwt.go               # JWT 工具函数
├── frontend/                    # 前端页面（原生 HTML）
│   ├── login.html               # 登录页面
│   ├── index.html               # 问卷列表页
│   ├── survey.html              # 问卷填写页
│   ├── answer-list.html         # 答卷列表页
│   └── css/
│       └── style.css            # 全局样式
└── sql/
    └── init.sql                 # 数据库初始化脚本
```

---

## 二、数据库设计

### 2.1 建表 SQL

完整建表脚本位于 `sql/init.sql`，包含三张表：

| 表名 | 说明 |
|------|------|
| `users` | 用户表，存储账号和 bcrypt 密码哈希 |
| `surveys` | 问卷表，题目以 JSONB 格式存储 |
| `survey_answers` | 答卷表，联合唯一约束保证一用户一问卷一份答卷 |

### 2.2 初始化用户（重要！）

系统**不提供注册功能**，只有预先在数据库中写入的 2 个账号可用。

**步骤一：生成 bcrypt 密码哈希**

你需要为两个用户的密码生成 bcrypt 哈希值。有以下几种方式：

**方式 A：使用 Go 代码生成（推荐）**

创建一个临时 Go 文件：

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "your_password_here"
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Println("生成失败:", err)
        return
    }
    fmt.Println(string(hash))
}
```

运行：`go run generate_hash.go`

**方式 B：使用在线工具**

访问 https://bcrypt-generator.com/ 在线生成。

**步骤二：执行插入 SQL**

将生成的哈希值填入以下 SQL 并执行：

```sql
INSERT INTO users (username, password_hash) VALUES
    ('user1', '<user1的bcrypt哈希值>'),
    ('user2', '<user2的bcrypt哈希值>');
```

---

## 三、环境变量说明

后端使用以下环境变量：

| 变量名 | 必需 | 说明 | 示例 |
|--------|------|------|------|
| `DATABASE_URL` | 是 | PostgreSQL 连接字符串 | `postgres://user:pass@localhost:5432/survey_db?sslmode=disable` |
| `JWT_SECRET` | 是 | JWT 签名密钥（生产环境请使用强随机字符串） | `your-strong-secret-key-here` |
| `PORT` | 否 | 服务监听端口，默认 `8080` | `8080` |

---

## 四、本地运行完整步骤

### 4.1 准备 PostgreSQL 数据库

1. 安装并启动 PostgreSQL（本地或使用 Docker）
2. 创建数据库，例如 `survey_db`
3. 执行 `sql/init.sql` 建表
4. 按上面的说明插入 2 个用户

### 4.2 启动 Go 后端

```bash
# 进入后端目录
cd survey-system/backend

# 下载依赖
go mod download

# 设置环境变量（Windows PowerShell）
$env:DATABASE_URL="postgres://user:password@localhost:5432/survey_db?sslmode=disable"
$env:JWT_SECRET="your-secret-key"

# 或者设置环境变量（Windows CMD）
# set DATABASE_URL=postgres://...
# set JWT_SECRET=your-secret-key

# 启动服务
go run main.go
```

服务启动后，API 地址为 `http://localhost:8080/api`

### 4.3 运行前端

前端是纯静态 HTML 文件，有两种运行方式：

**方式一：直接打开（最简单）**

直接用浏览器打开 `frontend/login.html` 即可。

**方式二：使用静态文件服务器（推荐）**

由于直接打开 file:// 协议可能导致某些浏览器限制，建议使用简单的 HTTP 服务器：

```bash
# 使用 Python 内置服务器（需安装 Python）
cd survey-system/frontend
python -m http.server 3000

# 或者使用 Node.js 的 http-server
npx http-server frontend -p 3000
```

然后访问 `http://localhost:3000/login.html`

### 4.4 配置前端 API 地址

如果后端地址不是 `http://localhost:8080`，需要修改每个 HTML 文件中的 `API_BASE` 常量：

```javascript
const API_BASE = 'http://你的后端地址/api';
```

涉及文件：`login.html`、`index.html`、`survey.html`、`answer-list.html`

---

## 五、API 接口说明

所有需要登录的接口需在请求头中携带：
```
Authorization: Bearer <JWT令牌>
```

### 5.1 登录接口

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/login` | 否 | 用户登录，返回 JWT |

请求体：
```json
{
    "username": "user1",
    "password": "password123"
}
```

响应：
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 1,
    "username": "user1"
}
```

### 5.2 用户信息

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/me` | 是 | 获取当前登录用户信息 |

### 5.3 问卷相关

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/surveys` | 是 | 获取全部问卷列表 |
| POST | `/api/surveys` | 是 | 创建新问卷 |
| GET | `/api/surveys/:id` | 是 | 获取单份问卷详情（含题目） |

**创建问卷请求体示例**：
```json
{
    "title": "用户满意度调查",
    "description": "请填写以下问卷",
    "questions": [
        {"id": "q1", "title": "您的姓名", "type": "text"},
        {"id": "q2", "title": "您的建议", "type": "text"}
    ]
}
```

### 5.4 答卷相关

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| PUT / POST | `/api/surveys/:id/answer` | 是 | 提交/更新当前用户答卷（Upsert） |
| GET | `/api/surveys/:id/my-answer` | 是 | 获取当前用户自己的答卷 |
| GET | `/api/surveys/:id/answers` | 是 | 获取问卷所有答卷（**仅创建者可访问，否则403**） |

**提交答卷请求体示例**：
```json
{
    "answer": {
        "q1": "张三",
        "q2": "服务很棒"
    }
}
```

---

## 六、业务权限说明

| 功能 | 谁可以操作 |
|------|-----------|
| 查看问卷列表 | 所有登录用户 |
| 发布问卷 | 所有登录用户 |
| 填写/修改问卷 | 所有登录用户（每人每问卷一份答卷，可反复修改） |
| 查看自己的答卷 | 所有登录用户 |
| 查看所有答卷列表 | **仅该问卷的发布者**（后端严格校验，前端仅控制按钮显示） |

> **注意**：答卷查看权限的真实校验在后端 `GetSurveyAnswers` 函数中完成，前端按钮只是 UI 控制，不能作为安全依据。

---

## 七、部署方案

### 7.1 数据库：Neon PostgreSQL

1. 注册 [Neon](https://neon.tech/) 账号
2. 创建一个新的 PostgreSQL 项目
3. 复制连接字符串（DATABASE_URL）
4. 使用任意 SQL 客户端连接，执行 `sql/init.sql` 建表和插入用户

### 7.2 后端：Render

1. 注册 [Render](https://render.com/) 账号
2. 选择 "New" → "Web Service"
3. 连接你的代码仓库
4. 配置：
   - **Runtime**: Go
   - **Build Command**: `cd backend && go build -o main .`
   - **Start Command**: `./backend/main`
   - **Root Directory**: 留空或设为项目根目录
5. 添加环境变量：
   - `DATABASE_URL`: Neon 提供的连接字符串
   - `JWT_SECRET`: 自定义强随机字符串
   - `PORT`: 10000（Render 默认端口）
6. 部署完成后，记下后端服务的 URL

### 7.3 前端：Vercel

1. 注册 [Vercel](https://vercel.com/) 账号
2. 选择 "Add New" → "Project"
3. 导入代码仓库
4. 配置：
   - **Framework Preset**: Other
   - **Build Command**: 留空
   - **Output Directory**: `frontend`
5. 部署前，先修改 4 个 HTML 文件中的 `API_BASE` 为 Render 上的后端地址
6. 部署完成后，通过 Vercel 提供的域名访问

---

## 八、安全注意事项

1. 生产环境务必使用强 `JWT_SECRET`（建议 32 位以上随机字符串）
2. 生产环境数据库连接务必启用 SSL（Neon 默认启用）
3. 前端部署时确保使用 HTTPS
4. 两个初始用户的密码请务必修改为强密码
5. CORS 目前设置为允许所有来源（`*`），生产环境建议限制为具体域名
