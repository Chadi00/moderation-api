package routes

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func transcriptAudio(audioPath string) (string, error) {
	fmt.Println("Starting transcription process...")
	url := "https://api.openai.com/v1/audio/transcriptions"

	fmt.Printf("Opening audio file at: %s\n", audioPath)
	file, err := os.Open(audioPath)
	if err != nil {
		fmt.Println("Error: failed to open the audio file.")
		return "", fmt.Errorf("failed to open the audio file: %v", err)
	}
	defer file.Close()

	fmt.Println("Creating multipart form data buffer...")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fmt.Println("Adding audio file to form...")
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", fmt.Errorf("cannot create form file: %v", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("cannot copy file content: %v", err)
	}

	fmt.Println("Adding model specification to form...")
	if err := writer.WriteField("model", "whisper-1"); err != nil {
		return "", fmt.Errorf("cannot write model field: %v", err)
	}

	fmt.Println("Closing multipart writer...")
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("cannot close writer: %v", err)
	}

	fmt.Println("Creating HTTP POST request...")
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %v", err)
	}

	fmt.Println("Loading API key from environment...")
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error: failed to load .env file.")
		return "", fmt.Errorf("failed to load .env: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: API key is not set in environment.")
		return "", fmt.Errorf("API key is not set")
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	fmt.Println("Sending HTTP request to OpenAI...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Reading response from server...")
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Println("Deleting audio file after processing...")
	err = deleteAudio(audioPath)
	if err != nil {
		fmt.Printf("Warning: failed to delete audio file: %v\n", err)
		return "", err
	}

	fmt.Println("Transcription process completed successfully.")
	return string(respBody), nil
}

func deleteAudio(audioPath string) error {
	fmt.Printf("Attempting to delete audio file at: %s\n", audioPath)
	err := os.Remove(audioPath)
	if err != nil {
		return fmt.Errorf("failed to delete audio file: %v", err)
	}
	fmt.Println("Audio file deleted successfully.")
	return nil
}
