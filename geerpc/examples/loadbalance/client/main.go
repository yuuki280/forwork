package main

import (
	"context"
	"log"
	"sync"
	"time"

	"geerpc/pkg/discovery"
	"geerpc/pkg/xclient"
)

type Args struct {
	A, B int
}

type Reply struct {
	Sum   int
	Port  int
	Delay time.Duration
}

func main() {
	// 创建注册中心客户端
	d := discovery.NewRegistryDiscovery("http://localhost:9999/_geerpc_/registry", 0)

	// 使用随机策略的XClient
	randomXClient := xclient.NewXClient(d, discovery.RandomSelect, nil)
	defer randomXClient.Close()

	// 使用轮询策略的XClient
	roundRobinXClient := xclient.NewXClient(d, discovery.RoundRobinSelect, nil)
	defer roundRobinXClient.Close()

	// 调用次数
	n := 5
	var wg sync.WaitGroup
	wg.Add(2 * n)

	// 使用随机策略调用n次
	go func() {
		for i := 0; i < n; i++ {
			testCall(randomXClient, context.Background(), "BalanceService.Add", &Args{A: 10, B: i}, "Random")
			wg.Done()
		}
	}()

	// 使用轮询策略调用n次
	go func() {
		for i := 0; i < n; i++ {
			testCall(roundRobinXClient, context.Background(), "BalanceService.Add", &Args{A: 10, B: i}, "RoundRobin")
			wg.Done()
		}
	}()

	wg.Wait()

	// 测试广播调用
	log.Println("测试广播调用...")
	testBroadcast(randomXClient)
}

func testCall(xc *xclient.XClient, ctx context.Context, serviceMethod string, args *Args, strategy string) {
	reply := &Reply{}
	start := time.Now()
	if err := xc.Call(ctx, serviceMethod, args, reply); err != nil {
		log.Printf("%s 调用失败: %v\n", strategy, err)
	} else {
		log.Printf("%s 调用成功 - [port:%d] %d + %d = %d, 服务端延迟:%v, 总耗时:%v\n",
			strategy, reply.Port, args.A, args.B, reply.Sum, reply.Delay, time.Since(start))
	}
}

func testBroadcast(xc *xclient.XClient) {
	args := &Args{A: 10, B: 20}
	reply := &Reply{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := xc.Broadcast(ctx, "BalanceService.Add", args, reply); err != nil {
		log.Printf("广播调用失败: %v\n", err)
	} else {
		log.Printf("广播调用成功 - 最后一个返回: [port:%d] %d + %d = %d\n",
			reply.Port, args.A, args.B, reply.Sum)
	}
}
