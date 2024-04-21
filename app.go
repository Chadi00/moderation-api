package main

import (
	"fmt"
	"moderation_api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("Hello World")

	server := gin.Default()

	routes.GenerateRoutes(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server.Run(":" + port)
}
