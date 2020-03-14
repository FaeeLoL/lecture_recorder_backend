package routes

import (
	"github.com/faeelol/lecture_recorder_backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	apiGroup := router.Group("/api")



	topics := apiGroup.Group("/topics")
	{
		topicsController := new(controllers.TopicsController)
		topics.POST("/", topicsController.Create)
		topics.GET("/", topicsController.List)
		topics.GET("/:id", topicsController.Get)
		topics.PUT("/:id", topicsController.Put)
		topics.DELETE("/:id", topicsController.Delete)
	}
	return router
}