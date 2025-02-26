# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o summarizer-app ./cmd/api

# Runtime stage
FROM alpine:latest

# Install SQLite and CA certificates
RUN apk add --no-cache sqlite ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/summarizer-app .

# Copy the database directory
COPY --from=builder /app/db ./db

# Copy the .env file
COPY .env ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./summarizer-app"]
