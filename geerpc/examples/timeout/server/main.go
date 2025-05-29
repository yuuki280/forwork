package main

import (
	"log"
	"net"
	"time"

	"geerpc/pkg/server"
)

type Args struct {
	A, B int
}

type Reply struct {
	Sum int
}

// 定义服务
type TimeoutService struct{}

// 普通方法，立即返回
func (s *TimeoutService) NoDelay(args Args, reply *Reply) error {
	reply.Sum = args.A + args.B
	return nil
}

// 延迟方法，模拟耗时操作
func (s *TimeoutService) SlowAdd(args Args, reply *Reply) error {
	// 模拟耗时操作
	time.Sleep(time.Second * 2)
	reply.Sum = args.A + args.B
	return nil
}

func main() {
	// 创建一个新的RPC服务器
	rpcServer := server.NewServer()
	// 注册服务
	if err := rpcServer.Register(&TimeoutService{}); err != nil {
		log.Fatal("注册服务出错:", err)
	}

	// 监听TCP端口
	l, err := net.Listen("tcp", ":9998")
	if err != nil {
		log.Fatal("网络错误:", err)
	}
	log.Println("超时示例RPC服务器启动在 :9998")

	// 启动服务器
	rpcServer.Accept(l)
}
