package main

import (
	"context"
	"log"
	"time"

	"geerpc/pkg/client"
)

type Args struct {
	A, B int
}

type Reply struct {
	Sum int
}

func main() {
	// 连接服务器
	c, err := client.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	defer c.Close()

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用Add方法
	args := &Args{A: 10, B: 20}
	reply := &Reply{}
	if err := c.Call(ctx, "Calculator.Add", args, reply); err != nil {
		log.Fatal("调用Add错误:", err)
	}
	log.Printf("Add: %d + %d = %d\n", args.A, args.B, reply.Sum)

	// 调用Mul方法
	args = &Args{A: 10, B: 20}
	reply = &Reply{}
	if err := c.Call(ctx, "Calculator.Mul", args, reply); err != nil {
		log.Fatal("调用Mul错误:", err)
	}
	log.Printf("Mul: %d * %d = %d\n", args.A, args.B, reply.Sum)
}
