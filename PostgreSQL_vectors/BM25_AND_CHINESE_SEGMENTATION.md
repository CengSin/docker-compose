# BM25 和中文分词使用指南

## 功能说明

本 PostgreSQL 配置已包含以下扩展：

1. **pgvector** - 向量相似度搜索
2. **pg_bm25** - BM25 全文搜索算法（由 ParadeDB 提供）
3. **zhparser** - 中文分词支持

## 使用方法

### 1. 启动服务

```bash
# 构建并启动容器
docker-compose up -d --build

# 查看日志确认扩展已加载
docker-compose logs postgres | grep -i extension
```

### 2. 验证扩展安装

```bash
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "\dx"
```

应该看到以下扩展：
- vector
- pg_bm25
- zhparser

### 3. 使用 BM25 全文搜索

#### 创建支持 BM25 的表

```sql
-- 创建示例表
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    title TEXT,
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 使用 pg_bm25 创建搜索索引
CREATE INDEX idx_documents_bm25 ON documents 
USING bm25 (title, content);

-- 插入测试数据
INSERT INTO documents (title, content) VALUES
('PostgreSQL 教程', 'PostgreSQL 是一个强大的开源关系型数据库'),
('向量搜索', 'pgvector 提供了高效的向量相似度搜索功能'),
('全文搜索', 'BM25 算法提供了优秀的全文搜索性能');
```

#### 执行 BM25 搜索

```sql
-- 使用 BM25 搜索
SELECT * FROM documents 
WHERE bm25_search(title, content, 'PostgreSQL') 
ORDER BY bm25_rank(title, content, 'PostgreSQL') DESC;
```

### 4. 使用中文分词全文搜索

#### 创建中文全文搜索索引

```sql
-- 创建中文内容表
CREATE TABLE chinese_docs (
    id SERIAL PRIMARY KEY,
    title TEXT,
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建中文全文搜索索引（使用 chinese_zh 配置）
CREATE INDEX idx_chinese_docs_fts ON chinese_docs 
USING GIN (to_tsvector('chinese_zh', title || ' ' || content));

-- 插入中文测试数据
INSERT INTO chinese_docs (title, content) VALUES
('数据库教程', 'PostgreSQL 是一个功能强大的开源数据库系统'),
('向量搜索技术', 'pgvector 扩展支持高效的向量相似度搜索'),
('全文检索', '中文分词使得全文搜索更加准确和高效');
```

#### 执行中文全文搜索

```sql
-- 使用中文分词搜索
SELECT 
    id,
    title,
    content,
    ts_rank(to_tsvector('chinese_zh', title || ' ' || content), 
            to_tsquery('chinese_zh', '数据库')) AS rank
FROM chinese_docs
WHERE to_tsvector('chinese_zh', title || ' ' || content) 
      @@ to_tsquery('chinese_zh', '数据库')
ORDER BY rank DESC;
```

#### 高级中文搜索示例

```sql
-- 搜索多个关键词（AND）
SELECT * FROM chinese_docs
WHERE to_tsvector('chinese_zh', title || ' ' || content) 
      @@ to_tsquery('chinese_zh', '数据库 & 搜索');

-- 搜索多个关键词（OR）
SELECT * FROM chinese_docs
WHERE to_tsvector('chinese_zh', title || ' ' || content) 
      @@ to_tsquery('chinese_zh', '数据库 | 向量');

-- 使用短语搜索
SELECT * FROM chinese_docs
WHERE to_tsvector('chinese_zh', title || ' ' || content) 
      @@ phraseto_tsquery('chinese_zh', '全文检索');
```

### 5. 组合使用 BM25 和中文分词

```sql
-- 创建同时支持 BM25 和中文分词的混合搜索表
CREATE TABLE hybrid_docs (
    id SERIAL PRIMARY KEY,
    title TEXT,
    content TEXT,
    embedding vector(1536),  -- 向量字段（用于向量搜索）
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建 BM25 索引
CREATE INDEX idx_hybrid_bm25 ON hybrid_docs 
USING bm25 (title, content);

-- 创建中文全文搜索索引
CREATE INDEX idx_hybrid_fts ON hybrid_docs 
USING GIN (to_tsvector('chinese_zh', title || ' ' || content));

-- 创建向量索引
CREATE INDEX idx_hybrid_vector ON hybrid_docs 
USING ivfflat (embedding vector_cosine_ops);
```

## 配置说明

### 中文分词配置

默认的中文全文搜索配置 `chinese_zh` 已自动创建，包含以下词性映射：
- `n` - 名词
- `v` - 动词
- `a` - 形容词
- `i` - 成语
- `e` - 习用语
- `l` - 习用语

如需自定义配置，可以修改 `init-scripts/01-enable-pgvector.sql` 中的配置。

### 性能优化建议

1. **BM25 搜索**：适合大规模文档的全文搜索，性能优于传统全文搜索
2. **中文分词**：对于中文内容，建议使用 `chinese_zh` 配置而不是默认的 `english`
3. **索引优化**：根据查询模式选择合适的索引类型

## 故障排查

### 扩展未加载

```bash
# 检查扩展是否已安装
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "\dx"

# 手动创建扩展
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "CREATE EXTENSION IF NOT EXISTS pg_bm25;"
docker exec -it postgres-vectors psql -U postgres -d vectors_db -c "CREATE EXTENSION IF NOT EXISTS zhparser;"
```

### zhparser 编译失败

如果 zhparser 编译失败，可以尝试：
1. 检查日志：`docker-compose logs postgres`
2. 确保 libscws-dev 已安装
3. 重新构建镜像：`docker-compose build --no-cache`

## 参考资源

- [ParadeDB 文档](https://docs.paradedb.com/)
- [pg_bm25 文档](https://docs.paradedb.com/documentation/search/bm25)
- [zhparser GitHub](https://github.com/amutu/zhparser)
- [PostgreSQL 全文搜索文档](https://www.postgresql.org/docs/current/textsearch.html)
