// Package models 数据模型定义
package models

import (
	"time"
)

// User 用户模型，对应 users 表
type User struct {
	ID            uint      `json:"id"`             // 用户ID
	Username      string    `json:"username"`       // 用户名
	PasswordHash  string    `json:"-"`              // bcrypt 密码哈希，JSON序列化时忽略
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名，必填
	Password string `json:"password" binding:"required"` // 密码，必填
}

// LoginResponse 登录响应体
type LoginResponse struct {
	Token    string `json:"token"`     // JWT 令牌
	UserID   uint   `json:"user_id"`   // 用户ID
	Username string `json:"username"`  // 用户名
}
