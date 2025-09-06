# Docker Services Integration

This document describes the complete Docker services stack for the Lazy Rabbit Reminder application, including PostgreSQL with pgvector, NATS message queue, pgweb database admin, and PlantUML server.

## Overview

NATS is a lightweight, high-performance messaging system that serves as the message queue for asynchronous operations in the application.

## Docker Compose Configuration

The NATS service is configured in `docker-compose.yml` with the following features:

### Service Configuration
```yaml
nats:
  container_name: lazy_rabbit_nats
  image: nats:2.10.20-alpine
  restart: always
  ports:
    - "4222:4222"   # Client connections
    - "8222:8222"   # HTTP monitoring port  
    - "6222:6222"   # Cluster routing port
```

### Features Enabled
- **JetStream**: Persistent messaging and stream processing
- **HTTP Monitoring**: Web-based monitoring on port 8222
- **Data Persistence**: Messages stored in `/data` volume
- **Health Checks**: Automated health monitoring

### Port Configuration
| Port | Purpose | Description |
|------|---------|-------------|
| 4222 | Client Connections | Main NATS protocol port |
| 8222 | HTTP Monitoring | Web UI and metrics endpoint |
| 6222 | Cluster Routing | Inter-cluster communication |

## Environment Variables

Configure NATS connection in your `.env` file:

```bash
# NATS Configuration
NATS_URL=nats://localhost:4222
NATS_CLUSTER_ID=lazy-rabbit-cluster
NATS_CLIENT_ID=lazy-rabbit-client
NATS_MONITORING_PORT=8222
NATS_MAX_PAYLOAD=1048576
NATS_MAX_CONNECTIONS=1000
```

## Usage Examples

### Starting the Services
```bash
# Start all services including NATS
docker-compose up -d

# Start only NATS
docker-compose up -d nats
```

### Accessing NATS Monitoring
Open your browser and navigate to:
- **Local Development**: http://localhost:8222
- **Docker Environment**: http://localhost:8222

### Checking NATS Status
```bash
# Check NATS container logs
docker-compose logs nats

# Check NATS health
curl http://localhost:8222/healthz

# View NATS server info
curl http://localhost:8222/varz
```

## Integration Points

### Potential Use Cases in Lazy Rabbit Reminder

1. **Email Queue**: Asynchronous email sending
   - Registration confirmations
   - Approval notifications
   - Password reset emails

2. **Task Processing**: Background job processing
   - Scheduled reminders
   - Data cleanup tasks
   - Report generation

3. **Event Streaming**: Real-time event processing
   - User activity tracking
   - Audit logging
   - System monitoring

4. **Inter-Service Communication**: Microservice messaging
   - Service-to-service communication
   - Event-driven architecture
   - Decoupled system components

### Example Go Integration

```go
package main

import (
    "log"
    "github.com/nats-io/nats.go"
)

func main() {
    // Connect to NATS
    nc, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    // Publish a message
    nc.Publish("email.send", []byte("Hello NATS!"))

    // Subscribe to messages
    nc.Subscribe("email.send", func(m *nats.Msg) {
        log.Printf("Received: %s", string(m.Data))
    })
}
```

## Monitoring and Management

### NATS Monitoring Dashboard
The NATS server provides a built-in monitoring interface accessible at `http://localhost:8222` with the following endpoints:

- `/` - General server information
- `/connz` - Connection information
- `/routez` - Route information
- `/subsz` - Subscription information
- `/varz` - Server variables and statistics

### Health Checks
The Docker configuration includes automatic health checks:
```yaml
healthcheck:
  test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8222/healthz"]
  interval: 10s
  timeout: 5s
  retries: 3
  start_period: 10s
```

## Data Persistence

NATS JetStream data is persisted in a Docker volume:
```yaml
volumes:
  - nats_data:/data
```

This ensures message durability across container restarts.

## Security Considerations

For production environments, consider:

1. **Authentication**: Enable NATS authentication
2. **TLS/SSL**: Configure encrypted connections
3. **Network Security**: Restrict access to NATS ports
4. **Resource Limits**: Set appropriate memory and CPU limits

## Troubleshooting

### Common Issues

1. **Connection Refused**
   ```bash
   # Check if NATS is running
   docker-compose ps nats
   
   # Check NATS logs
   docker-compose logs nats
   ```

2. **Port Conflicts**
   ```bash
   # Check if ports are already in use
   netstat -tulpn | grep :4222
   netstat -tulpn | grep :8222
   ```

3. **Memory Issues**
   ```bash
   # Check NATS memory usage
   curl http://localhost:8222/varz | jq '.mem'
   ```

### Useful Commands

```bash
# Restart NATS service
docker-compose restart nats

# View NATS configuration
docker-compose exec nats cat /etc/nats/nats.conf

# Check NATS version
docker-compose exec nats nats-server --version

# Clean up NATS data (WARNING: This will delete all messages)
docker-compose down
docker volume rm deploy_nats_data
```

## Next Steps

To fully integrate NATS into the application:

1. Add NATS client library to `go.mod`
2. Create NATS connection manager
3. Implement message publishers and subscribers
4. Add NATS configuration to application config
5. Create message handlers for different event types

## References

- [NATS Official Documentation](https://docs.nats.io/)
- [NATS Docker Hub](https://hub.docker.com/_/nats)
- [JetStream Documentation](https://docs.nats.io/jetstream)
- [NATS Go Client](https://github.com/nats-io/nats.go)
