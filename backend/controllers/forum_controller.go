package controllers

import (
	"strconv"
	"studyexam/database"
	"studyexam/models"
	"studyexam/utils"

	"github.com/gin-gonic/gin"
)

func GetForums(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	var forums []models.Forum
	var total int64

	query := database.DB.Model(&models.Forum{}).Where("status = ?", 1)

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&forums).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, forums, total, page, pageSize)
}

func GetForumDetail(c *gin.Context) {
	id := c.Param("id")

	var forum models.Forum
	if err := database.DB.Preload("User").First(&forum, id).Error; err != nil {
		utils.NotFound(c, "帖子不存在")
		return
	}

	database.DB.Model(&forum).Update("view_count", forum.ViewCount+1)

	utils.Success(c, forum)
}

func GetForumReplies(c *gin.Context) {
	forumID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var replies []models.ForumReply
	var total int64

	query := database.DB.Model(&models.ForumReply{}).Where("forum_id = ? AND status = ?", forumID, 1)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at asc").Offset(offset).Limit(pageSize).Find(&replies).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, replies, total, page, pageSize)
}

func CreateForum(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	forum := models.Forum{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  1,
	}

	if err := database.DB.Create(&forum).Error; err != nil {
		utils.InternalError(c, "发布失败")
		return
	}

	utils.SuccessWithMessage(c, "发布成功", gin.H{
		"id": forum.ID,
	})
}

func CreateReply(c *gin.Context) {
	userID := c.GetUint("user_id")
	forumID := c.Param("id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	forumIDUint, _ := strconv.ParseUint(forumID, 10, 64)

	var forum models.Forum
	if err := database.DB.First(&forum, forumIDUint).Error; err != nil {
		utils.NotFound(c, "帖子不存在")
		return
	}

	reply := models.ForumReply{
		ForumID:  uint(forumIDUint),
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
		Status:   1,
	}

	if err := database.DB.Create(&reply).Error; err != nil {
		utils.InternalError(c, "回复失败")
		return
	}

	database.DB.Model(&forum).Update("reply_count", forum.ReplyCount+1)

	utils.SuccessWithMessage(c, "回复成功", nil)
}
