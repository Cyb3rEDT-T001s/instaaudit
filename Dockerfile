# InstaAudit Docker Image
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o instaaudit cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates nmap curl openssl

# Create non-root user
RUN adduser -D -s /bin/sh instaaudit

# Set working directory
WORKDIR /home/instaaudit

# Copy binary from builder
COPY --from=builder /app/instaaudit .

# Copy documentation
COPY --from=builder /app/*.md ./docs/

# Change ownership
RUN chown -R instaaudit:instaaudit /home/instaaudit

# Switch to non-root user
USER instaaudit

# Expose common ports for web interface (if added later)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./instaaudit"]

# Default command
CMD ["--help"]