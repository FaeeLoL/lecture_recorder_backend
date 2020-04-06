package routes

import (
	"github.com/faeelol/lecture_recorder_backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	apiGroup := router.Group("/api/v1")

	publicGroup := router.Group("/api/v1")

	authGroup := apiGroup.Group("/auth")
	{
		authController := new(controllers.AuthController)
		authMiddleware := authController.Init()
		authGroup.POST("login", authMiddleware.LoginHandler)
		authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
		//authGroup.Use(authMiddleware.MiddlewareFunc())
		apiGroup.Use(authMiddleware.MiddlewareFunc())
	}

	users := apiGroup.Group("/users")
	{
		usersController := new(controllers.UserController)
		publicGroup.POST("/users", usersController.Create)
		users.GET("/self_info", usersController.SelfInfo)
	}

	topics := apiGroup.Group("/topics")
	{
		topicsController := new(controllers.TopicsController)
		topics.POST("/", topicsController.Create)
		topics.GET("/", topicsController.List)
		topics.GET("/:topic_id", topicsController.Get)
		topics.PUT("/:topic_id", topicsController.Put)
		topics.DELETE("/:topic_id", topicsController.Delete)
		courses := topics.Group("/:topic_id/courses")
		coursesController := new(controllers.CourseController)
		courses.POST("/", coursesController.Create)
		courses.GET("/", coursesController.List)
		courses.GET("/:course_id", coursesController.Get)
		courses.PUT("/:course_id", coursesController.Put)
		courses.DELETE("/:course_id", coursesController.Delete)
	}
	return router
}
