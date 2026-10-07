package models

import (
	"encoding/json"
	"time"
)

// SurveyAnswer 问卷答卷模型，对应 survey_answers 表
type SurveyAnswer struct {
	ID        uint            `json:"id"`         // 答卷ID
	SurveyID  uint            `json:"survey_id"`  // 所属问卷ID
	UserID    uint            `json:"user_id"`    // 填写用户ID
	Username  string          `json:"username"`   // 填写用户名（查询时JOIN得到）
	Answer    json.RawMessage `json:"answer"`     // 答案内容（JSONB），格式: {"q1":"答案1","q2":"答案2"}
	UpdatedAt time.Time       `json:"updated_at"` // 最后更新时间
}

// SubmitAnswerRequest 提交答卷请求体
type SubmitAnswerRequest struct {
	Answer json.RawMessage `json:"answer" binding:"required"` // 答案内容，必填
}

// AnswerListResponse 答卷列表响应
type AnswerListResponse struct {
	ID        uint            `json:"id"`         // 答卷ID
	SurveyID  uint            `json:"survey_id"`  // 问卷ID
	UserID    uint            `json:"user_id"`    // 填写用户ID
	Username  string          `json:"username"`   // 填写用户名
	Answer    json.RawMessage `json:"answer"`     // 答案内容
	UpdatedAt time.Time       `json:"updated_at"` // 更新时间
}
