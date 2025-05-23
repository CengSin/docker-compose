# etcd 微服务注册中心示例

本项目演示了如何基于 etcd 搭建微服务注册中心，并通过 go-micro 框架实现服务注册与发现。

## 目录结构

- go-service/  —— 服务端，注册服务到 etcd
- go-client/   —— 客户端，调用注册在 etcd 中的服务
- docker-compose.yml —— etcd 集群部署配置

## 组件说明

### 1. etcd 集群

使用 docker-compose 启动 3 节点 etcd 集群，作为微服务的注册与配置中心。

- etcd1: 2379/2380 端口
- etcd2: 2381/2382 端口
- etcd3: 2383/2384 端口

集群通过 etcd_net 网络互联，数据持久化到本地卷。

### 2. go-service（服务端）

- 使用 go-micro 框架，注册名为 `my.custom.service` 的服务到 etcd。
- 服务实现了一个简单的 `Start` 方法，接收请求并返回问候语。
- 主要文件：
  - main.go：服务注册、启动逻辑
  - server.go：服务实现
  - plugin.go：etcd 插件引入

### 3. go-client（客户端）

- 作为服务消费者，通过 etcd 发现并调用 `my.custom.service` 服务。
- 主要文件：
  - main.go：服务发现与远程调用逻辑
  - plugin.go：etcd 插件引入

## 快速开始

1. 启动 etcd 集群

```bash
docker-compose up -d
```

2. 启动服务端

```bash
cd go-service
# 构建并运行服务
# go build -o service .
# ./service
```

3. 启动客户端

```bash
cd go-client
# 构建并运行客户端
# go build -o client .
# ./client
```

## 依赖

- Go 1.24+
- go-micro v5
- etcd v3.5+

## 参考
- [go-micro 官方文档](https://go-micro.dev/)
- [etcd 官方文档](https://etcd.io/)