package main

import (
	"log"
	"net"

	"geerpc/pkg/server"
)

type Args struct {
	A, B int
}

type Reply struct {
	Sum int
}

// 定义服务
type Calculator struct{}

// 方法必须满足以下条件：
// - 方法是导出的（首字母大写）
// - 有两个参数，都是导出类型或内建类型
// - 第二个参数是指针
// - 有一个error类型的返回值
func (c *Calculator) Add(args Args, reply *Reply) error {
	reply.Sum = args.A + args.B
	return nil
}

func (c *Calculator) Mul(args Args, reply *Reply) error {
	reply.Sum = args.A * args.B
	return nil
}

func main() {
	// 创建一个新的RPC服务器
	rpcServer := server.NewServer()
	// 注册服务
	if err := rpcServer.Register(&Calculator{}); err != nil {
		log.Fatal("注册服务出错:", err)
	}

	// 监听TCP端口
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("网络错误:", err)
	}
	log.Println("RPC服务器启动在 :9999")

	// 启动服务器
	rpcServer.Accept(l)
}
