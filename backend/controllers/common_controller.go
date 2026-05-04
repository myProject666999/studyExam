package controllers

import (
	"strconv"
	"studyexam/database"
	"studyexam/models"
	"studyexam/utils"

	"github.com/gin-gonic/gin"
)

func GetBanners(c *gin.Context) {
	var banners []models.Banner
	if err := database.DB.Where("status = ?", 1).Order("sort asc").Find(&banners).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.Success(c, banners)
}

func GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var announcements []models.Announcement
	var total int64

	query := database.DB.Model(&models.Announcement{}).Where("status = ?", 1)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("is_top desc, sort asc, created_at desc").Offset(offset).Limit(pageSize).Find(&announcements).Error; err != nil {
		utils.InternalError(c, "查询失败")
		return
	}

	utils.SuccessPage(c, announcements, total, page, pageSize)
}

func GetAnnouncementDetail(c *gin.Context) {
	id := c.Param("id")

	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	database.DB.Model(&announcement).Update("view_count", announcement.ViewCount+1)

	utils.Success(c, announcement)
}
