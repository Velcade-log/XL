// Package utils 工具函数集合
package utils

import (
	"errors"
	"survey-system/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims JWT 自定义声明，包含用户ID和用户名
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT 令牌
// userID: 用户ID
// username: 用户名
// 返回: 签名字符串和可能的错误
// 令牌有效期 24 小时
func GenerateJWT(userID uint, username string) (string, error) {
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "survey-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(config.AppConfig.JWTSecret)
	return token.SignedString(secret)
}

// ParseJWT 解析并验证 JWT 令牌
// tokenString: JWT 令牌字符串
// 返回: 用户声明和可能的错误
func ParseJWT(tokenString string) (*UserClaims, error) {
	secret := []byte(config.AppConfig.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}
