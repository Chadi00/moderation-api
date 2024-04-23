package routes

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func downloadVideo(videoURL string, outputPath string) error {
	fmt.Println("Starting download...")
	fmt.Println("Video URL:", videoURL)
	fmt.Println("Output Path:", outputPath)

	// Specify the full path to the Python executable
	cmd := exec.Command("python", "utils/downloadVideo.py", videoURL, outputPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to download video: %v", err)
		log.Printf("stdout: %s", stdout.String())
		log.Printf("stderr: %s", stderr.String())
		return err
	}

	fmt.Println("Download completed successfully")
	return nil
}
