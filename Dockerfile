# Use the official Golang image to build the project
FROM golang:1.22.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o bot .

# Use a minimal base image for running the application
FROM debian:bullseye-slim

# Set the working directory in the runtime container
WORKDIR /app
COPY --from=builder /app/bot /app/bot
EXPOSE 8080
CMD ["./bot"]
