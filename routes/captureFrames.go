package routes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func captureFrames(videoPath string) (int, error) {
	// Define the output path for frames
	framesPath := filepath.Join("/app/downloads", "frames")

	// Ensure the frames directory exists
	if err := os.MkdirAll(framesPath, 0755); err != nil {
		return 0, fmt.Errorf("failed to create frames directory: %v", err)
	}

	// Define the FFmpeg command to run
	// Using the select filter to detect scene changes with a threshold of 0.4
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", "select='gt(scene,0.4)'",
		"-vsync", "vfr",
		"-frame_pts", "true",
		filepath.Join(framesPath, "frame_%04d.jpg"),
	)

	// Run the FFmpeg command
	if output, err := cmd.CombinedOutput(); err != nil {
		return 0, fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	// Count the number of frames extracted
	files, err := os.ReadDir(framesPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read frames directory: %v", err)
	}

	// Return the number of frames
	return len(files), nil
}
