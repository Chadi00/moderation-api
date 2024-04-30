# moderation-api

## Overview
The `moderation-api` is a robust video content moderation tool designed to assist social media platforms and similar services in complying with terms of service by analyzing video content and providing moderation descriptions and ratings. The API uses a combination of Golang (Gin), Python (Pytube), Docker, Redis, and AWS hosting, integrated with AI/ML technologies like AWS Rekognition for frame analysis, Whisper for audio transcription, and GPT-3.5 for moderation descriptions and ratings.

## Features
- **Video Analysis**: Accepts a YouTube video URL and returns a moderation rating and descriptive analysis.
- **Rating Scale**: Provides a rating from 1 (family-friendly) to 4 (should not be published), helping platforms identify content suitability.
- **Technology Stack**: Utilizes Docker for easy deployment, Redis for efficient caching, and AWS services for robust hosting and AI analysis.

## Target Audience
This API is ideal for social media platforms and any other service needing robust video content moderation to ensure that explicit or inappropriate content is not published.

## Getting Started
### Prerequisites
- Docker
- An AWS account
- An OpenAI API key

### Setup
1. Clone the repository:
   ```
   git clone https://github.com/yourgithub/moderation-api.git
   ```
2. Navigate to the project directory and create a .env file with your AWS and OpenAI credentials:
  ```
  AWS_ACCESS_KEY_ID=your_access_key
  AWS_SECRET_ACCESS_KEY=your_secret_key
  OPENAI_API_KEY=your_openai_key
  ```
3. Build and run the Docker container:
   ```
   docker-compose up --build
   ```

### Usage
To analyze a video, make a POST request to the /analyze-video endpoint with the video URL:

```
Copy code
POST /analyze-video
Content-Type: application/json

{
  "VideoURL": "https://youtube.com/watch?v=example"
}
```
Expected response:
```
{
  "description": "This video is family-friendly. The audio transcript contains lyrics about going to a party, dancing, and having a good time. There is no explicit language or mature content, making it suitable for all audiences.",
  "rating": "1"
}
```
### Error Handling
Errors are returned directly in the response body with appropriate HTTP status codes, providing clear information about the issue encountered.

### License
This project is free to use. Feel free to fork and adapt for your needs.
