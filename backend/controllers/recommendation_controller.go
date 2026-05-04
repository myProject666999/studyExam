package controllers

import (
	"strconv"
	"studyexam/recommendation"
	"studyexam/utils"

	"github.com/gin-gonic/gin"
)

func GetRecommendations(c *gin.Context) {
	userID := c.GetUint("user_id")
	topN, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations := recommendation.GetRecommendedCourses(userID, topN)

	utils.Success(c, recommendations)
}

func GetHotCourses(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations := recommendation.GetRecommendedCourses(0, limit)

	utils.Success(c, recommendations)
}
