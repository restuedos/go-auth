# Use official Go image as base
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install required dependencies
RUN apk add --no-cache make git

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Expose the application port
EXPOSE 8080

# Build the project using Makefile
RUN make build

# Start the app using Air (you can optionally change this to another make command)
CMD ["make", "run"]