# 故障树存储建议

主文件固定为：

`projects/<project_id>/fault_tree/.../tree.json`

推荐结构：

- `schema_version`
- `tree`
- `nodes`
- `edges`
- `evidence`
- `meta`

## 关键约束

- 使用 `nodes + edges`，不做 `children` 嵌套树
- 逻辑门显式建模为节点
- 节点类型固定为：
  - `top_event`
  - `intermediate_event`
  - `basic_event`
  - `gate`
- 第一阶段 gate 只支持：
  - `AND`
  - `OR`
- `evidence` 只引用 RAG 返回的 `chunk_id`
- `meta` 只保留必要图信息：
  - `version`
  - `based_on_version`
  - `generated_at`
  - `source_chunk_ids`

## 不建议写入的内容

不要在 `tree.json` 中混入：

- 用户信息
- 项目路径
- 运行日志
- 模型原始回复
- 任务过程状态
