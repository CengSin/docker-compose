#!/bin/bash
# 拷贝 Flink 及 Kafka 相关 jar 包到 flink 和 taskmanager 容器

# 定义 jar 包路径
JAR_DIR="/Users/cengsin/docker_workspace/flink_mysql"

# 依次拷贝 jar 包到 flink 和 taskmanager

docker cp "$JAR_DIR/kafka-clients-2.8.0.jar" flink:/opt/flink/lib/
docker cp "$JAR_DIR/kafka-clients-2.8.0.jar" taskmanager:/opt/flink/lib/

docker cp "$JAR_DIR/flink-connector-kafka-1.17.0.jar" flink:/opt/flink/lib/
docker cp "$JAR_DIR/flink-connector-kafka-1.17.0.jar" taskmanager:/opt/flink/lib/

docker cp "$JAR_DIR/flink-sql-connector-mysql-cdc-2.3.0.jar" flink:/opt/flink/lib/
docker cp "$JAR_DIR/flink-sql-connector-mysql-cdc-2.3.0.jar" taskmanager:/opt/flink/lib/

echo "JAR 包已全部拷贝完成！" 