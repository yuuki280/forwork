package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"geerpc/pkg/registry"
	"geerpc/pkg/server"
)

type Args struct {
	A, B int
}

type Reply struct {
	Sum   int
	Port  int
	Delay time.Duration
}

// 服务实现
type BalanceService struct {
	addr string
	port int
}

// 加法方法
func (s *BalanceService) Add(args Args, reply *Reply) error {
	// 随机延迟，模拟不同服务器处理时间不同
	delay := time.Duration(200+s.port%3*100) * time.Millisecond
	time.Sleep(delay)

	reply.Sum = args.A + args.B
	reply.Port = s.port
	reply.Delay = delay
	return nil
}

func main() {
	// 通过命令行参数指定端口，便于启动多个服务器
	var port int
	flag.IntVar(&port, "port", 9001, "服务端口")
	flag.Parse()

	addr := fmt.Sprintf(":%d", port)
	// 创建一个新的RPC服务器
	rpcServer := server.NewServer()
	// 注册服务
	service := &BalanceService{
		addr: addr,
		port: port,
	}
	if err := rpcServer.Register(service); err != nil {
		log.Fatal("注册服务出错:", err)
	}

	// 启动注册中心的HTTP服务
	registry.HandleHTTP()

	// 监听TCP端口
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("网络错误:", err)
	}
	log.Printf("RPC服务器启动在 %s\n", addr)

	// 向注册中心注册服务
	registry.Heartbeat("http://localhost:9999/_geerpc_/registry", "tcp@"+addr, 0)

	// 启动服务器
	rpcServer.Accept(l)
}
