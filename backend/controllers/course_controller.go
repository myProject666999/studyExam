package controllers

import (
	"strconv"
	"studyexam/database"
	"studyexam/models"
	"studyexam/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCourseTypes(c *gin.Context) {
	var types []models.CourseType
	if err := database.DB.Where("status = ?", 1).Order("sort asc").Find(&types).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.Success(c, types)
}

func GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	typeID := c.Query("type_id")
	keyword := c.Query("keyword")

	var courses []models.Course
	var total int64

	query := database.DB.Model(&models.Course{}).Where("status = ?", 1)

	if typeID != "" {
		query = query.Where("type_id = ?", typeID)
	}

	if keyword != "" {
		query = query.Where("title LIKE ? OR author LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("CourseType").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&courses).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, courses, total, page, pageSize)
}

func GetCourseDetail(c *gin.Context) {
	id := c.Param("id")

	var course models.Course
	if err := database.DB.Preload("CourseType").First(&course, id).Error; err != nil {
		utils.NotFound(c, "课程不存在")
		return
	}

	database.DB.Model(&course).Update("view_count", course.ViewCount+1)

	var chapters []models.Chapter
	database.DB.Where("course_id = ?", course.ID).Preload("Sections").Order("sort asc").Find(&chapters)

	utils.Success(c, gin.H{
		"course":   course,
		"chapters": chapters,
	})
}

func GetSectionDetail(c *gin.Context) {
	id := c.Param("id")

	var section models.Section
	if err := database.DB.First(&section, id).Error; err != nil {
		utils.NotFound(c, "小节不存在")
		return
	}

	database.DB.Model(&section).Update("view_count", section.ViewCount+1)

	userID, exists := c.Get("user_id")
	if exists {
		var behavior models.UserBehavior
		database.DB.Where("user_id = ? AND course_id = ? AND behavior_type = ?", userID, section.CourseID, "view").First(&behavior)

		if behavior.ID == 0 {
			behavior = models.UserBehavior{
				UserID:       userID.(uint),
				CourseID:     section.CourseID,
				BehaviorType: "view",
				Rating:       2,
			}
			database.DB.Create(&behavior)
		} else {
			database.DB.Model(&behavior).Update("rating", behavior.Rating+1)
		}
	}

	utils.Success(c, section)
}

func GetSectionComments(c *gin.Context) {
	sectionID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{}).Where("section_id = ? AND status = ?", sectionID, 1)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, comments, total, page, pageSize)
}

func AddComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	sectionID := c.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	sectionIDUint, _ := strconv.ParseUint(sectionID, 10, 64)

	var section models.Section
	if err := database.DB.First(&section, sectionIDUint).Error; err != nil {
		utils.NotFound(c, "视频不存在")
		return
	}

	comment := models.Comment{
		UserID:    userID,
		SectionID: uint(sectionIDUint),
		CourseID:  section.CourseID,
		Content:   req.Content,
		Status:    1,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		utils.InternalError(c, "评论失败")
		return
	}

	var behavior models.UserBehavior
	database.DB.Where("user_id = ? AND course_id = ? AND behavior_type = ?", userID, section.CourseID, "comment").First(&behavior)

	if behavior.ID == 0 {
		behavior = models.UserBehavior{
			UserID:       userID,
			CourseID:     section.CourseID,
			BehaviorType: "comment",
			Rating:       4,
		}
		database.DB.Create(&behavior)
	}

	utils.SuccessWithMessage(c, "评论成功", nil)
}

func ToggleFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		FavoriteType string `json:"favorite_type" binding:"required"`
		TargetID     uint   `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var favorite models.Favorite
	result := database.DB.Where("user_id = ? AND favorite_type = ? AND target_id = ?", userID, req.FavoriteType, req.TargetID).First(&favorite)

	if result.Error == nil {
		database.DB.Delete(&favorite)

		if req.FavoriteType == "course" {
			database.DB.Model(&models.Course{}).Where("id = ?", req.TargetID).UpdateColumn("collect_count", gorm.Expr("collect_count - 1"))
		}

		utils.SuccessWithMessage(c, "取消收藏成功", gin.H{"is_favorite": false})
		return
	}

	favorite = models.Favorite{
		UserID:       userID,
		FavoriteType: req.FavoriteType,
		TargetID:     req.TargetID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		utils.InternalError(c, "收藏失败")
		return
	}

	if req.FavoriteType == "course" {
		database.DB.Model(&models.Course{}).Where("id = ?", req.TargetID).UpdateColumn("collect_count", gorm.Expr("collect_count + 1"))

		var course models.Course
		database.DB.First(&course, req.TargetID)

		var behavior models.UserBehavior
		database.DB.Where("user_id = ? AND course_id = ? AND behavior_type = ?", userID, course.ID, "collect").First(&behavior)

		if behavior.ID == 0 {
			behavior = models.UserBehavior{
				UserID:       userID,
				CourseID:     course.ID,
				BehaviorType: "collect",
				Rating:       5,
			}
			database.DB.Create(&behavior)
		}
	}

	utils.SuccessWithMessage(c, "收藏成功", gin.H{"is_favorite": true})
}

func GetMyFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	favoriteType := c.DefaultQuery("type", "course")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var favorites []models.Favorite
	var total int64

	query := database.DB.Model(&models.Favorite{}).Where("user_id = ? AND favorite_type = ?", userID, favoriteType)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&favorites).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	var targetIDs []uint
	for _, f := range favorites {
		targetIDs = append(targetIDs, f.TargetID)
	}

	var result interface{}
	if favoriteType == "course" && len(targetIDs) > 0 {
		var courses []models.Course
		database.DB.Where("id IN ?", targetIDs).Preload("CourseType").Find(&courses)
		result = courses
	}

	utils.SuccessPage(c, result, total, page, pageSize)
}

func CheckFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	favoriteType := c.Query("type")
	targetID := c.Query("target_id")

	var favorite models.Favorite
	result := database.DB.Where("user_id = ? AND favorite_type = ? AND target_id = ?", userID, favoriteType, targetID).First(&favorite)

	utils.Success(c, gin.H{
		"is_favorite": result.Error == nil,
	})
}
