// Package handlers HTTP 请求处理器
package handlers

import (
	"database/sql"
	"net/http"
	"survey-system/database"
	"survey-system/middleware"
	"survey-system/models"
	"survey-system/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login 用户登录接口
// POST /api/login
// 请求体: { "username": "xxx", "password": "xxx" }
// 响应: { "token": "jwt_token", "user_id": 1, "username": "xxx" }
// 密码使用 bcrypt 比对，成功则签发 JWT 令牌
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查询用户
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 验证密码（bcrypt 比对）
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 JWT
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, models.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

// GetCurrentUser 获取当前登录用户信息
// GET /api/me
// 用于前端获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	username := middleware.GetCurrentUsername(c)

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
	})
}
