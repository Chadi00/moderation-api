package routes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// extracts the audio from the video file and saves it as an MP3 file, return the mp3 file path.
func extractAudio(videoPath string) (string, error) {

	audioPath := filepath.Join("/app/downloads", "audio")

	// Ensure the audio directory exists
	if err := os.MkdirAll(audioPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create audio directory: %v", err)
	}

	// Construct the path for the output MP3 file
	outputMP3 := filepath.Join(audioPath, "output_audio.mp3")

	// Define the FFmpeg command to extract audio
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vn",                   // No video.
		"-acodec", "libmp3lame", // Use MP3 encoder
		"-ac", "2", // Set audio channels to 2
		"-q:a", "4", // Quality for MP3
		"-y", // Overwrite output files without asking
		outputMP3,
	)

	// Run the FFmpeg command
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg failed to extract audio: %v, output: %s", err, string(output))
	}

	// Return the path of the created MP3 file
	return outputMP3, nil
}
