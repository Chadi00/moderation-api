package routes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Detect scene changes in the video and capture frame at every change
func captureFrames(videoPath string) (string, error) {
	framesPath := filepath.Join("/app/downloads", "frames")
	fmt.Println("Creating frames directory...")
	if err := os.MkdirAll(framesPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create frames directory: %v", err)
	}

	fmt.Println("Running FFmpeg command...")
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", "select='eq(n\\,0)+gt(scene\\,0.3)',setpts=N/(FRAME_RATE*TB)",
		"-vsync", "vfr",
		"-frame_pts", "true",
		filepath.Join(framesPath, "frame_%04d.jpg"),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	files, err := os.ReadDir(framesPath)
	if err != nil {
		return "", fmt.Errorf("failed to read frames directory: %v", err)
	}

	descriptionChannel := make(chan string, len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))

	fmt.Printf("Processing %d frames...\n", len(files))
	for _, file := range files {
		framePath := filepath.Join(framesPath, file.Name())
		go func(path string) {
			defer wg.Done()
			frameDescription, err := analyzeFrame(path)
			if err != nil {
				fmt.Printf("Failed to analyze frame %s: %v\n", path, err)
				descriptionChannel <- ""
				return
			}
			descriptionChannel <- frameDescription
			if err := deleteFrame(path); err != nil {
				fmt.Printf("Failed to delete frame %s: %v\n", path, err)
			}
		}(framePath)
	}

	go func() {
		wg.Wait()
		close(descriptionChannel)
	}()

	var videoDescription string
	for description := range descriptionChannel {
		videoDescription += description + " "
	}

	return videoDescription, nil
}

func deleteFrame(framePath string) error {
	err := os.Remove(framePath)
	if err != nil {
		return fmt.Errorf("failed to delete frame")
	} else {
		fmt.Println("Deleted frame")
		return nil
	}
}
