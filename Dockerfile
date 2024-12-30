FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bot .
FROM ubuntu:20.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    libc6 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/bot .
COPY .env .
COPY whitelist.json .
RUN chmod +x /app/bot
CMD ["./bot"]
