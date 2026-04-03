# graph

图管理微服务当前先落两层基础能力：

- 图项目的常规信息管理
- 图项目聊天记录的存储与查询

## 当前接口职责

- `AddGraph`
  - 新建图项目元信息
  - 同时创建项目目录、`versions/`、`tmp/` 和 `current`
- `UpdateGraph`
  - 更新图项目基础信息
- `GetGraph`
  - 查询单个图项目基础信息
- `ListGraph`
  - 分页列出图项目
- `DeleteGraph`
  - 删除图项目
  - 当前产品约定为数据库和项目目录一起删
- `AddGraphMessage`
  - 新增聊天记录
- `ListGraphMessage`
  - 按 `graph_id + last_create_time` 查询聊天记录
- `StartEdit`
  - 进入编辑态，准备或复用 `tmp`
- `Save`
  - 把某个 `from_version` 的 `tmp` 提交到某个 `to_version`
- `CreateVersion / DeleteVersion / RenameVersion / ListVersion`
  - 管理正式版本

## 数据与目录约定

- 图项目目录根路径来自配置项 `graph.root_dir`
- 当前默认值是 `/graph`
- 每个图项目目录名使用 `uuid`
- 项目目录最终形态是：
  - `/graph/<uuid>`

## 当前 DAO 层

`biz/model` 现在先实现两类数据：

- `Graph`
  - 存图项目元信息
  - 包含项目 `uuid` 和目录路径
- `Message`
  - 存图项目聊天记录
  - 支持基于 `last_create_time` 的深分页查询

MySQL 初始化时会自动建表：

- `graphs`
- `messages`

## utils 约定

`utils` 当前承接的是会被多个 service 复用、但不应该留在 service 里的小工具：

- `project.go`
  - 生成项目 `uuid`
  - 计算项目目录路径
  - 初始化项目目录结构
  - 维护 `versions/`、`tmp/`、`current`
  - 列出版本目录
- `time.go`
  - 统一时间格式化
  - 解析 `last_create_time` 游标

## 当前微服务边界

`graph` 微服务当前更适合负责：

- 图项目元信息管理
- 项目目录初始化
- 聊天记录查询
- 版本管理与 `tmp` 生命周期管理

前端高频编辑图内容这件事，当前仍然建议放 API 层直接写 `tmp`，而不是把每次拖拽和改点都拆成 RPC。
