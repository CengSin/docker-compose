# 基于 ParadeDB 镜像（包含 pg_bm25 和 pgvector），添加中文分词支持
FROM paradedb/paradedb:latest

# 设置环境变量
ENV DEBIAN_FRONTEND=noninteractive

# 安装编译工具和依赖（用于编译 zhparser）
# 注意：ParadeDB 使用 PostgreSQL 18，但为了兼容性，我们尝试安装所有版本
RUN apt-get update && \
    (apt-get install -y postgresql-server-dev-18 || \
     apt-get install -y postgresql-server-dev-17 || \
     apt-get install -y postgresql-server-dev-16 || \
     apt-get install -y postgresql-server-dev-all) && \
    apt-get install -y \
    build-essential \
    git \
    libscws-dev \
    make \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# 安装 zhparser（中文分词扩展）
RUN cd /tmp && \
    git clone --depth 1 https://github.com/amutu/zhparser.git && \
    cd zhparser && \
    make && \
    make install && \
    cd / && \
    rm -rf /tmp/zhparser

# 清理编译工具（保留运行时依赖）
RUN apt-get purge -y build-essential git make gcc && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
