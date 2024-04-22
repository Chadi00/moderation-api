package routes

import (
	"moderation_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func analyzeVideo(ctx *gin.Context) {
	var req models.RequestBody

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Did not find the video url in the request body"})
		return
	}

	if req.VideoURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "VideoURL is empty"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Video URL": req.VideoURL})
}
