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

- 读取请求版本对应的工作图文件

规则：

- 版本始终来自请求上下文，不接受模型自行指定
- 优先读该版本的 `tmp/<version>/tree.json`
- 如果该版本没有 `tmp`，则回退读取 `versions/<version>/tree.json` 作为初始内容
- 返回 `file_path`、`line_count`、`numbered_content`
- AI 必须基于返回的行号做修改

### `write_tmp_graph`

作用：

- 按行编辑 `tmp/<version>/tree.json`

规则：

- 只允许写 `tmp`
- 不直接碰目录版本
- 只支持 `insert` / `delete`
- 单次插入内容大小由提示词约束，避免模型输出过长被截断
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

## Eino ReAct 与千问的稳定性说明

当前项目使用 `cloudwego/eino` 的 ReAct agent，并通过 `eino-ext/components/model/qwen` 接入千问。

需要注意一个典型问题：

- 千问在流式场景下可能先输出少量自然语言或 reasoning
- 这会让 ReAct tool calling 更容易跑偏

当前项目的稳定性策略不是改流式判定逻辑，而是优先约束模型输出纪律：

- Prompt 层：明确要求“若下一步是工具调用，则 assistant 消息只能是工具调用本身”
- 模型层：默认关闭 thinking，并降低采样随机性，优先保证工具调用稳定性

这样可以尽量保留前端流式体验，同时减少“先说话、后调工具”的概率。

当前代码默认策略：

- `MODEL_ENABLE_THINKING=false`
- 较低温度，优先稳定执行

如需切换模型，仍然通过环境变量控制：

- `MODEL_NAME`
- `MODEL_BASE_URL`
- `DASHSCOPE_API_KEY`
- `MODEL_ENABLE_THINKING`
- `MODEL_TEMPERATURE`
- `MODEL_TOP_P`

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
