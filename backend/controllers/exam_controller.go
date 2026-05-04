package controllers

import (
	"strconv"
	"studyexam/database"
	"studyexam/models"
	"studyexam/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetExams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var exams []models.Exam
	var total int64

	query := database.DB.Model(&models.Exam{}).Where("status = ?", 1)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("Paper").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&exams).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, exams, total, page, pageSize)
}

func GetExamDetail(c *gin.Context) {
	id := c.Param("id")

	var exam models.Exam
	if err := database.DB.Preload("Paper").First(&exam, id).Error; err != nil {
		utils.NotFound(c, "考试不存在")
		return
	}

	var paperQuestions []models.PaperQuestion
	database.DB.Where("paper_id = ?", exam.PaperID).Order("sort asc").Find(&paperQuestions)

	var questionIDs []uint
	for _, pq := range paperQuestions {
		questionIDs = append(questionIDs, pq.QuestionID)
	}

	var questions []models.Question
	database.DB.Where("id IN ?", questionIDs).Find(&questions)

	questionMap := make(map[uint]models.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	type QuestionWithScore struct {
		models.Question
		Score int `json:"score"`
	}

	var questionsWithScore []QuestionWithScore
	for _, pq := range paperQuestions {
		if q, exists := questionMap[pq.QuestionID]; exists {
			questionsWithScore = append(questionsWithScore, QuestionWithScore{
				Question: q,
				Score:    pq.Score,
			})
		}
	}

	utils.Success(c, gin.H{
		"exam":       exam,
		"questions":  questionsWithScore,
	})
}

func StartExam(c *gin.Context) {
	userID := c.GetUint("user_id")
	examID := c.Param("id")

	examIDUint, _ := strconv.ParseUint(examID, 10, 64)

	var exam models.Exam
	if err := database.DB.First(&exam, examIDUint).Error; err != nil {
		utils.NotFound(c, "考试不存在")
		return
	}

	var existingRecord models.ExamRecord
	database.DB.Where("user_id = ? AND exam_id = ?", userID, examIDUint).First(&existingRecord)

	if existingRecord.ID > 0 && existingRecord.Status != "pending" {
		utils.BadRequest(c, "您已完成该考试")
		return
	}

	now := time.Now()

	if existingRecord.ID == 0 {
		existingRecord = models.ExamRecord{
			UserID:     userID,
			ExamID:     uint(examIDUint),
			PaperID:    exam.PaperID,
			TotalScore: exam.TotalScore,
			StartTime:  &now,
			Status:     "pending",
		}
		database.DB.Create(&existingRecord)
	}

	utils.Success(c, gin.H{
		"record_id": existingRecord.ID,
		"start_time": existingRecord.StartTime,
	})
}

func SubmitExam(c *gin.Context) {
	userID := c.GetUint("user_id")
	recordID := c.Param("record_id")

	recordIDUint, _ := strconv.ParseUint(recordID, 10, 64)

	var record models.ExamRecord
	if err := database.DB.First(&record, recordIDUint).Error; err != nil {
		utils.NotFound(c, "考试记录不存在")
		return
	}

	if record.UserID != userID {
		utils.Forbidden(c, "无权操作")
		return
	}

	if record.Status != "pending" {
		utils.BadRequest(c, "考试已提交")
		return
	}

	var req struct {
		Answers []struct {
			QuestionID uint   `json:"question_id"`
			Answer     string `json:"answer"`
		} `json:"answers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var paperQuestions []models.PaperQuestion
	database.DB.Where("paper_id = ?", record.PaperID).Find(&paperQuestions)

	questionMap := make(map[uint]models.PaperQuestion)
	for _, pq := range paperQuestions {
		questionMap[pq.QuestionID] = pq
	}

	var questionIDs []uint
	for _, pq := range paperQuestions {
		questionIDs = append(questionIDs, pq.QuestionID)
	}

	var questions []models.Question
	database.DB.Where("id IN ?", questionIDs).Find(&questions)

	answerMap := make(map[uint]models.Question)
	for _, q := range questions {
		answerMap[q.ID] = q
	}

	now := time.Now()
	var totalScore int

	for _, answer := range req.Answers {
		if pq, exists := questionMap[answer.QuestionID]; exists {
			if q, exists := answerMap[answer.QuestionID]; exists {
				isCorrect := 0
				score := 0

				if q.Answer == answer.Answer {
					isCorrect = 1
					score = pq.Score
					totalScore += score
				}

				userAnswer := models.UserAnswer{
					RecordID:   uint(recordIDUint),
					QuestionID: answer.QuestionID,
					UserAnswer: answer.Answer,
					IsCorrect:  isCorrect,
					Score:      score,
				}
				database.DB.Create(&userAnswer)
			}
		}
	}

	isPass := 0
	passScore := record.TotalScore * 60 / 100
	if totalScore >= passScore {
		isPass = 1
	}

	timeSpent := 0
	if record.StartTime != nil {
		timeSpent = int(now.Sub(*record.StartTime).Seconds())
	}

	database.DB.Model(&record).Updates(map[string]interface{}{
		"score":      totalScore,
		"is_pass":    isPass,
		"submit_time": now,
		"time_spent":  timeSpent,
		"status":      "submitted",
	})

	utils.Success(c, gin.H{
		"score":       totalScore,
		"total_score": record.TotalScore,
		"is_pass":     isPass,
	})
}

func GetMyExamRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var records []models.ExamRecord
	var total int64

	query := database.DB.Model(&models.ExamRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("Exam").Preload("User").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, records, total, page, pageSize)
}

func GetExamRecordDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	recordID := c.Param("id")

	recordIDUint, _ := strconv.ParseUint(recordID, 10, 64)

	var record models.ExamRecord
	if err := database.DB.Preload("Exam").First(&record, recordIDUint).Error; err != nil {
		utils.NotFound(c, "记录不存在")
		return
	}

	if record.UserID != userID {
		utils.Forbidden(c, "无权查看")
		return
	}

	var userAnswers []models.UserAnswer
	database.DB.Where("record_id = ?", recordIDUint).Find(&userAnswers)

	utils.Success(c, gin.H{
		"record":  record,
		"answers": userAnswers,
	})
}
