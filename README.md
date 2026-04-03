# industrial-fault-tree-ai

基于知识的工业设备故障树智能生成与辅助构建系统。

当前仓库的方向已调整为：

- `user` 模块继续复用
- `app/ai` 已重新建立为新的 AI 微服务骨架
- AI 与文档服务优先先落接口、IDL 和目录约束
- 先统一 tool 边界、项目目录、故障树文件格式和文档服务接口

当前依赖根路径为：

`github.com/MoScenix/industrial-fault-tree-ai`

简历版项目介绍见：

- `/root/Projects/tmp/industrial-fault-tree-ai/docs/resume-project.md`

目前优先落地的内容：

- `app/ai`：新的 AI 微服务骨架与内部 tool 占位
- `app/user`：直接复用的用户微服务
- `idl`：当前包含 `user.proto`、`ai.proto`、`document.proto`
- `common`：公共基础能力
- `rpc_gen`：保留用户模块服务端生成代码与 client 调用封装
- `docs/ai-service-tools.md`：AI 服务 tool 设计
- `docs/document-service.md`：文档管理微服务接口
- `docs/project-fault-tree-layout.md`：项目目录设计
- `docs/fault-tree-storage.md`：故障树存储规范
- `specs/`：`project.json`、`tree.json`、`last_suggestion.json` 的 schema
- `examples/project-template/`：程序初始化目录参考模板

建议程序和前端优先对齐下面几份文件：

- `/root/Projects/tmp/industrial-fault-tree-ai/docs/ai-service-tools.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/docs/document-service.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/docs/project-fault-tree-layout.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/specs/fault-tree.schema.json`
