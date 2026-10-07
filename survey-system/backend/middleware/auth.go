package middleware

import (
	"net/http"
	"strings"
	"survey-system/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired JWT 鉴权中间件
// 从请求头 Authorization: Bearer <token> 中解析 JWT 令牌
// 验证通过后将 user_id 和 username 存入上下文，供后续处理器使用
// 验证失败返回 401 Unauthorized
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少授权令牌"})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "授权格式错误，应为 Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析并验证 JWT
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或已过期的令牌"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetCurrentUserID 从 Gin 上下文中获取当前用户ID
// 必须在 AuthRequired 中间件之后使用
func GetCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetCurrentUsername 从 Gin 上下文中获取当前用户名
// 必须在 AuthRequired 中间件之后使用
func GetCurrentUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}
