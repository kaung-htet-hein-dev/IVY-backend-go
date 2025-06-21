# Build stage
FROM golang:1.24.1-alpine AS build

# Install dependencies for building
RUN apk --no-cache add git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates and create non-root user
RUN apk --no-cache add ca-certificates && \
    addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/main .

# Copy environment file (create this file)
COPY .env.production .

# Make binary executable and change ownership
RUN chmod +x main && chown appuser:appgroup main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check (optional - remove if no health endpoint)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application with prod command
CMD ["./main", "prod"]
