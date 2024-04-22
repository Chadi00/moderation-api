package main

import (
	"moderation_api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	routes.GenerateRoutes(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server.Run(":" + port)
}
