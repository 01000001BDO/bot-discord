FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# No need to build, we will run directly using 'go run'
FROM golang:1.22

WORKDIR /root/

# Copy everything including Go code
COPY --from=builder /app .

# Ensure go run is available and execute the application
CMD ["go", "run", "."]
