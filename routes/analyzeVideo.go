package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type RequestBody struct {
	VideoURL string `json:"VideoURL"`
}

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "viaduct.proxy.rlwy.net:37481",
		Password: "",
		DB:       0,
	})
	return rdb
}

// Analyze video (with the audio of the audio) and return a moderation description of it
func analyzeVideo(ctx *gin.Context) {
	var req RequestBody

	rdb := initRedis()

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

	videoDescription, err := captureFrames(videoPath)
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

	transcript, err := transcriptAudio(audioPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to transcript audio"})
		return
	}

	rating, overallDescription, err := getVideoDescription(videoDescription, transcript)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to get video moderation description"})
		return
	}

	fields := map[string]interface{}{
		"rating":      rating,
		"description": overallDescription,
	}

	_, err = rdb.HSet(req.VideoURL, "mod", fields).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save data to Redis"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Video rating": rating, "Moderation description": overallDescription})
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
