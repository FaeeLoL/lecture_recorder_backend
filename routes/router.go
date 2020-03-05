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
	user := new(controllers.UserController)
	apiGroup.GET("/test", user.Test)
	return router
}