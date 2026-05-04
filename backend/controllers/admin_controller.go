package controllers

import (
	"strconv"
	"studyexam/database"
	"studyexam/models"
	"studyexam/utils"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	UserCount      int64         `json:"user_count"`
	CourseCount    int64         `json:"course_count"`
	ExamCount      int64         `json:"exam_count"`
	ForumCount     int64         `json:"forum_count"`
	VideoTypeStats []TypeStat    `json:"video_type_stats"`
	VideoStats     []VideoStat   `json:"video_stats"`
}

type TypeStat struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type VideoStat struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func GetDashboardStats(c *gin.Context) {
	var userCount, courseCount, examCount, forumCount int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Course{}).Where("status = ?", 1).Count(&courseCount)
	database.DB.Model(&models.Exam{}).Where("status = ?", 1).Count(&examCount)
	database.DB.Model(&models.Forum{}).Where("status = ?", 1).Count(&forumCount)

	var courseTypes []models.CourseType
	database.DB.Where("status = ?", 1).Find(&courseTypes)

	var videoTypeStats []TypeStat
	for _, ct := range courseTypes {
		var count int64
		database.DB.Model(&models.Course{}).Where("type_id = ? AND status = ?", ct.ID, 1).Count(&count)
		videoTypeStats = append(videoTypeStats, TypeStat{
			Name:  ct.Name,
			Value: count,
		})
	}

	var videoStats []VideoStat
	var courses []models.Course
	database.DB.Where("status = ?", 1).Order("view_count desc").Limit(10).Find(&courses)

	for _, c := range courses {
		videoStats = append(videoStats, VideoStat{
			Name:  c.Title,
			Value: int64(c.ViewCount),
		})
	}

	stats := DashboardStats{
		UserCount:      userCount,
		CourseCount:    courseCount,
		ExamCount:      examCount,
		ForumCount:     forumCount,
		VideoTypeStats: videoTypeStats,
		VideoStats:     videoStats,
	}

	utils.Success(c, stats)
}

func AdminGetUsers(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var users []models.User
	var total int64

	database.DB.Model(&models.User{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, users, total, page, pageSize)
}

func AdminUpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminGetCourses(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var courses []models.Course
	var total int64

	database.DB.Model(&models.Course{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Preload("CourseType").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, courses, total, page, pageSize)
}

func AdminCreateCourse(c *gin.Context) {
	var course models.Course

	if err := c.ShouldBindJSON(&course); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&course).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": course.ID})
}

func AdminUpdateCourse(c *gin.Context) {
	id := c.Param("id")

	var course models.Course
	if err := database.DB.First(&course, id).Error; err != nil {
		utils.NotFound(c, "课程不存在")
		return
	}

	if err := c.ShouldBindJSON(&course); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&course).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteCourse(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Course{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetCourseTypes(c *gin.Context) {
	var types []models.CourseType

	if err := database.DB.Order("sort asc").Find(&types).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.Success(c, types)
}

func AdminCreateCourseType(c *gin.Context) {
	var courseType models.CourseType

	if err := c.ShouldBindJSON(&courseType); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&courseType).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": courseType.ID})
}

func AdminUpdateCourseType(c *gin.Context) {
	id := c.Param("id")

	var courseType models.CourseType
	if err := database.DB.First(&courseType, id).Error; err != nil {
		utils.NotFound(c, "类型不存在")
		return
	}

	if err := c.ShouldBindJSON(&courseType); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&courseType).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteCourseType(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.CourseType{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetForums(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var forums []models.Forum
	var total int64

	database.DB.Model(&models.Forum{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Preload("User").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&forums).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, forums, total, page, pageSize)
}

func AdminUpdateForumStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Model(&models.Forum{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteForum(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Forum{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetAnnouncements(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var announcements []models.Announcement
	var total int64

	database.DB.Model(&models.Announcement{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Order("is_top desc, sort asc, created_at desc").Offset(offset).Limit(pageSize).Find(&announcements).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, announcements, total, page, pageSize)
}

func AdminCreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement

	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": announcement.ID})
}

func AdminUpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")

	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&announcement).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Announcement{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetBanners(c *gin.Context) {
	var banners []models.Banner

	if err := database.DB.Order("sort asc").Find(&banners).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.Success(c, banners)
}

func AdminCreateBanner(c *gin.Context) {
	var banner models.Banner

	if err := c.ShouldBindJSON(&banner); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&banner).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": banner.ID})
}

func AdminUpdateBanner(c *gin.Context) {
	id := c.Param("id")

	var banner models.Banner
	if err := database.DB.First(&banner, id).Error; err != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	if err := c.ShouldBindJSON(&banner); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&banner).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Banner{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetExams(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var exams []models.Exam
	var total int64

	database.DB.Model(&models.Exam{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Preload("Paper").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&exams).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, exams, total, page, pageSize)
}

func AdminGetExamRecords(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var records []models.ExamRecord
	var total int64

	database.DB.Model(&models.ExamRecord{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Preload("Exam").Preload("User").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, records, total, page, pageSize)
}

func AdminGetPapers(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	var papers []models.Paper
	var total int64

	database.DB.Model(&models.Paper{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&papers).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, papers, total, page, pageSize)
}

func AdminCreatePaper(c *gin.Context) {
	var paper models.Paper

	if err := c.ShouldBindJSON(&paper); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&paper).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": paper.ID})
}

func AdminUpdatePaper(c *gin.Context) {
	id := c.Param("id")

	var paper models.Paper
	if err := database.DB.First(&paper, id).Error; err != nil {
		utils.NotFound(c, "试卷不存在")
		return
	}

	if err := c.ShouldBindJSON(&paper); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&paper).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeletePaper(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Paper{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetQuestions(c *gin.Context) {
	page, pageSize := getPaginationParams(c)
	bankID := c.Query("bank_id")

	var questions []models.Question
	var total int64

	query := database.DB.Model(&models.Question{})
	if bankID != "" {
		query = query.Where("bank_id = ?", bankID)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&questions).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, questions, total, page, pageSize)
}

func AdminCreateQuestion(c *gin.Context) {
	var question models.Question

	if err := c.ShouldBindJSON(&question); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&question).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": question.ID})
}

func AdminUpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	var question models.Question
	if err := database.DB.First(&question, id).Error; err != nil {
		utils.NotFound(c, "试题不存在")
		return
	}

	if err := c.ShouldBindJSON(&question); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&question).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Question{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func AdminGetQuestionBanks(c *gin.Context) {
	var banks []models.QuestionBank

	if err := database.DB.Order("created_at desc").Find(&banks).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.Success(c, banks)
}

func AdminCreateQuestionBank(c *gin.Context) {
	var bank models.QuestionBank

	if err := c.ShouldBindJSON(&bank); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&bank).Error; err != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", gin.H{"id": bank.ID})
}

func AdminUpdateQuestionBank(c *gin.Context) {
	id := c.Param("id")

	var bank models.QuestionBank
	if err := database.DB.First(&bank, id).Error; err != nil {
		utils.NotFound(c, "题库不存在")
		return
	}

	if err := c.ShouldBindJSON(&bank); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&bank).Error; err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func AdminDeleteQuestionBank(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.QuestionBank{}, id).Error; err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func getPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}
