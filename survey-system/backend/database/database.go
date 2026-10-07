// Package database 数据库连接管理
package database

import (
	"database/sql"
	"log"
	"survey-system/config"

	_ "github.com/lib/pq" // PostgreSQL 驱动
)

// DB 全局数据库连接实例
var DB *sql.DB

// InitDB 初始化数据库连接
// 使用 lib/pq 驱动连接 PostgreSQL，连接字符串从环境变量 DATABASE_URL 读取
func InitDB() {
	var err error

	// 打开数据库连接
	DB, err = sql.Open("postgres", config.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 测试连接是否可用
	err = DB.Ping()
	if err != nil {
		log.Fatalf("数据库 ping 失败: %v", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(25)  // 最大打开连接数
	DB.SetMaxIdleConns(10)  // 最大空闲连接数

	log.Println("数据库连接成功")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("数据库连接已关闭")
	}
}
