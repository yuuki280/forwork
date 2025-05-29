package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"geerpc/pkg/registry"
	"geerpc/pkg/server"
)

// 用法信息
const usage = `geerpc是一个简单的RPC框架命令行工具

使用方法:
  geerpc [命令] [参数]

可用命令:
  server    启动RPC服务器
  registry  启动服务注册中心

示例:
  geerpc server -addr=:9999       # 在9999端口启动RPC服务器
  geerpc registry -addr=:9999     # 在9999端口启动服务注册中心
`

func main() {
	// 检查命令行参数数量
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	// 解析子命令
	switch os.Args[1] {
	case "server":
		serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
		serverAddr := serverCmd.String("addr", ":9999", "服务器监听地址")
		registryAddr := serverCmd.String("registry", "", "注册中心地址 (可选)")

		serverCmd.Parse(os.Args[2:])
		startServer(*serverAddr, *registryAddr)

	case "registry":
		registryCmd := flag.NewFlagSet("registry", flag.ExitOnError)
		addr := registryCmd.String("addr", ":9999", "注册中心监听地址")
		timeout := registryCmd.Duration("timeout", time.Minute*5, "服务过期时间")

		registryCmd.Parse(os.Args[2:])
		startRegistry(*addr, *timeout)

	default:
		fmt.Println(usage)
	}
}

// 启动RPC服务器
func startServer(addr, registryAddr string) {
	// 创建RPC服务器
	rpcServer := server.NewServer()

	// 监听TCP连接
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("监听错误:", err)
	}
	log.Printf("RPC服务器启动在 %s\n", addr)

	// 如果提供了注册中心地址，则向注册中心注册服务
	if registryAddr != "" {
		registry.Heartbeat(registryAddr, "tcp@"+addr, 0)
		log.Printf("已注册到服务中心 %s\n", registryAddr)
	}

	// 启动服务器
	rpcServer.Accept(l)
}

// 启动服务注册中心
func startRegistry(addr string, timeout time.Duration) {
	// 创建注册中心
	r := registry.New(timeout)
	// 处理HTTP请求
	r.HandleHTTP("/_geerpc_/registry")

	// 启动HTTP服务
	log.Printf("注册中心启动在 http://%s/_geerpc_/registry\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("注册中心启动错误:", err)
	}
}
