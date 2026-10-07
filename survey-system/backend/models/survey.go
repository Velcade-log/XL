package models

import (
	"encoding/json"
	"time"
)

// Question 问卷题目结构
type Question struct {
	ID    string `json:"id"`    // 题目唯一标识，如 "q1", "q2"
	Title string `json:"title"` // 题目标题
	Type  string `json:"type"`  // 题目类型，目前支持 "text"（文本输入）
}

// Survey 问卷模型，对应 surveys 表
type Survey struct {
	ID          uint            `json:"id"`           // 问卷ID
	Title       string          `json:"title"`        // 问卷标题
	Description string          `json:"description"`  // 问卷描述
	Questions   json.RawMessage `json:"questions"`    // 题目数组（JSONB），格式: [{"id":"q1","title":"...","type":"text"}]
	CreatorID   uint            `json:"creator_id"`   // 创建者用户ID
	CreatorName string          `json:"creator_name"` // 创建者用户名（查询时JOIN得到）
	CreatedAt   time.Time       `json:"created_at"`   // 创建时间
}

// CreateSurveyRequest 创建问卷请求体
type CreateSurveyRequest struct {
	Title       string          `json:"title" binding:"required"`       // 问卷标题，必填
	Description string          `json:"description"`                     // 问卷描述
	Questions   json.RawMessage `json:"questions" binding:"required"`   // 题目数组，必填
}

// SurveyListResponse 问卷列表响应（列表页不需要完整题目）
type SurveyListResponse struct {
	ID          uint      `json:"id"`           // 问卷ID
	Title       string    `json:"title"`        // 问卷标题
	Description string    `json:"description"`  // 问卷描述
	CreatorID   uint      `json:"creator_id"`   // 创建者ID
	CreatorName string    `json:"creator_name"` // 创建者用户名
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
}
