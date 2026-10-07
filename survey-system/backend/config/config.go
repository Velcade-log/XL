// Package config 配置管理
// 配置直接写在代码里，简单方便
package config

import (
	"os"
)

// Config 应用配置结构体
type Config struct {
	DatabaseURL string // PostgreSQL 连接字符串
	JWTSecret   string // JWT 签名密钥
	ServerPort  string // 服务监听端口
}

// AppConfig 全局配置实例
var AppConfig Config

// LoadConfig 加载配置
// 优先从环境变量读取（部署时用），没有则用代码里的默认值
func LoadConfig() {
	AppConfig = Config{
		// ========== 数据库连接字符串 ==========
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://neondb_owner:npg_Q5BfgWureh8E@ep-floral-band-b5tfczl0-pooler.c-7.us-east-2.aws.neon.tech/survey_db?sslmode=require"),

		// ========== JWT 签名密钥 ==========
		JWTSecret: getEnv("JWT_SECRET", "survey-system-secret-2024-abcxyz"),

		// ========== 服务端口 ==========
		ServerPort: getEnv("PORT", "88888"),
	}
}

// getEnv 读取环境变量，不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
