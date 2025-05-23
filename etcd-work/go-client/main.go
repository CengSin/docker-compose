package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/micro/plugins/v5/registry/etcd"
	"go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

type StartRequest struct {
	Name string
}

type StartResponse struct {
	Ans string
}

func main() {
	// 创建与服务端相同的注册中心配置
	r := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"http://etcd1.etcd-work.orb.local:2379"}
		options.Timeout = 5 * time.Second
		options.Secure = false
	})

	// 创建服务客户端
	service := micro.NewService(
		micro.Name("my.custom.service.client"),
		micro.Registry(r),
	)

	service.Init()

	// 创建服务的客户端
	cli := service.Client()

	// 创建请求对象
	request := &StartRequest{
		Name: "测试客户端",
	}

	// 创建响应对象
	response := &StartResponse{}

	// 调用远程服务
	req := cli.NewRequest("my.custom.service", "MyCustomServer.Start", request)
	err := cli.Call(context.Background(), req, response)
	if err != nil {
		log.Fatal("调用服务失败：", err)
	}

	fmt.Printf("服务响应：%v\n", response.Ans)
}
