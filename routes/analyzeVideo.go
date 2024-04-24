package routes

import (
	"fmt"
	"moderation_api/models"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
)

// Analyze video (with the audio of the audio) and return a moderation description of it
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

	if !isValidURL(req.VideoURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "The video url is not valid, make sure it's a Youtube video url"})
		return
	}

	outputPath := filepath.Join("/app/downloads")
	fmt.Println("Output path is : ", outputPath)

	videoPath, err := downloadVideo(req.VideoURL, outputPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Error when trying to download the video, make sure it's a valid youtube url"})
		return
	}

	numberFrames, err := captureFrames(videoPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Error when trying to capture frames"})
		return
	}

	audioPath, err := extractAudio(videoPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Error when trying to extract audio"})
		return
	}

	err = deleteVideo(videoPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Error when trying to delete video"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Video URL": req.VideoURL, "Number of frames": numberFrames, "Audio path": audioPath})
}

// make sure VideoURL is a valid url for youtube and tiktok videos
func isValidURL(videoURL string) bool {
	parsedURL, err := url.ParseRequestURI(videoURL)
	if err != nil {
		return false
	}

	youtubeRegex := regexp.MustCompile(`^(https?:\/\/)?(www\.)?(youtube\.com\/watch\?v=|youtu\.be\/)[^\s]+$`)

	return youtubeRegex.MatchString(parsedURL.String())
}
