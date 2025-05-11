# syntax=docker/dockerfile:1

# 1. Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git for module downloads
RUN apk add --no-cache git

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o kvant-backend cmd/main.go

# 2. Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
# Copy built binary from builder
COPY --from=builder /app/kvant-backend .

# Expose the application port
EXPOSE 8080

# Launch the application
CMD ["./kvant-backend"]

