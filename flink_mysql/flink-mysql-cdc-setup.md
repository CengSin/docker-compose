# Flink MySQL CDC 集群部署与常见问题总结

## 1. 为什么需要单独创建 TaskManager 服务？

Flink 集群由 JobManager 和 TaskManager 组成：
- **JobManager** 负责作业调度和资源管理。
- **TaskManager** 负责实际的数据处理和任务执行。

在 Docker Compose 部署时，**只启动 JobManager 是无法运行实际作业的**，因为没有可用的 slot。必须单独创建 TaskManager 服务，JobManager 才能将作业分发给 TaskManager 执行。

## 2. 为什么要在 JobManager 和 TaskManager 都上传 CDC jar 包？如何上传？

Flink SQL CDC 连接器（如 `flink-sql-connector-mysql-cdc-2.3.0.jar`）用于捕获 MySQL 的变更数据。

- **原因**：Flink 作业在运行时，JobManager 会将作业的 class 分发到 TaskManager。如果 TaskManager 没有 CDC 相关 jar 包，会导致序列化/反序列化失败，出现 `StreamCorruptedException` 等错误。
- **结论**：**CDC jar 包必须同时放在 JobManager 和所有 TaskManager 的 `/opt/flink/lib` 目录下。**

### 上传方法
假设你已经下载好 CDC jar 包：

```sh
docker cp flink-sql-connector-mysql-cdc-2.3.0.jar flink:/opt/flink/lib/
docker cp flink-sql-connector-mysql-cdc-2.3.0.jar taskmanager:/opt/flink/lib/
```

上传后，建议重启 Flink 集群：

```sh
docker-compose down
docker-compose up -d
```

## 3. 本次问题总结与解决过程

### 问题现象
- Flink SQL 能查到 MySQL 初始快照，但 MySQL 表变更无法实时同步到 Flink。
- Flink SQL 提交作业时报 `NoResourceAvailableException`。
- Flink Web UI 的 Task Managers 页面为空。
- 后续出现 `StreamCorruptedException: unexpected block data`。

### 排查与解决过程
1. **未启动 TaskManager**：
   - 只启动了 JobManager，导致没有可用 slot，作业无法运行。
   - 解决：在 `docker-compose.yml` 中添加 `taskmanager` 服务。
2. **TaskManager 未注册**：
   - TaskManager 容器未正常启动或未能注册到 JobManager。
   - 解决：检查容器日志，确保网络和配置无误。
3. **CDC jar 包未同步到所有节点**：
   - 只在 JobManager 上传了 CDC jar，TaskManager 缺失导致作业运行时报错。
   - 解决：将 CDC jar 包同时上传到 JobManager 和 TaskManager 的 `/opt/flink/lib` 目录。
4. **最终验证**：
   - 重新提交 Flink SQL 任务，MySQL 表变更可以被实时捕获，问题解决。

---

如需更多帮助，欢迎随时提问！ 