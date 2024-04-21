package routes

import (
	"moderation_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func moderation(ctx *gin.Context) {
	var videoURL models.RequestBody

	err := ctx.BindJSON(&videoURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Did not find the video url in the request body"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "OK"})
}
