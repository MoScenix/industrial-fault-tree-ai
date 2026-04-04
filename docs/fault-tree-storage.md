# 故障树存储格式

当前前端、BFF、AI 共用同一份图文件格式，主文件固定为：

`tree.json`

文件位置：

- `versions/<version>/tree.json`
- `tmp/<version>/tree.json`

## 顶层结构

- `schema_version`
- `tree`
- `nodes`
- `meta`

## 节点模型

当前图采用邻接链表语义，不使用独立边数组。

每个节点至少包含：

- `node_id`
- `node_type`
- `label`
- `description`
- `gate_type`
- `points_to`
- `pointed_by`
- 可选 `position`

## 节点类型

- `top_event`
- `intermediate_event`
- `basic_event`
- `gate`

## 逻辑门

第一阶段支持：

- `AND`
- `OR`

逻辑门作为显式节点存在，不做 children 嵌套树。

## `tree`

建议保留：

- `name`
- `top_node_id`

## `meta`

建议保留：

- `version`
- `generated_at`
- 可选 `based_on_version`

## 说明

- `points_to`
  - 表示当前节点指向哪些节点
- `pointed_by`
  - 表示哪些节点指向当前节点
- `position`
  - 允许前端直接把画布坐标写回

## 不放进图文件的内容

不要混入：

- 用户信息
- 项目路径
- 聊天记录
- 模型原始回复
- 运行日志

这些由数据库、建议文件和服务日志承担。
