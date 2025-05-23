# 使用 Flink CDC 监控 MySQL 数据变更并发布到 Kafka

## 一、快速启动

1. 进入本项目文件夹，启动所有服务：

```sh
docker-compose up -d
```

2. 执行 [copy-flink-jars.sh](./copy-flink-jars.sh) 脚本，将 Flink CDC、Kafka 连接器、Kafka client 的 jar 包上传到 flink 和 taskmanager 容器：

```sh
bash copy-flink-jars.sh
```

3. 进入 Flink JobManager 容器，启动 Flink SQL 客户端，创建表：

```sql
CREATE TABLE user_cdc (
  id INT,
  name STRING,
  PRIMARY KEY (id) NOT ENFORCED
) WITH (
  'connector' = 'mysql-cdc',
  'hostname' = 'mysql',
  'port' = '3306',
  'username' = 'root',
  'password' = 'root',
  'database-name' = 'test_db',
  'table-name' = 'user'
);

CREATE TABLE kafka_sink (
  id INT,
  name STRING,
  PRIMARY KEY (id) NOT ENFORCED
) WITH (
  'connector' = 'upsert-kafka',
  'topic' = 'user_changes',
  'properties.bootstrap.servers' = 'kafka:9092',
  'key.format' = 'json',
  'value.format' = 'json'
);

INSERT INTO kafka_sink
SELECT id, name FROM user_cdc;
```

4. 创建 Kafka topic（如未自动创建）：

```sh
docker exec -it kafka kafka-topics.sh --create --topic user_changes --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

5. 启动 [golang](./golang/) 文件夹下的 main.go 文件，即可监听 MySQL 数据变动。

---

## 二、常见问题与解决方法

### 1. Kafka 容器内外访问方式
- **容器内（如 Flink）访问 Kafka**：`kafka:9092`
- **宿主机（如 Go 程序）访问 Kafka**：`localhost:9093`
- `docker-compose.yml` 中 Kafka 推荐配置：

```yaml
  kafka:
    image: wurstmeister/kafka
    container_name: kafka
    ports:
      - "9092:9092"
      - "9093:9093"
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,EXTERNAL://0.0.0.0:9093
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,EXTERNAL://localhost:9093
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_DELETE_TOPIC_ENABLE: 'true'
    depends_on:
      - zookeeper
```

### 2. Flink 任务报 `NoResourceAvailableException`
- 需确保 `taskmanager` 容器已启动并注册到 Flink 集群。
- Flink Web UI 的 Task Managers 页面应有可用 slot。
- 可通过重启 taskmanager 容器解决：
  ```sh
  docker-compose restart taskmanager
  ```

### 3. Flink SQL 报 `StreamCorruptedException: unexpected block data`
- 原因：JobManager 和 TaskManager 的 `/opt/flink/lib/` 目录下 jar 包不一致或有冲突。
- 解决：
  1. 重启 Flink 集群和 SQL Client。

### 4. Kafka topic 没有 leader 或无法消费
- 创建 topic 时 `--replication-factor` 必须为 1（单节点环境）。
- 删除错误 topic 后重建，或重启 Kafka 容器。
- 用如下命令检查 topic 状态：
  ```sh
  docker exec -it kafka kafka-topics.sh --describe --topic user_changes --bootstrap-server localhost:9092
  ```

### 5. 宿主机 Go 程序无法连接 Kafka
- Go 程序应连接 `localhost:9093`，不能用 `kafka:9092`。
- 确保 Kafka 的 `KAFKA_ADVERTISED_LISTENERS` 配置包含 `localhost:9093`。

### 6. 端口冲突或网络不通
- 确认 9092、9093 端口未被占用，防火墙未阻断。
- 容器内可用 `ping kafka` 测试网络。

---

如遇到其他问题，建议优先检查容器日志和 Flink/Kafka 的 Web UI 状态，或在 issue 区留言反馈。

