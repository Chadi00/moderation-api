# Start with the official Golang base image for the build stage
FROM golang:1.22-alpine as builder

# Install build dependencies (Alpine uses 'apk' for package management)
RUN apk add --no-cache git gcc g++ libc-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o moderation-api .

# Start a new stage from scratch for a smaller final image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache python3 py3-pip ffmpeg python3-dev

# Setup a virtual environment and install pytube
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip install --no-cache-dir pytube

# Set the working directory in the container
WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/moderation-api .

# Copy the rest of the necessary files
COPY --from=builder /app/utils ./utils
COPY --from=builder /app/models ./models
COPY --from=builder /app/routes ./routes

# Ensure the downloads folder is available
RUN mkdir -p downloads

# Command to run the executable
CMD ["./moderation-api"]
