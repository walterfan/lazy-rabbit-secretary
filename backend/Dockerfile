# Multi-architecture Dockerfile for pre-built binary deployment
# This expects the binary to be built locally and copied to the container
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    wget \
    sqlite

# Set timezone
ENV TZ=Asia/Shanghai

WORKDIR /root/

# Copy the pre-built static binary based on target architecture
# Build command: make build-static
# The script will determine the correct binary to copy
ARG TARGETARCH
ARG TARGETPLATFORM

# Copy the appropriate binary based on architecture
COPY build/lazy-rabbit-secretary-linux-${TARGETARCH}-static ./lazy-rabbit-secretary

# Copy config files
COPY config ./config

# Copy certs directory (for JWT keys)
COPY certs ./certs

# Create directories for data and logs
RUN mkdir -p data logs

# Make binary executable
RUN chmod +x lazy-rabbit-secretary

# Expose port (default application port)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Command to run the server
CMD ["./lazy-rabbit-secretary", "server"]
