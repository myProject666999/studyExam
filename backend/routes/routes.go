package routes

import (
	"studyexam/controllers"
	"studyexam/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/user/login", controllers.UserLogin)
		api.POST("/user/register", controllers.UserRegister)
		api.POST("/admin/login", controllers.AdminLogin)

		api.GET("/banners", controllers.GetBanners)
		api.GET("/announcements", controllers.GetAnnouncements)
		api.GET("/announcements/:id", controllers.GetAnnouncementDetail)
		api.GET("/course-types", controllers.GetCourseTypes)
		api.GET("/courses", controllers.GetCourses)
		api.GET("/courses/:id", controllers.GetCourseDetail)
		api.GET("/sections/:id", controllers.GetSectionDetail)
		api.GET("/sections/:id/comments", controllers.GetSectionComments)
		api.GET("/hot-courses", controllers.GetHotCourses)
		api.GET("/exams", controllers.GetExams)
		api.GET("/forums", controllers.GetForums)
		api.GET("/forums/:id", controllers.GetForumDetail)
		api.GET("/forums/:id/replies", controllers.GetForumReplies)

		user := api.Group("")
		user.Use(middleware.JWTAuth())
		{
			user.GET("/user/info", controllers.GetUserInfo)
			user.PUT("/user/info", controllers.UpdateUserInfo)
			user.PUT("/user/password", controllers.UpdatePassword)

			user.GET("/recommendations", controllers.GetRecommendations)

			user.POST("/favorites/toggle", controllers.ToggleFavorite)
			user.GET("/favorites", controllers.GetMyFavorites)
			user.GET("/favorites/check", controllers.CheckFavorite)

			user.POST("/sections/:id/comments", controllers.AddComment)

			user.GET("/exams/:id", controllers.GetExamDetail)
			user.POST("/exams/:id/start", controllers.StartExam)
			user.POST("/exam-records/:record_id/submit", controllers.SubmitExam)
			user.GET("/exam-records", controllers.GetMyExamRecords)
			user.GET("/exam-records/:id", controllers.GetExamRecordDetail)

			user.POST("/forums", controllers.CreateForum)
			user.POST("/forums/:id/replies", controllers.CreateReply)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			admin.GET("/dashboard", controllers.GetDashboardStats)

			admin.GET("/users", controllers.AdminGetUsers)
			admin.PUT("/users/:id/status", controllers.AdminUpdateUserStatus)

			admin.GET("/course-types", controllers.AdminGetCourseTypes)
			admin.POST("/course-types", controllers.AdminCreateCourseType)
			admin.PUT("/course-types/:id", controllers.AdminUpdateCourseType)
			admin.DELETE("/course-types/:id", controllers.AdminDeleteCourseType)

			admin.GET("/courses", controllers.AdminGetCourses)
			admin.POST("/courses", controllers.AdminCreateCourse)
			admin.PUT("/courses/:id", controllers.AdminUpdateCourse)
			admin.DELETE("/courses/:id", controllers.AdminDeleteCourse)

			admin.GET("/forums", controllers.AdminGetForums)
			admin.PUT("/forums/:id/status", controllers.AdminUpdateForumStatus)
			admin.DELETE("/forums/:id", controllers.AdminDeleteForum)

			admin.GET("/announcements", controllers.AdminGetAnnouncements)
			admin.POST("/announcements", controllers.AdminCreateAnnouncement)
			admin.PUT("/announcements/:id", controllers.AdminUpdateAnnouncement)
			admin.DELETE("/announcements/:id", controllers.AdminDeleteAnnouncement)

			admin.GET("/banners", controllers.AdminGetBanners)
			admin.POST("/banners", controllers.AdminCreateBanner)
			admin.PUT("/banners/:id", controllers.AdminUpdateBanner)
			admin.DELETE("/banners/:id", controllers.AdminDeleteBanner)

			admin.GET("/exams", controllers.AdminGetExams)
			admin.GET("/exam-records", controllers.AdminGetExamRecords)

			admin.GET("/papers", controllers.AdminGetPapers)
			admin.POST("/papers", controllers.AdminCreatePaper)
			admin.PUT("/papers/:id", controllers.AdminUpdatePaper)
			admin.DELETE("/papers/:id", controllers.AdminDeletePaper)

			admin.GET("/question-banks", controllers.AdminGetQuestionBanks)
			admin.POST("/question-banks", controllers.AdminCreateQuestionBank)
			admin.PUT("/question-banks/:id", controllers.AdminUpdateQuestionBank)
			admin.DELETE("/question-banks/:id", controllers.AdminDeleteQuestionBank)

			admin.GET("/questions", controllers.AdminGetQuestions)
			admin.POST("/questions", controllers.AdminCreateQuestion)
			admin.PUT("/questions/:id", controllers.AdminUpdateQuestion)
			admin.DELETE("/questions/:id", controllers.AdminDeleteQuestion)
		}
	}

	return r
}
