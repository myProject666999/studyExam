package recommendation

import (
	"math"
	"sort"
	"studyexam/database"
	"studyexam/models"
)

type UserRating struct {
	UserID   uint
	CourseID uint
	Rating   float64
}

type CourseScore struct {
	CourseID uint
	Score    float64
}

type Recommendation struct {
	Course models.Course
	Score  float64
}

func BuildUserItemMatrix() map[uint]map[uint]float64 {
	var behaviors []models.UserBehavior
	database.DB.Find(&behaviors)

	userItemMatrix := make(map[uint]map[uint]float64)

	for _, behavior := range behaviors {
		if userItemMatrix[behavior.UserID] == nil {
			userItemMatrix[behavior.UserID] = make(map[uint]float64)
		}

		var rating float64
		switch behavior.BehaviorType {
		case "view":
			rating = 2.0
		case "collect":
			rating = 5.0
		case "comment":
			rating = 4.0
		case "study":
			rating = 3.0
		default:
			rating = 1.0
		}

		existingRating := userItemMatrix[behavior.UserID][behavior.CourseID]
		userItemMatrix[behavior.UserID][behavior.CourseID] = math.Max(existingRating, rating)
	}

	return userItemMatrix
}

func CosineSimilarity(vec1, vec2 map[uint]float64) float64 {
	var dotProduct, norm1, norm2 float64

	for item, rating1 := range vec1 {
		if rating2, exists := vec2[item]; exists {
			dotProduct += rating1 * rating2
		}
		norm1 += rating1 * rating1
	}

	for _, rating2 := range vec2 {
		norm2 += rating2 * rating2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func FindSimilarUsers(userID uint, userItemMatrix map[uint]map[uint]float64, topN int) []uint {
	userVector := userItemMatrix[userID]
	if userVector == nil {
		return nil
	}

	type UserSimilarity struct {
		UserID     uint
		Similarity float64
	}

	var similarities []UserSimilarity

	for otherUserID, otherVector := range userItemMatrix {
		if otherUserID == userID {
			continue
		}

		similarity := CosineSimilarity(userVector, otherVector)
		if similarity > 0 {
			similarities = append(similarities, UserSimilarity{
				UserID:     otherUserID,
				Similarity: similarity,
			})
		}
	}

	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	var similarUsers []uint
	for i := 0; i < len(similarities) && i < topN; i++ {
		similarUsers = append(similarUsers, similarities[i].UserID)
	}

	return similarUsers
}

func RecommendForUser(userID uint, topN int) []Recommendation {
	userItemMatrix := BuildUserItemMatrix()

	userVector := userItemMatrix[userID]

	if userVector == nil || len(userVector) == 0 {
		return recommendPopularCourses(topN)
	}

	similarUsers := FindSimilarUsers(userID, userItemMatrix, 10)

	if len(similarUsers) == 0 {
		return recommendPopularCourses(topN)
	}

	itemScores := make(map[uint]float64)
	itemSimSum := make(map[uint]float64)

	for _, similarUserID := range similarUsers {
		similarVector := userItemMatrix[similarUserID]
		similarity := CosineSimilarity(userVector, similarVector)

		for courseID, rating := range similarVector {
			if _, exists := userVector[courseID]; exists {
				continue
			}

			itemScores[courseID] += similarity * rating
			itemSimSum[courseID] += similarity
		}
	}

	var courseScores []CourseScore
	for courseID, score := range itemScores {
		if itemSimSum[courseID] > 0 {
			normalizedScore := score / itemSimSum[courseID]
			courseScores = append(courseScores, CourseScore{
				CourseID: courseID,
				Score:    normalizedScore,
			})
		}
	}

	sort.Slice(courseScores, func(i, j int) bool {
		return courseScores[i].Score > courseScores[j].Score
	})

	var courseIDs []uint
	for i := 0; i < len(courseScores) && i < topN; i++ {
		courseIDs = append(courseIDs, courseScores[i].CourseID)
	}

	if len(courseIDs) == 0 {
		return recommendPopularCourses(topN)
	}

	var courses []models.Course
	database.DB.Where("id IN ? AND status = ?", courseIDs, 1).Preload("CourseType").Find(&courses)

	courseMap := make(map[uint]models.Course)
	for _, course := range courses {
		courseMap[course.ID] = course
	}

	var recommendations []Recommendation
	for i := 0; i < len(courseScores) && i < topN; i++ {
		if course, exists := courseMap[courseScores[i].CourseID]; exists {
			recommendations = append(recommendations, Recommendation{
				Course: course,
				Score:  courseScores[i].Score,
			})
		}
	}

	if len(recommendations) < topN {
		remaining := topN - len(recommendations)
		popular := recommendPopularCourses(remaining)
		recommendations = append(recommendations, popular...)
	}

	return recommendations
}

func recommendPopularCourses(topN int) []Recommendation {
	var courses []models.Course
	database.DB.Where("status = ?", 1).Preload("CourseType").
		Order("view_count desc, collect_count desc").
		Limit(topN).Find(&courses)

	var recommendations []Recommendation
	for _, course := range courses {
		score := float64(course.ViewCount)/1000.0 + float64(course.CollectCount)/100.0
		recommendations = append(recommendations, Recommendation{
			Course: course,
			Score:  score,
		})
	}

	return recommendations
}

func GetRecommendedCourses(userID uint, topN int) []Recommendation {
	return RecommendForUser(userID, topN)
}
