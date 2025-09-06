# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

ARG GOPROXY
ENV GOPROXY=$GOPROXY
ENV GO111MODULE=on
ENV CGO_ENABLED=1

# Install build dependencies for CGO (needed for SQLite and other database drivers)
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    sqlite-dev \
    postgresql-dev \
    pkgconfig \
    ca-certificates

# Update CA certificates
RUN update-ca-certificates

# Copy go mod files and vendor directory
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy source code
COPY . .

# Build the application with CGO enabled for database drivers
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o lazy-rabbit-reminder .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    wget

# Set timezone
ENV TZ=Asia/Shanghai

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/lazy-rabbit-reminder .

# Copy config files
COPY --from=builder /app/config ./config

# Copy certs directory (for JWT keys)
COPY --from=builder /app/certs ./certs

# Create directories for data and logs
RUN mkdir -p data logs

# Expose port (default application port)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Command to run the server
CMD ["./lazy-rabbit-reminder", "server"]
