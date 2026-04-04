# AI

当前 AI 微服务采用单 agent 架构。

## 当前接口

- `Chat`
  - 修改模式
  - 对话式辅助编辑
- `Validate`
  - 建议模式
  - 输出版本建议 Markdown
- `UpdatePrompt`
  - 更新 prompt 文件
- `GetPrompt`
  - 读取当前 prompt

## 当前模式

- `MODIFY_MODE`
  - 修改图
- `LOG_MODE`
  - 生成建议

## 当前 tool

- `get_project_context`
- `rag_search`
- `read_tmp_graph`
- `write_tmp_graph`

## 当前约束

- AI 只写 `tmp/<version>/tree.json`
- AI 校验只写 `suggestions/<version>.md`
- AI 不负责版本保存
- AI 不负责文档上传和解析
