# PostgreSQL with pgvector Docker 配置

## 功能特性

- ✅ PostgreSQL 18 数据库（基于 ParadeDB）
- ✅ pgvector 扩展（用于向量相似度搜索）
- ✅ pg_bm25 扩展（BM25 全文搜索算法）
- ✅ zhparser 扩展（中文分词支持）
- ✅ 数据持久化存储
- ✅ 配置文件外部挂载
- ✅ 支持从 MySQL 迁移数据

## 快速开始

### 1. 启动服务

```bash
# 首次启动需要构建自定义镜像（包含中文分词支持）
docker-compose up -d --build

# 后续启动
docker-compose up -d
```

### 2. 验证扩展安装

```bash
# 查看所有已安装的扩展
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "\dx"

# 应该看到以下扩展：
# - vector (pgvector)
# - pg_bm25 (BM25 全文搜索)
# - zhparser (中文分词)
```

### 3. 连接数据库

```bash
# 使用 psql 连接
docker exec -it postgres-vectors psql -U postgres -d vectors_db

# 或从外部连接（端口 5432）
psql -h localhost -p 5432 -U postgres -d vectors_db
```

默认密码：`postgres`

## 从 MySQL 迁移数据

### 方法 1: 使用 pgloader（推荐）

1. **安装 pgloader**（如果未安装）：
   ```bash
   # macOS
   brew install pgloader
   
   # Ubuntu/Debian
   sudo apt-get install pgloader
   ```

2. **创建迁移脚本** `migrate.load`：
   ```lisp
LOAD DATABASE
     FROM mysql://root:rootpassword@localhost:3306/ivanka_content
     INTO postgresql://postgres:postgres@localhost:5432/vectors_db

WITH include drop, create tables, create indexes, reset sequences

SET maintenance_work_mem to '256MB',
    work_mem to '128MB',
    search_path to 'public'

CAST 
     -- 核心修复：显式匹配带长度的 tinyint 和 boolean 别名
     type tinyint when (= 1 precision) to boolean drop typemod using tinyint-to-boolean,
     type tinyint to boolean using tinyint-to-boolean,
     
     -- 处理时间
     type datetime to timestamptz
          drop default drop not null using zero-dates-to-null,
     type date drop not null drop default using zero-dates-to-null,
     
     -- 其他转换
     type year to integer
;
   ```

3. **执行迁移**：
   ```bash
   pgloader migrate.load
   ```

### 方法 2: 使用 mysqldump + 手动转换

1. **导出 MySQL 数据**：
   ```bash
   mysqldump -u user -p source_db > mysql_dump.sql
   ```

2. **转换并导入**（需要手动调整 SQL 语法差异）

3. **或使用在线转换工具**：
   - https://www.sqlines.com/online
   - https://github.com/dalibo/sql-migrate

### 方法 3: 使用 Python 脚本迁移

将迁移脚本放在 `init-scripts/` 目录下，容器启动时会自动执行。

## 目录结构

```
PostgreSQL_vectors/
├── docker-compose.yml          # Docker Compose 配置
├── Dockerfile                  # 自定义镜像（添加中文分词支持）
├── data/                       # 数据持久化目录
│   └── postgres/              # PostgreSQL 数据文件
├── conf/                       # 配置文件目录
│   ├── postgresql.conf        # PostgreSQL 主配置
│   └── pg_hba.conf            # 客户端认证配置
├── init-scripts/              # 初始化脚本目录
│   └── 01-enable-pgvector.sql # 自动启用所有扩展
├── BM25_AND_CHINESE_SEGMENTATION.md  # BM25 和中文分词使用指南
└── README.md                   # 本文件
```

## 扩展功能

### BM25 全文搜索

本配置包含 `pg_bm25` 扩展，提供高性能的 BM25 全文搜索算法。详细使用方法请参考 [BM25_AND_CHINESE_SEGMENTATION.md](./BM25_AND_CHINESE_SEGMENTATION.md)。

### 中文分词

本配置包含 `zhparser` 扩展，支持中文全文搜索。默认创建了 `chinese_zh` 全文搜索配置。详细使用方法请参考 [BM25_AND_CHINESE_SEGMENTATION.md](./BM25_AND_CHINESE_SEGMENTATION.md)。

## 配置说明

### 修改数据库密码

编辑 `docker-compose.yml` 中的 `POSTGRES_PASSWORD` 环境变量。

### 自定义 PostgreSQL 配置

编辑 `conf/postgresql.conf`，修改后重启容器：
```bash
docker-compose restart postgres
```

### 添加初始化脚本

将 SQL 脚本放入 `init-scripts/` 目录，文件名按字母顺序执行（如 `01-xxx.sql`, `02-xxx.sql`）。

## 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f postgres

# 进入容器
docker exec -it postgres-vectors bash

# 备份数据库
docker exec postgres-vectors pg_dump -U postgres vectors_db > backup.sql

# 恢复数据库
docker exec -i postgres-vectors psql -U postgres vectors_db < backup.sql

# 查看数据库大小
docker exec postgres-vectors psql -U postgres -c "SELECT pg_size_pretty(pg_database_size('vectors_db'));"
```

## 健康检查

容器包含健康检查，可通过以下命令查看状态：
```bash
docker ps
```

## 注意事项

1. **首次启动**：数据目录会自动创建，所有扩展（pgvector、pg_bm25、zhparser）会自动启用
2. **镜像构建**：首次启动需要构建自定义镜像（包含中文分词支持），可能需要几分钟时间
3. **配置文件**：修改 `conf/` 目录下的配置文件后需要重启容器
4. **数据备份**：定期备份 `data/postgres/` 目录或使用 `pg_dump`
5. **生产环境**：建议修改默认密码和 `pg_hba.conf` 中的访问权限
6. **性能调优**：根据实际负载调整 `postgresql.conf` 中的内存和连接数设置
7. **中文搜索**：默认全文搜索配置已设置为 `chinese_zh`，适合中文内容搜索

## 故障排查

### 容器无法启动

检查端口是否被占用：
```bash
lsof -i :5432
```

### 扩展未加载

检查日志：
```bash
docker-compose logs postgres | grep -i extension
```

手动启用扩展：
```bash
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "CREATE EXTENSION IF NOT EXISTS vector;"
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "CREATE EXTENSION IF NOT EXISTS pg_bm25;"
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "CREATE EXTENSION IF NOT EXISTS zhparser;"
```

### zhparser 编译失败

如果构建镜像时 zhparser 编译失败：
1. 检查日志：`docker-compose logs postgres`
2. 重新构建：`docker-compose build --no-cache`
3. 确保有足够的磁盘空间和内存

### 配置文件未生效

确保配置文件路径正确，重启容器：
```bash
docker-compose restart postgres
```