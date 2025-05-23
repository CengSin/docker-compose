package main

import (
	"log"
	"time"

	"github.com/micro/plugins/v5/registry/etcd"
	"go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

func main() {
	r := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"http://etcd1.etcd-work.orb.local:2379"}
		options.Timeout = 5 * time.Second
		options.Secure = false
	})

	service := micro.NewService(
		micro.Name("my.custom.service"),
		micro.Handle(&MyCustomServer{}),
		micro.Registry(r),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal("Service run error: ", err)
	}
}
