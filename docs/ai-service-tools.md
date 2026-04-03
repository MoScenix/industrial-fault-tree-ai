# AI 服务 Tool 设计

AI 服务只消费“程序已经准备好的项目上下文”和“文档微服务已经完成的解析/索引结果”。

以下 tool 都由 AI 服务内部调用，但项目上下文由业务程序注入，AI 本身不传 `project_id`、用户信息、目录路径。

## Tool 列表

### 1. `get_project_context`

用途：

- 读取当前项目上下文

AI 入参：

- 无

程序注入上下文：

- `project_id`
- `device_name`
- `top_event`
- `current_version`
- `available_versions`

返回：

- `device_name`
- `top_event`
- `current_version`
- `tree_summary`
- `document_summary`

### 2. `rag_search`

用途：

- 在当前项目已解析文档中做检索

AI 入参：

- `query`
- `top_k`
- `filters`

程序隐式约束：

- 只能搜索当前项目
- 只能搜索已解析完成的文档索引

返回：

- `chunks[]`

每个 chunk 至少包含：

- `chunk_id`
- `document_name`
- `text`
- `score`

### 3. `read_fault_tree`

用途：

- 读取当前图或指定历史版本

AI 入参：

- 可选 `version`

默认行为：

- 不传时读取当前正式版本

返回：

- `tree`
- `version`
- `is_tmp_version`

### 4. `write_tmp_fault_tree`

用途：

- 写 AI 生成或修改后的中间图

AI 入参：

- `tree`
- `change_summary`

行为：

- 永远不覆盖正式版本
- 始终写入 `fault_tree/tmpversion/tree.json`
- 覆盖上一次临时版本

返回：

- `tmp_version_path`
- `based_on_version`

### 5. `write_last_suggestion`

用途：

- 保存 AI 的最后一次建议

AI 入参：

- `suggestions`

行为：

- 只保留最后一次建议
- 每次写入直接覆盖 `fault_tree/last_suggestion.json`

返回：

- `suggestion_path`

## 不属于 AI Tool 的职责

下面这些能力明确不进入 AI tool：

- 新建项目目录
- 文档上传
- 文档解析
- 文档切片
- 向量化建库
- 正式版本晋升
- 版本切换
- 权限控制

这些由业务程序或文档微服务负责。
