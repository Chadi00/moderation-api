package routes

import (
	"fmt"
	"log"
	"os/exec"
)

func downloadVideo(videoURL string, outputPath string) error {
	fmt.Println("Starting download...")
	fmt.Println("Video URL:", videoURL)
	fmt.Println("Output Path:", outputPath)

	cmd := exec.Command("python", "utils/downloadVideo.py", videoURL, outputPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to download video: %v", err)
		return err
	}

	fmt.Println("Download completed successfully")
	return nil
}
