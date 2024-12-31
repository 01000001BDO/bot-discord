FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev pkgconfig opus-dev make git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o bot .
FROM ubuntu:20.04
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    ffmpeg \
    espeak \
    iputils-ping \
    libopus0 \
    libopus-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/bot .
COPY .env .
COPY whitelist.json .
RUN chmod +x /app/bot
CMD ["./bot"]