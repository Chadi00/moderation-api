package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rdb *redis.Client

func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("not able to load .env")
	}

	REDISHOST := os.Getenv("REDISHOST")
	REDISPASSWORD := os.Getenv("REDISPASSWORD")
	REDISPORT := os.Getenv("REDISPORT")
	REDISUSER := os.Getenv("REDISUSER")

	// create Redis connection
	rdb = redis.NewClient(&redis.Options{
		Addr:     REDISHOST + ":" + REDISPORT,
		Username: REDISUSER,
		Password: REDISPASSWORD,
		DB:       0,
	})
}

type RequestBody struct {
	VideoURL string `json:"VideoURL"`
}

// Analyze video (with the audio of the video) and return a moderation description of it
func analyzeVideo(ctx *gin.Context) {

	var req RequestBody

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

	ctxR := context.Background()
	existingResults, err := rdb.HGetAll(ctxR, req.VideoURL).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve data from Redis"})
		return
	}

	if len(existingResults) > 0 {
		ctx.JSON(http.StatusOK, existingResults)
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

	// map containing the rating and video moderation description
	fields := map[string]interface{}{
		"rating":      rating,
		"description": overallDescription,
	}

	_, err = rdb.HMSet(ctxR, req.VideoURL, fields).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to store data in Redis"})
		return
	}

	ctx.JSON(http.StatusOK, fields)
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
