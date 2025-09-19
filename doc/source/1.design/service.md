# Docker Services Stack

This document describes the complete Docker services stack for the Lazy Rabbit Secretary application.

## Architecture Overview

The application consists of the following services orchestrated via Docker Compose:

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose Stack                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Nginx     │  │     Web     │  │    NATS     │         │
│  │   :1980     │  │   :8000     │  │   :4222     │         │
│  │   :1981     │  │             │  │   :8222     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                          │                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  PostgreSQL │  │    Redis    │  │   pgweb     │         │
│  │  (pgvector) │  │   :6379     │  │   :8081     │         │
│  │   :5432     │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                             │
│  ┌─────────────┐                                           │
│  │  PlantUML   │                                           │
│  │   :8002     │                                           │
│  └─────────────┘                                           │
│                                                             │
│  Network: pgvector-net                                      │
└─────────────────────────────────────────────────────────────┘
```

## Services Configuration

### 1. PostgreSQL with pgvector

**Container**: `pgvector`  
**Image**: `ankane/pgvector`  
**Port**: `5432`  
**Purpose**: Primary database with vector extension support

```yaml
pgvector:
  image: ankane/pgvector
  container_name: pgvector
  environment:
    POSTGRES_USER: ${DB_USER}
    POSTGRES_PASSWORD: ${DB_PASS}
    POSTGRES_DB: ${DB_NAME}
  ports:
    - "5432:5432"
  volumes:
    - pgvector-data:/var/lib/postgresql/data
  networks:
    - pgvector-net
  restart: unless-stopped
```

**Features**:
- Full PostgreSQL database functionality
- pgvector extension for vector operations
- Persistent data storage
- Optimized for AI/ML workloads

### 2. Redis Cache

**Container**: `redis`  
**Image**: `redis:7-alpine`  
**Port**: `6379`  
**Purpose**: Caching and session storage

```yaml
redis:
  image: redis:7-alpine
  container_name: redis
  ports:
    - "6379:6379"
  volumes:
    - redis-data:/data
  networks:
    - pgvector-net
  restart: unless-stopped
  command: redis-server --appendonly yes
```

**Features**:
- AOF (Append Only File) persistence enabled
- High-performance key-value store
- Session management
- Caching layer

### 3. NATS Message Queue

**Container**: `lazy_rabbit_nats`  
**Image**: `nats:2.10.20-alpine`  
**Ports**: `4222` (client), `8222` (monitoring), `6222` (cluster)  
**Purpose**: Message queue and event streaming

```yaml
nats:
  container_name: lazy_rabbit_nats
  image: nats:2.10.20-alpine
  restart: always
  ports:
    - "4222:4222"   # Client connections
    - "8222:8222"   # HTTP monitoring port
    - "6222:6222"   # Cluster routing port
  command: 
    - "--jetstream"
    - "--store_dir=/data"
    - "--http_port=8222"
    - "--max_payload=1MB"
    - "--max_connections=1000"
  volumes:
    - nats_data:/data
  networks:
    - pgvector-net
  environment:
    - TZ=Asia/Shanghai
  healthcheck:
    test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8222/healthz"]
    interval: 10s
    timeout: 5s
    retries: 3
    start_period: 10s
```

**Features**:
- JetStream for persistent messaging
- Built-in monitoring dashboard
- High-performance messaging
- Event-driven architecture support

### 4. pgweb Database Admin

**Container**: `pgweb`  
**Image**: `sosedoff/pgweb`  
**Port**: `8081`  
**Purpose**: Web-based PostgreSQL administration

```yaml
pgweb:
  image: sosedoff/pgweb
  container_name: pgweb
  expose:
    - "8081"
  ports:
    - "8081:8081"
  environment:
    - PGWEB_DATABASE_URL=${DB_URL}
  depends_on:
    - pgvector
  networks:
    - pgvector-net
  restart: unless-stopped
```

**Features**:
- Web-based database administration
- SQL query execution
- Database schema visualization
- Table data browsing and editing

### 5. PlantUML Server

**Container**: `plantuml-server`  
**Image**: `plantuml/plantuml-server:jetty`  
**Port**: `8002`  
**Purpose**: Diagram generation service

```yaml
plantuml:
  image: plantuml/plantuml-server:jetty
  container_name: plantuml-server
  expose:
    - "8080"
  ports:
    - "8002:8080"
  networks:
    - pgvector-net
  restart: unless-stopped
```

**Features**:
- UML diagram generation
- REST API for diagram creation
- Multiple output formats (PNG, SVG, etc.)
- Integration with documentation systems

### 6. Web Application

**Container**: `lazy_rabbit_reminder`  
**Build**: `../backend`  
**Port**: `8000`  
**Purpose**: Main application backend

```yaml
web:
  container_name: lazy_rabbit_reminder
  build: ../backend
  ports:
    - "8000:8080"
  environment:
    - TZ=Asia/Shanghai
  depends_on:
    - pgvector
    - redis
    - nats
  networks:
    - pgvector-net
```

### 7. Nginx Reverse Proxy

**Container**: `lazy_rabbit_nginx`  
**Image**: `nginx:1.13-alpine`  
**Ports**: `1980` (HTTP), `1981` (HTTPS)  
**Purpose**: Reverse proxy and static file serving

## Environment Variables

Configure the following environment variables in your `.env` file:

```bash
# PostgreSQL Database Configuration (with pgvector support)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_PASS=password
DB_NAME=lazy_rabbit_reminder
DB_SSLMODE=disable
DB_URL=postgres://postgres:password@pgvector:5432/lazy_rabbit_reminder?sslmode=disable

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_PORT=6379

# NATS Configuration
NATS_URL=nats://localhost:4222
NATS_CLUSTER_ID=lazy-rabbit-cluster
NATS_CLIENT_ID=lazy-rabbit-client
NATS_MONITORING_PORT=8222
NATS_MAX_PAYLOAD=1048576
NATS_MAX_CONNECTIONS=1000

# pgweb Configuration
PGWEB_DATABASE_URL=postgres://postgres:password@pgvector:5432/lazy_rabbit_reminder?sslmode=disable

# PlantUML Server Configuration
PLANTUML_SERVER_URL=http://localhost:8002
```

## Service Access Points

| Service | URL | Purpose |
|---------|-----|---------|
| **Web App** | http://localhost:8000 | Main application |
| **Frontend** | http://localhost:1980 | Nginx-served frontend |
| **PostgreSQL** | localhost:5432 | Database connection |
| **Redis** | localhost:6379 | Cache connection |
| **NATS Client** | nats://localhost:4222 | Message queue |
| **NATS Monitor** | http://localhost:8222 | NATS dashboard |
| **pgweb** | http://localhost:8081 | Database admin |
| **PlantUML** | http://localhost:8002 | Diagram service |

## Network Configuration

All services communicate through a custom bridge network:

```yaml
networks:
  pgvector-net:
    driver: bridge
```

**Benefits**:
- Service isolation
- Internal DNS resolution
- Improved security
- Container-to-container communication

## Persistent Storage

```yaml
volumes:
  pgvector-data:
    driver: local      # PostgreSQL data
  redis-data:
    driver: local      # Redis persistence
  nats_data:
    driver: local      # NATS JetStream data
```

## Usage Commands

### Starting Services

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d pgvector

# View logs
docker-compose logs -f [service_name]

# Check service status
docker-compose ps
```

### Database Operations

```bash
# Connect to PostgreSQL
docker-compose exec pgvector psql -U postgres -d lazy_rabbit_reminder

# Access pgweb admin
open http://localhost:8081

# Redis CLI
docker-compose exec redis redis-cli
```

### NATS Operations

```bash
# View NATS monitoring
open http://localhost:8222

# Check NATS health
curl http://localhost:8222/healthz

# View server stats
curl http://localhost:8222/varz
```

### PlantUML Usage

```bash
# Test PlantUML server
curl -X POST http://localhost:8002/plantuml/png \
  -H "Content-Type: text/plain" \
  -d "@startuml
      Alice -> Bob: Hello
      @enduml"
```

## Integration Examples

### PostgreSQL with pgvector

```go
import (
    "github.com/pgvector/pgvector-go"
    "gorm.io/gorm"
)

type Document struct {
    ID       uint
    Content  string
    Embedding pgvector.Vector `gorm:"type:vector(1536)"`
}

// Create vector index
db.Exec("CREATE INDEX ON documents USING ivfflat (embedding vector_cosine_ops)")
```

### NATS Messaging

```go
import "github.com/nats-io/nats.go"

// Connect to NATS
nc, _ := nats.Connect("nats://localhost:4222")

// Publish message
nc.Publish("email.send", []byte("Hello NATS!"))

// Subscribe to messages
nc.Subscribe("email.send", func(m *nats.Msg) {
    log.Printf("Received: %s", string(m.Data))
})
```

### Redis Caching

```go
import "github.com/go-redis/redis/v8"

rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Set cache
rdb.Set(ctx, "key", "value", time.Hour)

// Get cache
val := rdb.Get(ctx, "key")
```

## Monitoring and Health Checks

### Service Health

```bash
# Check all services
docker-compose ps

# Individual service health
docker-compose exec pgvector pg_isready
docker-compose exec redis redis-cli ping
curl http://localhost:8222/healthz  # NATS
curl http://localhost:8081/         # pgweb
curl http://localhost:8002/         # PlantUML
```

### Resource Monitoring

```bash
# View resource usage
docker stats

# Service-specific stats
docker stats pgvector redis lazy_rabbit_nats
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check port usage
   netstat -tulpn | grep :5432
   netstat -tulpn | grep :6379
   netstat -tulpn | grep :4222
   ```

2. **Network Issues**
   ```bash
   # Inspect network
   docker network inspect deploy_pgvector-net
   
   # Test connectivity
   docker-compose exec web ping pgvector
   ```

3. **Volume Issues**
   ```bash
   # List volumes
   docker volume ls
   
   # Inspect volume
   docker volume inspect deploy_pgvector-data
   ```

### Service Restart

```bash
# Restart specific service
docker-compose restart pgvector

# Recreate service
docker-compose up -d --force-recreate pgvector

# Clean restart all services
docker-compose down && docker-compose up -d
```

## Security Considerations

### Production Deployment

1. **Database Security**
   - Use strong passwords
   - Enable SSL/TLS connections
   - Restrict network access
   - Regular backups

2. **Network Security**
   - Use internal networks only
   - Implement firewall rules
   - Enable container security scanning

3. **Secrets Management**
   - Use Docker secrets or external secret management
   - Avoid plain text passwords in environment files
   - Rotate credentials regularly

## Performance Optimization

### Database Tuning

```sql
-- PostgreSQL optimization
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
```

### Redis Optimization

```bash
# Redis memory optimization
redis-cli CONFIG SET maxmemory 256mb
redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

### NATS Tuning

```bash
# NATS server optimization (in docker-compose)
command:
  - "--max_payload=2MB"
  - "--max_connections=2000"
  - "--max_subscriptions=1000"
```

## Migration from MariaDB

If migrating from the previous MariaDB setup:

1. **Export existing data**
   ```bash
   docker-compose exec db mysqldump -u root -p lazy_rabbit_reminder > backup.sql
   ```

2. **Convert to PostgreSQL format**
   ```bash
   # Use tools like pgloader or manual conversion
   pgloader mysql://user:pass@localhost/lazy_rabbit_reminder postgresql://user:pass@localhost/lazy_rabbit_reminder
   ```

3. **Update application configuration**
   - Change database driver from MySQL to PostgreSQL
   - Update connection strings
   - Test all database operations

## Next Steps

1. **Vector Operations**: Implement pgvector for AI/ML features
2. **Message Queues**: Integrate NATS for async operations
3. **Monitoring**: Add Prometheus/Grafana for metrics
4. **Backup Strategy**: Implement automated backups
5. **Load Balancing**: Scale services horizontally

## References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [pgvector Extension](https://github.com/pgvector/pgvector)
- [NATS Documentation](https://docs.nats.io/)
- [Redis Documentation](https://redis.io/documentation)
- [pgweb Documentation](https://github.com/sosedoff/pgweb)
- [PlantUML Server](https://github.com/plantuml/plantuml-server)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
