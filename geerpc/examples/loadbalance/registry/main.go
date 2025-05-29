package main

import (
	"log"
	"net/http"
	"time"

	"geerpc/pkg/registry"
)

func main() {
	// 创建注册中心，设置服务超时时间为5分钟
	r := registry.New(time.Minute * 5)
	// 处理HTTP请求
	r.HandleHTTP("/_geerpc_/registry")

	// 启动注册中心服务
	log.Println("注册中心启动在 http://localhost:9999/_geerpc_/registry")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("注册中心启动错误:", err)
	}
}
