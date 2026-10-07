package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"survey-system/database"
	"survey-system/middleware"
	"survey-system/models"

	"github.com/gin-gonic/gin"
)

// SubmitAnswer 提交/更新当前用户的答卷（Upsert）
// PUT /api/surveys/:id/answer
// POST /api/surveys/:id/answer
// 如果已有答卷则更新，没有则新增（基于 UNIQUE(survey_id, user_id) 约束）
func SubmitAnswer(c *gin.Context) {
	idParam := c.Param("id")
	surveyID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问卷ID"})
		return
	}

	// 检查问卷是否存在
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM surveys WHERE id = $1)", surveyID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查问卷失败"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "问卷不存在"})
		return
	}

	var req models.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	userID := middleware.GetCurrentUserID(c)

	// 使用 PostgreSQL 的 INSERT ... ON CONFLICT 实现 Upsert
	var answerID int
	err = database.DB.QueryRow(`
		INSERT INTO survey_answers (survey_id, user_id, answer, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (survey_id, user_id) DO UPDATE
		SET answer = EXCLUDED.answer, updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, surveyID, userID, req.Answer).Scan(&answerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交答卷失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      answerID,
		"message": "答卷提交成功",
	})
}

// GetMyAnswer 获取当前用户在该问卷的答卷
// GET /api/surveys/:id/my-answer
func GetMyAnswer(c *gin.Context) {
	idParam := c.Param("id")
	surveyID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问卷ID"})
		return
	}

	userID := middleware.GetCurrentUserID(c)

	var answer models.SurveyAnswer
	err = database.DB.QueryRow(`
		SELECT id, survey_id, user_id, answer, updated_at
		FROM survey_answers
		WHERE survey_id = $1 AND user_id = $2
	`, surveyID, userID).Scan(&answer.ID, &answer.SurveyID, &answer.UserID, &answer.Answer, &answer.UpdatedAt)

	if err == sql.ErrNoRows {
		// 没有答卷返回空对象，前端可以据此判断是否首次填写
		c.JSON(http.StatusOK, gin.H{"answer": nil})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询答卷失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}

// GetSurveyAnswers 获取问卷下所有答卷
// GET /api/surveys/:id/answers
// 【权限校验】仅问卷创建者可以访问，否则返回 403 Forbidden
// 这是核心权限控制，必须在后端校验，不能依赖前端
func GetSurveyAnswers(c *gin.Context) {
	idParam := c.Param("id")
	surveyID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问卷ID"})
		return
	}

	currentUserID := middleware.GetCurrentUserID(c)

	// ========== 权限校验 ==========
	// 查询问卷的创建者，判断当前用户是否为创建者
	var creatorID int
	err = database.DB.QueryRow(
		"SELECT creator_id FROM surveys WHERE id = $1",
		surveyID,
	).Scan(&creatorID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "问卷不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询问卷信息失败"})
		return
	}

	// 核心权限判断：只有问卷创建者才能查看所有答卷
	if uint(creatorID) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看该问卷的答卷"})
		return
	}
	// ================================

	// 查询所有答卷，关联用户表获取用户名
	rows, err := database.DB.Query(`
		SELECT sa.id, sa.survey_id, sa.user_id, sa.answer, sa.updated_at, u.username
		FROM survey_answers sa
		JOIN users u ON sa.user_id = u.id
		WHERE sa.survey_id = $1
		ORDER BY sa.updated_at DESC
	`, surveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询答卷列表失败"})
		return
	}
	defer rows.Close()

	var answers []models.AnswerListResponse
	for rows.Next() {
		var a models.AnswerListResponse
		err := rows.Scan(&a.ID, &a.SurveyID, &a.UserID, &a.Answer, &a.UpdatedAt, &a.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析答卷数据失败"})
			return
		}
		answers = append(answers, a)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "遍历答卷数据失败"})
		return
	}

	// 确保返回空数组而不是 null
	if answers == nil {
		answers = []models.AnswerListResponse{}
	}

	c.JSON(http.StatusOK, gin.H{"answers": answers})
}
