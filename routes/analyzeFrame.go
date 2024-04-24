package routes

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/joho/godotenv"
)

// Analyze frame with AWS Rekognition for moderation
func analyzeFrame(framePath string) {
	fmt.Printf("Loading environment variables for frame: %s\n", framePath)
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Loading AWS configuration...")
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	fmt.Printf("Creating Rekognition client for frame: %s\n", framePath)
	// Create a Rekognition client using the loaded configuration
	svc := rekognition.NewFromConfig(cfg)

	fmt.Printf("Reading image file from: %s\n", framePath)
	// Read the image from the filesystem
	imageBytes, err := os.ReadFile(framePath)
	if err != nil {
		log.Fatalf("failed to read image file, %v", err)
	}

	// Define the input parameters for the DetectModerationLabels operation
	input := &rekognition.DetectModerationLabelsInput{
		Image: &types.Image{
			Bytes: imageBytes,
		},
		MinConfidence: aws.Float32(75.0),
	}

	fmt.Println("Detecting moderation labels...")
	// Call the DetectModerationLabels operation
	resp, err := svc.DetectModerationLabels(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to detect moderation labels, %v", err)
	}

	// After calling the DetectModerationLabels operation
	if err != nil {
		log.Fatalf("failed to detect moderation labels, %v", err)
	} else {
		if len(resp.ModerationLabels) == 0 {
			fmt.Printf("No moderation labels detected for frame: %s\n", framePath)
		} else {
			fmt.Printf("Moderation labels detected for frame: %s\n", framePath)
			for _, label := range resp.ModerationLabels {
				fmt.Printf("Label: %s, Confidence: %.2f%%\n", *label.Name, *label.Confidence)
			}
		}
	}
}
