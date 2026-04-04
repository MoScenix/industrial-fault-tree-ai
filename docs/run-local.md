# 本地运行说明

## 1. 启动外部依赖

在仓库根目录执行：

```bash
docker compose up -d
```

当前依赖：

- MySQL
- Redis
- Consul
- Etcd
- MinIO
- Milvus
- Jaeger

## 2. 默认凭据

### MySQL

- host: `127.0.0.1`
- port: `3306`
- database: `industrial_fault_tree_ai`
- user: `app`
- password: `app123456`

### Redis

- host: `127.0.0.1`
- port: `6379`
- password: `redis123456`

### MinIO

- endpoint: `127.0.0.1:9000`
- console: `127.0.0.1:9001`
- user: `minioadmin`
- password: `minioadmin123`

### Consul

- ui: `http://127.0.0.1:8500`

### Milvus

- endpoint: `127.0.0.1:19530`
- username: `root`
- password: 空字符串

## 3. 服务配置

当前服务使用：

- `conf/test/conf.yaml`
- 每个服务目录下的 `.env`

已包含本地默认值的服务：

- `app/user/.env`
- `app/graph/.env`
- `app/document/.env`
- `app/ai/.env`
- `app/bff/.env`
- `app/frontend/.env`

## 4. 启动顺序

建议顺序：

1. `user`
2. `graph`
3. `document`
4. `ai`
5. `bff`
6. `frontend`

## 5. 本地主流程验收

建议至少走一遍：

1. 用户注册 / 登录
2. 新建项目
3. 进入工作台
4. 准备暂存
5. 上传项目文档
6. 发起一次 AI 对话
7. 读取当前建议
8. 保存版本
9. 上传个人文档
10. 管理员登录并查看项目总览 / 提示词管理

## 6. 当前重要说明

- 项目目录根路径默认 `/graph`
- 文档全局目录默认 `/document`
- 普通用户只能操作自己的项目
- 管理员可查看和操作全部项目
- 当前 `ChatToModifyGraph` 还是一次请求聚合返回，不是真正的流式 UI
