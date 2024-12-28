# Step 1: Build stage
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application (bot)
RUN go build -o bot .

# Step 2: Final image using Ubuntu
FROM ubuntu:20.04

# Install necessary dependencies for the final image
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/bot .

# Copy .env and whitelist.json to the container
COPY .env .
COPY whitelist.json .

# Ensure the binary is executable
RUN chmod +x /app/bot

# Set the environment variable for the bot token (optional if passed during runtime)
# ENV TOKEN=your-discord-bot-token

# Run the bot
CMD ["./bot"]
