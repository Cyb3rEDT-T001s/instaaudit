# InstaAudit Docker Image
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN go mod tidy && \
    go build -o instaaudit cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/instaaudit .

# Make it executable
RUN chmod +x instaaudit

# Expose no ports (tool makes outbound connections only)

# Set entrypoint
ENTRYPOINT ["./instaaudit"]

# Default help command
CMD ["--help"]