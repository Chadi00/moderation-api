package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ChatMessage represents a single message in the conversation with the chat model
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request payload for the chat API
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse represents the JSON structure of the response from OpenAI's chat API
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func getVideoDescription(framesDescription string, transcript string) (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", fmt.Errorf("failed to load .env")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	url := "https://api.openai.com/v1/chat/completions"

	// Setup request body
	requestBody1 := ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a video moderation assistant, the user share with you moderation labels on the frames of the video (that can be empty) and transcription of the audio of the video. Your goal is to respond with a rating from 1 to 4 (1 being family friendly, 2 being content that contains some profanity or a few violence, 3 being content that should be age restricted and 4 being explicit content that should not be published like porn or gore videos) and a moderation description of the video, saying if the video is a fit for anyone or if there is any explicit or implicit content or audio. You should describe why a video is family friendly or why it should not be published on public websites. It is crucial that you response begin with the rating then the video description separated by a pipe (|) for example '2|video description'"},
			{Role: "user", Content: "Moderation labels : " + framesDescription + "\n Audio transcript : " + transcript},
		},
	}

	jsonData1, err := json.Marshal(requestBody1)
	if err != nil {
		return "", "", fmt.Errorf("failed to marsher request number 1")
	}

	req1, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData1))
	if err != nil {
		return "", "", fmt.Errorf("failed to create new request number 1")
	}

	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+apiKey)

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		return "", "", fmt.Errorf("failed to send the request")
	}
	defer resp1.Body.Close()

	var chatResponse1 ChatResponse
	err = json.NewDecoder(resp1.Body).Decode(&chatResponse1)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse the response: %v", err)
	}

	// content from the first response
	var content1 string
	if len(chatResponse1.Choices) > 0 {
		content1 = chatResponse1.Choices[0].Message.Content
	}

	fmt.Println("rating : ", content1)

	return string(content1[0]), string(content1[2:]), nil

}
