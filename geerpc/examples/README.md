# GeeRPC 框架示例

本目录包含GeeRPC框架的几个使用示例，展示了框架的不同功能。

## 基础示例 (basic)

展示GeeRPC最基本的服务注册和调用方式。

运行方法：
```bash
# 终端1: 启动服务器
cd examples/basic/server
go run main.go

# 终端2: 运行客户端
cd examples/basic/client
go run main.go
```

## 超时处理示例 (timeout)

展示GeeRPC的超时处理机制，包括连接超时和处理超时。

运行方法：
```bash
# 终端1: 启动服务器
cd examples/timeout/server
go run main.go

# 终端2: 运行客户端
cd examples/timeout/client
go run main.go
```

## 负载均衡示例 (loadbalance)

展示GeeRPC的服务发现和负载均衡功能，包括随机选择和轮询策略，以及广播调用。

运行方法：
```bash
# 终端1: 启动注册中心
cd examples/loadbalance/registry
go run main.go

# 终端2: 启动服务器1
cd examples/loadbalance/server
go run main.go -port=9001

# 终端3: 启动服务器2
cd examples/loadbalance/server
go run main.go -port=9002

# 终端4: 启动服务器3
cd examples/loadbalance/server
go run main.go -port=9003

# 终端5: 运行客户端
cd examples/loadbalance/client
go run main.go
```

## 示例说明

1. **基础示例**：展示了如何定义服务、注册服务以及调用远程方法。
2. **超时处理示例**：展示了如何设置和处理调用超时，增强系统的健壮性。
3. **负载均衡示例**：展示了如何实现服务注册、发现和负载均衡，适用于分布式环境。

每个示例都包含服务端和客户端代码，可以独立运行。请确保先启动服务器再运行客户端。 