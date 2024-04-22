package routes

import "github.com/gin-gonic/gin"

func GenerateRoutes(server *gin.Engine) {
	server.GET("/hello", hello)

	server.GET("/analyze-video", analyzeVideo)
}
