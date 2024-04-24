package routes

import (
	"fmt"
	"os"
)

// Delete video from the video path
func deleteVideo(videoPath string) error {
	// Attempt to delete the video file
	err := os.Remove(videoPath)
	if err != nil {
		// If there was an error deleting the file, return the error
		return fmt.Errorf("failed to delete video file %s: %v", videoPath, err)
	}

	// If successful, log the deletion and return nil indicating no error
	fmt.Printf("Successfully deleted video file: %s\n", videoPath)
	return nil
}
