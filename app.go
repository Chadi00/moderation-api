package main

import (
	"fmt"
	"moderation_api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("Hello World")

	server := gin.Default()

	routes.GenerateRoutes(server)

	server.Run(":8080")
}
