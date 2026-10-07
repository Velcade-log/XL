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

// GetSurveys 获取所有问卷列表
// GET /api/surveys
// 所有登录用户可访问，返回所有账号发布的问卷
func GetSurveys(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT s.id, s.title, s.description, s.creator_id, s.created_at, u.username
		FROM surveys s
		JOIN users u ON s.creator_id = u.id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询问卷列表失败"})
		return
	}
	defer rows.Close()

	var surveys []models.SurveyListResponse
	for rows.Next() {
		var s models.SurveyListResponse
		err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.CreatorID, &s.CreatedAt, &s.CreatorName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析问卷数据失败"})
			return
		}
		surveys = append(surveys, s)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "遍历问卷数据失败"})
		return
	}

	// 确保返回空数组而不是 null
	if surveys == nil {
		surveys = []models.SurveyListResponse{}
	}

	c.JSON(http.StatusOK, gin.H{"surveys": surveys})
}

// GetSurvey 获取单份问卷详情（含题目）
// GET /api/surveys/:id
func GetSurvey(c *gin.Context) {
	idParam := c.Param("id")
	surveyID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问卷ID"})
		return
	}

	var survey models.Survey
	err = database.DB.QueryRow(`
		SELECT s.id, s.title, s.description, s.questions, s.creator_id, s.created_at, u.username
		FROM surveys s
		JOIN users u ON s.creator_id = u.id
		WHERE s.id = $1
	`, surveyID).Scan(
		&survey.ID, &survey.Title, &survey.Description, &survey.Questions,
		&survey.CreatorID, &survey.CreatedAt, &survey.CreatorName,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "问卷不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询问卷失败"})
		return
	}

	c.JSON(http.StatusOK, survey)
}

// CreateSurvey 创建新问卷
// POST /api/surveys
// 请求体: { "title": "xxx", "description": "xxx", "questions": [...] }
func CreateSurvey(c *gin.Context) {
	var req models.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证题目格式（简单检查是否为数组）
	if len(req.Questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问卷至少需要一道题目"})
		return
	}

	creatorID := middleware.GetCurrentUserID(c)

	var surveyID int
	err := database.DB.QueryRow(`
		INSERT INTO surveys (title, description, questions, creator_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.Title, req.Description, req.Questions, creatorID).Scan(&surveyID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建问卷失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      surveyID,
		"message": "问卷创建成功",
	})
}
