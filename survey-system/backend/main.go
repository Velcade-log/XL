// main.go 问卷系统后端服务入口
// 技术栈: Go + Gin + PostgreSQL
package main

import (
	"log"
	"survey-system/config"
	"survey-system/database"
	"survey-system/handlers"
	"survey-system/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	config.LoadConfig()

	// 校验必要的环境变量
	if config.AppConfig.DatabaseURL == "" {
		log.Fatal("错误: 环境变量 DATABASE_URL 未设置")
	}
	if config.AppConfig.JWTSecret == "" {
		log.Fatal("错误: 环境变量 JWT_SECRET 未设置")
	}

	// 2. 初始化数据库连接
	database.InitDB()
	defer database.CloseDB()

	// 3. 创建 Gin 引擎
	r := gin.Default()

	// 4. 全局中间件
	r.Use(middleware.CORS()) // CORS 跨域支持

	// 5. 路由分组
	api := r.Group("/api")
	{
		// 公开接口（不需要登录）
		api.POST("/login", handlers.Login)

		// 需要鉴权的接口
		auth := api.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// 用户信息
			auth.GET("/me", handlers.GetCurrentUser)

			// 问卷相关
			auth.GET("/surveys", handlers.GetSurveys)       // 获取问卷列表
			auth.POST("/surveys", handlers.CreateSurvey)     // 创建问卷
			auth.GET("/surveys/:id", handlers.GetSurvey)     // 获取问卷详情

			// 答卷相关
			auth.PUT("/surveys/:id/answer", handlers.SubmitAnswer)    // 提交/更新答卷
			auth.POST("/surveys/:id/answer", handlers.SubmitAnswer)   // 兼容 POST 方式
			auth.GET("/surveys/:id/my-answer", handlers.GetMyAnswer)  // 获取自己的答卷
			auth.GET("/surveys/:id/answers", handlers.GetSurveyAnswers) // 获取所有答卷（仅创建者）
		}
	}

	// 6. 启动服务
	port := config.AppConfig.ServerPort
	log.Printf("服务启动成功，监听端口 %s", port)
	log.Printf("API 基础路径: http://localhost:%s/api", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
