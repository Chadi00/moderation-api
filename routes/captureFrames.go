package routes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Detect scene changes in the video and capture frame at every change
func captureFrames(videoPath string) (int, error) {
	// Define the output path for frames
	framesPath := filepath.Join("/app/downloads", "frames")
	fmt.Println("Creating frames directory...")

	// Ensure the frames directory exists
	if err := os.MkdirAll(framesPath, 0755); err != nil {
		return 0, fmt.Errorf("failed to create frames directory: %v", err)
	}

	fmt.Println("Running FFmpeg command...")
	// Define the FFmpeg command to run
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", "select='gt(scene,0.4)'",
		"-vsync", "vfr",
		"-frame_pts", "true",
		filepath.Join(framesPath, "frame_%04d.jpg"),
	)

	// Run the FFmpeg command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	// Count the number of frames extracted
	files, err := os.ReadDir(framesPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read frames directory: %v", err)
	}

	fmt.Printf("Processing %d frames...\n", len(files))
	for _, file := range files {
		framePath := filepath.Join(framesPath, file.Name())
		go analyzeFrame(framePath) // Call analyzeFrame in a goroutine for concurrency
	}

	// Return the number of frames processed
	return len(files), nil
}
