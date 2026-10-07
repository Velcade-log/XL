// Package middleware Gin 中间件集合
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
// 允许所有来源的跨域请求（开发环境友好），支持常见的 HTTP 方法和头部
// 同时处理 OPTIONS 预检请求
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的来源，* 表示允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 设置是否允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理 OPTIONS 预检请求，直接返回 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
