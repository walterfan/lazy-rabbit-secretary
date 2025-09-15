# Redis Mini Manual

A comprehensive guide to Redis setup, configuration, and usage in the Lazy Rabbit Reminder project.

## Table of Contents

1. [What is Redis?](#what-is-redis)
2. [Installation & Setup](#installation--setup)
3. [Authentication](#authentication)
4. [Basic Commands](#basic-commands)
5. [Data Types & Operations](#data-types--operations)
6. [Advanced Features](#advanced-features)
7. [Application Integration](#application-integration)
8. [Performance & Monitoring](#performance--monitoring)
9. [Security Recommendations](#security-recommendations)
10. [Troubleshooting](#troubleshooting)
11. [Quick Reference](#quick-reference)
12. [Summary](#summary)

## What is Redis?

Redis (Remote Dictionary Server) is an in-memory data structure store used as:
- **Database**: Fast key-value storage
- **Cache**: High-performance caching layer
- **Message Broker**: Pub/Sub messaging system
- **Session Store**: User session management

### Key Features
- ⚡ **In-Memory**: Extremely fast read/write operations
- 🔄 **Persistence**: Optional data durability (AOF, RDB)
- 🏗️ **Data Structures**: Strings, Lists, Sets, Hashes, Sorted Sets
- 📡 **Pub/Sub**: Real-time messaging
- 🔧 **Atomic Operations**: Thread-safe operations
- 📈 **Scalability**: Master-slave replication, clustering

## Installation & Setup

### Docker Setup (Current Project)

Our project uses Redis in Docker with simple password authentication.

**Files:**
- `deploy/redis/redis.conf` - Redis configuration
- `deploy/docker-compose.yml` - Container orchestration
- `deploy/env.example` - Environment variables template

### Environment Configuration

Create a `.env` file in the `deploy/` directory:

```bash
# Copy the example file
cp env.example .env

# Configure Redis credentials
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your-secure-password
```

### Docker Compose Configuration

The Redis service is configured with:

```yaml
redis:
  image: redis:7-alpine
  container_name: redis
  ports:
    - "6379:6379"
  volumes:
    - redis-data:/data
    - ./redis/redis.conf:/usr/local/etc/redis/redis.conf:ro
  command: |
    sh -c "
      redis-server /usr/local/etc/redis/redis.conf --requirepass $${REDIS_PASSWORD}
    "
  environment:
    - REDIS_PASSWORD=${REDIS_PASSWORD:-lazy-rabbit}
```

### Starting Redis

```bash
# Start all services
docker-compose up -d

# Start only Redis
docker-compose up -d redis

# Check Redis status
docker-compose ps redis
```

## Authentication

### Simple Password Authentication

Our setup uses `requirepass` with environment variables for security.

**How it works:**
1. Redis starts with `--requirepass $REDIS_PASSWORD` flag
2. Password is loaded from environment variable
3. No hardcoded credentials in config files

### Testing Authentication

```bash
# Connect to Redis container
docker exec -it redis redis-cli

# Try command without auth (will fail)
127.0.0.1:6379> PING
(error) NOAUTH Authentication required.

# Authenticate with password
127.0.0.1:6379> AUTH your-secure-password
OK

# Now commands work
127.0.0.1:6379> PING
PONG
```

## Basic Commands

### Connection & Info Commands

```bash
# Connect to Redis
redis-cli -h redis -p 6379 -a your-password

# Basic connectivity
PING                    # Test connection
ECHO "Hello Redis"      # Echo message
INFO                    # Server information
CONFIG GET "*"          # Get configuration
CLIENT LIST             # List connected clients
```

### Key Management

```bash
# Key operations
SET mykey "Hello World"     # Set a key
GET mykey                   # Get a key value
EXISTS mykey               # Check if key exists
DEL mykey                  # Delete a key
KEYS *                     # List all keys (avoid in production)
SCAN 0                     # Iterate keys safely
EXPIRE mykey 60            # Set expiration (60 seconds)
TTL mykey                  # Check time to live
PERSIST mykey              # Remove expiration
```

### Database Operations

```bash
# Database selection
SELECT 0                   # Select database 0 (default)
SELECT 1                   # Select database 1

# Database management
FLUSHDB                    # Clear current database
FLUSHALL                   # Clear all databases (dangerous!)
DBSIZE                     # Number of keys in current DB
```

## Data Types & Operations

### 1. Strings

Most basic Redis data type - binary safe strings up to 512MB.

```bash
# Basic string operations
SET user:1:name "John Doe"          # Set string
GET user:1:name                     # Get string
MSET key1 "val1" key2 "val2"       # Set multiple
MGET key1 key2                     # Get multiple

# String with expiration
SETEX session:abc123 3600 "user_data"  # Set with TTL
SETNX lock:resource1 "locked"           # Set if not exists

# Numeric operations
SET counter 0                       # Set number
INCR counter                        # Increment by 1
INCRBY counter 5                    # Increment by 5
DECR counter                        # Decrement by 1
DECRBY counter 3                    # Decrement by 3

# String manipulation
APPEND mykey " World"               # Append to string
STRLEN mykey                        # Get string length
GETRANGE mykey 0 4                 # Get substring
```

### 2. Lists

Ordered collections of strings (linked lists).

```bash
# List operations
LPUSH tasks "task1" "task2"         # Push to left (head)
RPUSH tasks "task3"                 # Push to right (tail)
LPOP tasks                          # Pop from left
RPOP tasks                          # Pop from right

# List access
LLEN tasks                          # List length
LINDEX tasks 0                      # Get by index
LRANGE tasks 0 -1                   # Get range (all)
LRANGE tasks 0 2                    # Get first 3 items

# List modification
LSET tasks 0 "new_task"            # Set by index
LINSERT tasks BEFORE "task2" "new"  # Insert before/after
LREM tasks 1 "task1"               # Remove elements
LTRIM tasks 0 9                     # Keep only range
```

### 3. Sets

Unordered collections of unique strings.

```bash
# Set operations
SADD tags "redis" "database" "cache"   # Add members
SREM tags "cache"                      # Remove member
SMEMBERS tags                          # Get all members
SCARD tags                             # Count members
SISMEMBER tags "redis"                 # Check membership

# Set operations between sets
SINTER set1 set2                       # Intersection
SUNION set1 set2                       # Union
SDIFF set1 set2                        # Difference
SINTERSTORE result set1 set2           # Store intersection
```

### 4. Hashes

Maps of field-value pairs (like objects/dictionaries).

```bash
# Hash operations
HSET user:1 name "John" age 30 email "john@example.com"
HGET user:1 name                       # Get field
HMGET user:1 name age                  # Get multiple fields
HGETALL user:1                         # Get all fields
HKEYS user:1                           # Get all field names
HVALS user:1                           # Get all values

# Hash modification
HINCRBY user:1 age 1                   # Increment numeric field
HDEL user:1 email                      # Delete field
HEXISTS user:1 name                    # Check field exists
HLEN user:1                            # Count fields
```

### 5. Sorted Sets

Sets with scores for ordering.

```bash
# Sorted set operations
ZADD leaderboard 100 "player1" 85 "player2" 120 "player3"
ZRANGE leaderboard 0 -1                # Get by rank (asc)
ZREVRANGE leaderboard 0 -1             # Get by rank (desc)
ZRANGE leaderboard 0 -1 WITHSCORES     # Include scores

# Score-based queries
ZRANGEBYSCORE leaderboard 90 200       # Get by score range
ZCOUNT leaderboard 90 200              # Count in score range
ZSCORE leaderboard "player1"           # Get member score
ZRANK leaderboard "player1"            # Get member rank

# Sorted set modification
ZINCRBY leaderboard 10 "player1"       # Increment score
ZREM leaderboard "player2"             # Remove member
ZCARD leaderboard                      # Count members
```

## Advanced Features

### Pub/Sub Messaging

Redis supports publish/subscribe messaging for real-time communication.

```bash
# Publisher (in one terminal)
PUBLISH notifications "New message arrived"
PUBLISH user:123:alerts "Your order is ready"

# Subscriber (in another terminal)
SUBSCRIBE notifications user:123:alerts    # Subscribe to channels
PSUBSCRIBE user:*:alerts                   # Pattern subscription
UNSUBSCRIBE notifications                  # Unsubscribe
```

### Transactions

Redis supports atomic transactions with MULTI/EXEC.

```bash
# Transaction example
MULTI                           # Start transaction
SET key1 "value1"
INCR counter
EXPIRE key1 60
EXEC                           # Execute all commands atomically

# Conditional transaction with WATCH
WATCH mykey                    # Watch key for changes
MULTI
SET mykey "new_value"
EXEC                          # Will fail if mykey changed
```

### Lua Scripting

Execute Lua scripts atomically on the server.

```bash
# Simple Lua script
EVAL "return redis.call('SET', KEYS[1], ARGV[1])" 1 mykey myvalue

# More complex script
EVAL "
local current = redis.call('GET', KEYS[1])
if current == false then
    return redis.call('SET', KEYS[1], ARGV[1])
else
    return current
end
" 1 mykey myvalue
```

### Persistence

Redis offers two persistence options:

**RDB (Redis Database Backup)**
- Point-in-time snapshots
- Compact single file
- Good for backups

**AOF (Append Only File)**
- Logs every write operation
- Better durability
- Larger file size

```bash
# Persistence commands
SAVE                           # Synchronous save (blocks)
BGSAVE                         # Background save
BGREWRITEAOF                   # Rewrite AOF file
LASTSAVE                       # Last successful save timestamp
```

## Application Integration

### Go Redis Client

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/go-redis/redis/v8"
)

func main() {
    // Create Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0, // Default database
    })
    
    ctx := context.Background()
    
    // Test connection
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Redis connection failed:", err)
    }
    fmt.Println("Redis connected:", pong)
    
    // Basic operations
    err = rdb.Set(ctx, "key", "value", time.Hour).Err()
    if err != nil {
        log.Fatal("Set failed:", err)
    }
    
    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        log.Fatal("Get failed:", err)
    }
    fmt.Println("Retrieved:", val)
    
    // Hash operations
    err = rdb.HSet(ctx, "user:1", map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    }).Err()
    
    user := rdb.HGetAll(ctx, "user:1").Val()
    fmt.Printf("User: %+v\n", user)
}
```

### Common Use Cases in Our Project

```go
// Session management
func StoreSession(rdb *redis.Client, sessionID, userID string) error {
    return rdb.Set(context.Background(), 
        "session:"+sessionID, userID, 24*time.Hour).Err()
}

// Cache user data
func CacheUser(rdb *redis.Client, userID string, userData map[string]interface{}) error {
    return rdb.HSet(context.Background(), 
        "user:"+userID, userData).Err()
}

// Rate limiting
func CheckRateLimit(rdb *redis.Client, userID string, limit int) (bool, error) {
    key := "rate_limit:" + userID
    count, err := rdb.Incr(context.Background(), key).Result()
    if err != nil {
        return false, err
    }
    
    if count == 1 {
        rdb.Expire(context.Background(), key, time.Minute)
    }
    
    return count <= int64(limit), nil
}

// Task queue
func EnqueueTask(rdb *redis.Client, queueName, taskData string) error {
    return rdb.LPush(context.Background(), 
        "queue:"+queueName, taskData).Err()
}

func DequeueTask(rdb *redis.Client, queueName string) (string, error) {
    return rdb.BRPop(context.Background(), 
        time.Second*10, "queue:"+queueName).Result()[1], nil
}
```

## Performance & Monitoring

### Performance Commands

```bash
# Performance monitoring
INFO stats                     # Performance statistics
INFO memory                    # Memory usage
INFO replication              # Replication info
INFO clients                  # Client connections
SLOWLOG GET 10                # Show slow queries
MONITOR                       # Real-time command monitoring (dev only)

# Memory analysis
MEMORY USAGE mykey            # Memory used by key
MEMORY STATS                  # Memory statistics
DEBUG OBJECT mykey            # Object debugging info

# Performance tuning
CONFIG SET maxmemory 256mb    # Set memory limit
CONFIG SET maxmemory-policy allkeys-lru  # Eviction policy
```

### Key Performance Metrics

**Memory Usage**
- Monitor with `INFO memory`
- Set appropriate `maxmemory` limit
- Choose right eviction policy

**Connection Count**
- Monitor with `INFO clients`
- Set `maxclients` if needed
- Use connection pooling

**Command Statistics**
- Monitor with `INFO commandstats`
- Identify slow operations
- Use `SLOWLOG` for debugging

### Best Practices

**Key Naming**
```bash
# Good: Hierarchical, descriptive
user:123:profile
session:abc123:data
cache:product:456

# Bad: Unclear, no structure
u123
data
temp
```

**Memory Optimization**
- Use appropriate data types
- Set expiration on temporary data
- Use compression for large values
- Monitor memory usage regularly

**Performance Tips**
- Use pipelining for multiple commands
- Prefer SCAN over KEYS in production
- Use Lua scripts for complex operations
- Set appropriate timeout values

## Security Recommendations

1. **Strong Authentication**: Use complex passwords (20+ characters)
2. **Environment Variables**: Never hardcode credentials
3. **Network Security**: Redis only accessible within Docker network
4. **Regular Updates**: Keep Redis image updated
5. **Password Rotation**: Change passwords regularly
6. **Disable Dangerous Commands**: Use `rename-command` in production
7. **Monitor Access**: Use `CONFIG SET protected-mode yes`

## Troubleshooting

### Common Issues

**1. Authentication Failed**
```bash
# Check password in .env file
cat .env | grep REDIS_PASSWORD

# Test with correct password
docker exec -it redis redis-cli -a your-password ping
```

**2. Connection Refused**
```bash
# Check if Redis is running
docker-compose ps redis

# Check Redis logs
docker-compose logs redis

# Restart Redis if needed
docker-compose restart redis
```

**3. Memory Issues**
```bash
# Check memory usage
docker exec -it redis redis-cli INFO memory

# Check if memory limit is reached
docker exec -it redis redis-cli CONFIG GET maxmemory
```

**4. Performance Issues**
```bash
# Check slow queries
docker exec -it redis redis-cli SLOWLOG GET 10

# Monitor commands in real-time (development only)
docker exec -it redis redis-cli MONITOR
```

### Debug Commands

```bash
# Container management
docker-compose ps redis                    # Check status
docker-compose logs redis                  # View logs
docker-compose restart redis               # Restart service

# Connection testing
docker exec -it redis redis-cli ping      # Test from inside container
docker exec -it redis redis-cli -a $REDIS_PASSWORD ping  # With auth

# Performance debugging
docker exec -it redis redis-cli INFO all  # Complete server info
docker exec -it redis redis-cli CLIENT LIST  # Active connections
docker exec -it redis redis-cli CONFIG GET "*"  # All configuration
```

### Emergency Recovery

**Clear all data (use with caution)**
```bash
docker exec -it redis redis-cli FLUSHALL
```

**Reset Redis completely**
```bash
docker-compose down
docker volume rm lazy-rabbit-reminder_redis-data
docker-compose up -d redis
```

## Quick Reference

### Essential Commands Cheat Sheet

```bash
# Connection
redis-cli -h host -p port -a password

# Basic operations
SET key value                  GET key
MSET k1 v1 k2 v2              MGET k1 k2
DEL key                       EXISTS key
EXPIRE key 60                 TTL key

# Lists
LPUSH list item               RPUSH list item
LPOP list                     RPOP list
LRANGE list 0 -1              LLEN list

# Sets
SADD set member               SREM set member
SMEMBERS set                  SCARD set

# Hashes
HSET hash field value         HGET hash field
HGETALL hash                  HDEL hash field

# Sorted Sets
ZADD zset score member        ZRANGE zset 0 -1
ZSCORE zset member            ZRANK zset member

# Server
INFO                          CONFIG GET pattern
SAVE                          FLUSHDB
```

---

## Summary

This Redis manual covers:
- ✅ **Setup & Configuration** - Docker-based deployment with authentication
- ✅ **Basic Commands** - Essential Redis operations
- ✅ **Data Types** - Strings, Lists, Sets, Hashes, Sorted Sets
- ✅ **Advanced Features** - Pub/Sub, Transactions, Lua scripting
- ✅ **Go Integration** - Client setup and common patterns
- ✅ **Performance** - Monitoring and optimization
- ✅ **Security** - Best practices and recommendations
- ✅ **Troubleshooting** - Common issues and solutions

Redis is now configured and ready for use in the Lazy Rabbit Reminder project! 🚀
