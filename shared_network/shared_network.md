# shared_network

shared_network是一个自建的docker内部共享网络，用来让不同的docker-compose内的容器和其他自建的容器进行通信。

## 创建

```shell
docker network create shared-network
```

## 其他容器启动时指定网络

```shell

docker run -d --name container1 --network shared-network my-image
docker run -d --name container2 --network shared-network my-image
```

## docker-compose指定网络

修改docker-compose.yml文件，添加以下内容

```shell
networks:
  shared-net:
    external: true
    name: shared-network
```

这样就可以在容器内部直接使用container_name:port来访问服务了。