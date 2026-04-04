# AI 服务 Tool 说明

当前 AI 服务采用单 agent 架构，只有两种模式：

- `MODIFY_MODE`
  - 用于对话式修改图
- `LOG_MODE`
  - 用于生成建议和校验说明

两种模式共用同一套 agent 工厂，差异只有：

- system prompt
- 使用场景

## 当前 Tool

### `get_project_context`

作用：

- 读取项目运行上下文

返回信息：

- `project_id`
- `current_version`
- `tmp_version_ready`
- `document_summary`

### `rag_search`

作用：

- 在当前项目的已解析文档中检索相关内容

输入：

- `query`
- `top_k`
- `filters`

返回：

- `chunk_id`
- `document_name`
- `text`
- `score`

### `read_tmp_graph`

作用：

- 读取当前工作图

规则：

- 指定 `version` 时，优先读该版本的 `tmp/<version>/tree.json`
- 如果该版本没有 `tmp`，则读取 `versions/<version>/tree.json`
- 不传 `version` 时，默认读取 `current` 指向的版本

### `write_tmp_graph`

作用：

- 将 AI 整理好的完整图写回 `tmp/<version>/tree.json`

规则：

- 只允许写 `tmp`
- 不直接碰目录版本
- 写回后由工程师决定是否保存

## 当前 AI 业务语义

- `Chat`
  - 使用 `MODIFY_MODE`
  - 允许读写 `tmp`
- `Validate`
  - 使用 `LOG_MODE`
  - 只生成建议 Markdown
  - 建议文件写入 `suggestions/<version>.md`
- `UpdatePrompt`
  - 只改 prompt 文件
  - 当前仅支持：
    - `MODIFY_MODE`
    - `LOG_MODE`

## 不属于 AI Tool 的职责

下面这些能力不属于 AI tool：

- 新建项目目录
- 版本创建 / 删除 / 重命名
- 项目保存
- 上传文档
- 文档解析
- 用户权限控制
- 提示词页面编排

这些分别由：

- Graph 微服务
- Document 微服务
- BFF
- Frontend

共同负责。
