# document

文档管理微服务骨架由 `cwgo` 生成，当前只保留接口和查询层占位，不填具体业务实现。

## 当前职责

- 解析指定 PDF 为个人文档
- 解析指定 PDF 为项目文档
- 查询文档详情与列表
- 根据 `user_id + query` 或 `user_id + project_id + query` 返回搜索结果

## 目录说明

| catalog | introduce |
| ---- | ---- |
| `conf` | 配置文件 |
| `main.go` | 启动入口 |
| `handler.go` | RPC handler |
| `biz/service` | service 层占位 |
| `biz/dal` | `cwgo` 默认生成的底层初始化目录 |
| `biz/model` | 文档表结构、查询方法和搜索接口占位 |

## 数据层说明

- 当前骨架采用 `cwgo` 生成 RPC 工程
- `biz/model` 放结构体和直接查询方法，风格对齐 `ai-code` 里的 `Query` 写法
- 搜索后端当前方向固定为：
  - `Milvus`
  - `CloudWeGo Eino Retriever/Indexer`

当前 `biz/dal/milvus` 负责数据库初始化占位，和之前 `mysql` 初始化目录承担同类职责。
