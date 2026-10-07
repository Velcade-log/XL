-- ============================================
-- 问卷系统数据库初始化脚本
-- PostgreSQL 12+
-- ============================================

-- 启用 pgcrypto 扩展（用于 gen_random_uuid()，PostgreSQL 13+ 内置可用）
-- CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- --------------------------------------------
-- 1. 用户表 users
-- --------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- --------------------------------------------
-- 2. 问卷表 surveys
--    questions 字段存储 JSONB 格式的题目数组
--    格式示例: [{"id": "q1", "title": "您的姓名", "type": "text"}]
-- --------------------------------------------
CREATE TABLE IF NOT EXISTS surveys (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    questions JSONB NOT NULL DEFAULT '[]'::jsonb,
    creator_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_surveys_creator_id ON surveys(creator_id);

-- --------------------------------------------
-- 3. 问卷答卷表 survey_answers
--    answer 字段存储 JSONB 格式的答案
--    格式示例: {"q1": "张三", "q2": "答案内容"}
-- --------------------------------------------
CREATE TABLE IF NOT EXISTS survey_answers (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    answer JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 联合唯一约束：同一用户对同一份问卷只能有一份答卷
    UNIQUE(survey_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_survey_answers_survey_id ON survey_answers(survey_id);
CREATE INDEX IF NOT EXISTS idx_survey_answers_user_id ON survey_answers(user_id);

-- ============================================
-- 初始化用户说明
-- ============================================
-- 系统预设 2 个用户，不提供注册功能。
-- 密码使用 bcrypt 哈希，需要先生成哈希值再插入。
--
-- 生成 bcrypt 密码的方式：
-- 方式一：使用 Go 代码生成（推荐）
--   package main
--   import (
--       "fmt"
--       "golang.org/x/crypto/bcrypt"
--   )
--   func main() {
--       hash, _ := bcrypt.GenerateFromPassword([]byte("你的密码"), bcrypt.DefaultCost)
--       fmt.Println(string(hash))
--   }
--
-- 方式二：使用在线 bcrypt 生成工具（如 https://bcrypt-generator.com/）
--
-- 方式三：使用 PostgreSQL 的 pgcrypto 扩展（需先安装）
--   SELECT crypt('your_password', gen_salt('bf'));
--
-- 以下为示例 SQL（请替换为实际生成的 bcrypt 哈希值）：
-- ============================================

-- 示例：插入用户 admin（密码 admin123）和 user（密码 user123）
-- 注意：下面的哈希值是示例，请替换为你自己生成的哈希！
--
-- INSERT INTO users (username, password_hash) VALUES
--     ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
--     ('user',  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy');
--
-- 上面哈希值对应的明文密码都是 "password"，仅作示例，请勿在生产环境使用！

-- 建议使用下面的占位符，生成哈希后再执行：
-- INSERT INTO users (username, password_hash) VALUES
--     ('user1', '<替换为user1密码的bcrypt哈希>'),
--     ('user2', '<替换为user2密码的bcrypt哈希>');

-- ============================================
-- 测试数据（可选）
-- ============================================

-- 插入一份示例问卷（需要先有用户，假设用户ID为1）
-- INSERT INTO surveys (title, description, questions, creator_id) VALUES (
--     '用户满意度调查',
--     '请填写以下问卷，帮助我们改进服务',
--     '[{"id":"q1","title":"您的姓名","type":"text"},{"id":"q2","title":"您对我们的服务满意吗？","type":"text"},{"id":"q3","title":"您有什么建议？","type":"text"}]'::jsonb,
--     1
-- );
