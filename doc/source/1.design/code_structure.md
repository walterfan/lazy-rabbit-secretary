# Code structure

## Overview

## Main Application Flow

### 1. **Entry Point (`main.go`)**
```go
func main() {
    cmd.Execute()
}
```
The application starts by calling `cmd.Execute()`, which delegates to the Cobra command-line framework.

### 2. **Root Command Initialization (`cmd/root.go`)**

The root command sets up the application infrastructure through a **PersistentPreRun** hook that executes before any subcommand:

#### **Initialization Sequence:**
1. **`initLogger()`** - Sets up structured logging with Zap
2. **`initConfig()`** - Loads configuration from YAML files using Viper
3. **`initRedis()`** - Connects to Redis database and initializes client
4. **`initAuth()`** - Sets up authentication service with JWT manager

#### **Key Global Variables:**
- `logger *zap.Logger` - Application logger
- `rdb *redis.Client` - Redis client for data storage
- `authService *auth.AuthService` - Authentication service

### 3. **Command Structure**

The application uses Cobra to provide multiple subcommands:

#### **Available Commands:**

1. **`server`** (`cmd/server.go`)
   - Starts HTTP/HTTPS web server
   - Runs two concurrent services:
     - **Web API Service**: Handles HTTP requests with Gin framework
     - **Task Manager**: Monitors and manages scheduled tasks
   - Graceful shutdown on SIGINT/SIGTERM

2. **`blog`** (`cmd/blog.go`)
   - Generates daily technical blog posts
   - Supports multiple languages (English/Chinese)
   - Uses LLM (OpenAI) for content generation
   - Configurable via command-line flags

3. **`check-tasks`** (`cmd/cronjob.go`)
   - Runs scheduled task monitoring
   - Manages task expiry using Redis
   - Standalone cron job functionality

4. **`version`** (`cmd/version.go`)
   - Displays application version information

### 4. **Web Server Flow (`cmd/server.go`)**

When running the `server` command:

```go
// Start Web API Service
webService := api.NewWebApiService(logger, rdb, GetAuthService())
go webService.Run()

// Start Task Manager
tm := task.NewTaskManager(logger, rdb)
go tm.CheckTasks()

// Wait for shutdown signal
signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
<-signalChan
```

### 5. **Web API Service (`internal/api/web_service.go`)**

The web service provides:
- **Health check**: `/ping` endpoint
- **News API**: `/api/v1/news` - Fetches news from Redis
- **Commands API**: `/api/v1/commands` - Lists available commands
- **Static file serving**: Serves frontend files
- **Authentication middleware**: JWT-based auth for protected routes

### 6. **Configuration Management**

- **Config files**: `config/config.yaml` (default) or custom via `--config` flag
- **Environment variables**: Loaded from `.env` file
- **Redis configuration**: Host, port, password from environment
- **JWT certificates**: Private/public keys for token signing

### 7. **Data Flow**

```
main.go → cmd.Execute() → rootCmd.PersistentPreRun → [initLogger, initConfig, initRedis, initAuth] → subcommand execution
```

### 8. **Key Dependencies**

- **Cobra**: Command-line interface framework
- **Viper**: Configuration management
- **Zap**: Structured logging
- **Redis**: Data storage and caching
- **Gin**: HTTP web framework
- **JWT**: Authentication tokens
- **OpenAI**: LLM for blog generation

### 9. **Service Architecture**

The application follows a **microservice-like** pattern within a single binary:
- **Web Service**: HTTP API endpoints
- **Task Service**: Background task processing
- **Auth Service**: Authentication and authorization
- **LLM Service**: AI-powered content generation

This architecture allows the application to serve as both a web server and a task processing system, with shared infrastructure (Redis, logging, config) across all components.


## code structure

```
lazy-rabbit-reminder/backend/
```


```

go-project/            # 项目根目录
├── cmd/               # 程序入口（命令行/服务启动）
│   └── main.go        # 初始化依赖、启动服务（仅负责“组装”，不写业务逻辑）
├── internal/          # 私有代码（不对外暴露的业务核心，Go 编译时禁止外部导入）
│   ├── domain/        # 领域层（业务实体/核心规则，与框架无关）
│   │   └── user/      # 按业务模块拆分（如用户模块、订单模块）
│   │       ├── model.go    # 业务实体（User 结构体，定义核心属性和规则）
│   │       └── service.go  # 领域服务（纯业务逻辑，不依赖数据访问细节）
│   ├── repository/    # 数据访问层（Repository 模式，隔离数据来源）
│   │   └── user/
│   │       ├── repo.go     # 数据访问接口（定义“做什么”，如 UserRepo 接口）
│   │       └── mysql.go    # 接口实现（MySQL 实现，“怎么做”，依赖具体数据库）
│   ├── service/       # 应用服务层（协调领域层与数据层，处理跨模块逻辑）
│   │   └── user/
│   │       └── service.go  # 应用服务（调用 domain 业务逻辑 + repository 数据操作）
│   └── handler/       # 接口适配层（处理 HTTP/gRPC 等外部请求，MVC 中的 Controller）
│       └── user/
│           └── handler.go  # 处理 HTTP 请求（参数校验、调用 service、返回响应）
├── pkg/               # 公共代码（可对外复用的工具/组件，非业务相关）
│   ├── db/            # 数据库工具（MySQL 连接池、初始化）
│   ├── http/          # HTTP 工具（路由封装、中间件）
│   ├── logger/        # 日志工具（统一日志格式、输出）
│   └── config/        # 配置工具（读取环境变量、配置文件）
├── api/               # 接口定义（对外暴露的 API 契约，如 OpenAPI/Swagger、gRPC proto）
│   └── user.proto     # gRPC 接口定义（若用 gRPC）
├── configs/           # 配置文件（yaml/json，不包含敏感信息）
├── scripts/           # 脚本（部署、数据库迁移、构建脚本）
└── go.mod/go.sum      # Go 模块依赖

```