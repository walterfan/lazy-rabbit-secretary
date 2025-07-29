```{contents} Table of Contents
:depth: 3
```

# Overview

## 📚 基本功能

* 

## 📁  项目结构

```
your-project/
|-- main.go               # 主应用入口点
├── cmd/
│   │── version.go         # version command
│   └── server.go          # server command
├── internal/             # 内部私有代码，只能在本模块内使用（Go 1.4+特性）
│   └── service/
├── pkg/                  # 可被外部项目导入使用的公共代码
│   └── utils/
├── api/                  # API 相关定义（如 OpenAPI/Protobuf 等）
├── configs/              # 配置文件
├── scripts/              # 启动、构建、部署脚本
├── build/                # Dockerfile、CI/CD 配置等构建相关
├── test/                 # 额外测试数据或集成测试
├── go.mod
├── go.sum
└── README.md
```

---

## 📦 各目录说明及好处

### `cmd/`

* **说明**：每个文件对应一个可构建的主程序（如 `version.go`）。
* **好处**：

  * 支持多个命令行工具或微服务。
  * 避免一个庞大的 `main.go` 雪崩。

### `internal/`

* **说明**：放置模块私有代码，其他模块无法导入。
* **好处**：

  * 强制封装。Go 会阻止 `internal/` 之外的模块访问它。
  * 有助于模块解耦。

### `pkg/`

* **说明**：可被外部项目导入的公共库代码。
* **好处**：

  * 明确哪些代码是复用的 API。
  * 与 `internal/` 分开，降低耦合。

### `api/`

* **说明**：协议定义，如 REST 的 OpenAPI、gRPC 的 proto 文件。
* **好处**：

  * 清晰的接口契约。
  * 支持自动代码生成。

### `configs/`

* **说明**：配置模板或默认配置。
* **好处**：

  * 配置管理更清晰。
  * 容易与 k8s/Consul 等配置中心对接。

### `scripts/`

* **说明**：构建、测试、部署脚本。
* **好处**：

  * 自动化部署，CI/CD 更容易维护。

### `build/`

* **说明**：放置构建工具文件，例如 Dockerfile、Makefile、Helm Chart。
* **好处**：

  * CI/CD 更加清晰，构建流程标准化。

### `test/`

* **说明**：非单元测试的其他测试资源，如集成测试或测试数据。
* **好处**：

  * 分离测试数据和业务逻辑代码。

---

## 🔧 示例应用场景

比如一个微服务项目：

```
my-service/
├── cmd/my-service/main.go        // 程序入口
├── internal/order/handler.go     // 业务逻辑
├── pkg/logger/logger.go          // 通用日志库
├── api/proto/order.proto         // gRPC 接口定义
├── configs/config.yaml           // 配置文件模板
├── build/Dockerfile              // Docker 构建文件
├── scripts/run-dev.sh            // 启动脚本
└── go.mod
```




