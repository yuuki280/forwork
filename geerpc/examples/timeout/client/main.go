package main

import (
	"context"
	"log"
	"time"

	"geerpc/pkg/client"
	"geerpc/pkg/protocol"
)

type Args struct {
	A, B int
}

type Reply struct {
	Sum int
}

func main() {
	// 创建客户端连接选项
	option := &protocol.Option{
		ConnectTimeout: time.Second, // 连接超时设置为1秒
	}

	// 连接服务器
	c, err := client.Dial("tcp", "localhost:9998", option)
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	defer c.Close()

	// 测试无延迟方法（应该成功）
	args := &Args{A: 10, B: 20}
	reply := &Reply{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := c.Call(ctx, "TimeoutService.NoDelay", args, reply); err != nil {
		log.Printf("调用NoDelay错误: %v\n", err)
	} else {
		log.Printf("NoDelay成功: %d + %d = %d\n", args.A, args.B, reply.Sum)
	}

	// 测试慢速方法（应该超时）
	args = &Args{A: 10, B: 20}
	reply = &Reply{}
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	if err := c.Call(ctx, "TimeoutService.SlowAdd", args, reply); err != nil {
		log.Printf("调用SlowAdd超时错误（预期行为）: %v\n", err)
	} else {
		log.Printf("SlowAdd成功: %d + %d = %d\n", args.A, args.B, reply.Sum)
	}

	// 测试慢速方法（给足够时间，应该成功）
	args = &Args{A: 30, B: 40}
	reply = &Reply{}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := c.Call(ctx, "TimeoutService.SlowAdd", args, reply); err != nil {
		log.Printf("调用SlowAdd错误: %v\n", err)
	} else {
		log.Printf("SlowAdd成功（使用更长超时）: %d + %d = %d\n", args.A, args.B, reply.Sum)
	}
}
