# Use the official Golang image as a builder
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o my-go-app .

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/my-go-app .

# Command to run the executable
ENTRYPOINT ["/my-go-app"]

# Expose port 8080 to the outside world
EXPOSE 8080