# industrial-fault-tree-ai

工业设备故障树智能生成与辅助分析平台。

项目仓库：

- https://github.com/MoScenix/industrial-fault-tree-ai

## 系统概览

当前仓库采用微服务 + BFF + 前端的结构：

- `app/user`
  - 用户注册、登录、用户信息管理
- `app/graph`
  - 图项目元信息、聊天记录、版本管理、`tmp` 生命周期
- `app/document`
  - 个人/项目 PDF 解析、文档搜索
- `app/ai`
  - 单 agent AI 服务，支持修改模式与建议模式
- `app/bff`
  - 对前端暴露聚合接口，负责权限校验与服务编排
- `app/frontend`
  - Vue 3 + Ant Design Vue + Vue Flow 前端

## 当前业务规则

- 数据库只保存：
  - 项目基础信息
  - 聊天记录
- 故障树版本、暂存版本、建议文件全部走文件系统
- 普通用户只能操作自己的项目
- `admin` 可以查看和操作全部项目
- AI 修改图时只写 `tmp/<version>/tree.json`
- AI 校验只生成 `suggestions/<version>.md`，不改图
- 工程师保存时：
  - `useTmp=true`：直接提交后端已有暂存图
  - `useTmp=false`：先把前端当前内容写入 `tmp/<fromVersion>/tree.json`，再提交

## 项目目录

每个项目使用独立 `uuid` 目录，根目录默认 `/graph`：

```text
/graph/<project_uuid>/
  current
  documents/
    *.pdf
  versions/
    v001/
      tree.json
  tmp/
    v001/
      tree.json
  suggestions/
    v001.md
```

说明：

- `current`
  - 当前版本指针
- `versions/<version>/tree.json`
  - 目录版本
- `tmp/<version>/tree.json`
  - 当前编辑态
- `suggestions/<version>.md`
  - AI 校验建议
- `documents/*.pdf`
  - 项目归档文档

## 前端页面

- `/`
  - 首页导航
- `/graph/manage`
  - 我的项目
- `/graph/workspace/:id`
  - 图工作台
- `/admin/graphManage`
  - 管理员项目总览
- `/admin/userManage`
  - 用户管理
- `/admin/promptManage`
  - 提示词管理

工作台说明：

- 左侧三按钮：
  - 对话
  - 版本
  - 校验
- 中间主区域：
  - Vue Flow 故障树编辑器
- 右上角操作：
  - 保存
  - 准备暂存
  - 上传项目文档
  - 导入图
  - 导出当前图
  - 放弃暂存

## 运行依赖

根目录 `docker-compose.yml` 当前包含：

- MySQL
- Redis
- Consul
- Etcd
- MinIO
- Milvus
- Jaeger

详细启动方式见：

- `/root/Projects/tmp/industrial-fault-tree-ai/docs/run-local.md`

## 核心文档

- `/root/Projects/tmp/industrial-fault-tree-ai/docs/project-fault-tree-layout.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/docs/fault-tree-storage.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/docs/ai-service-tools.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/docs/document-service.md`
- `/root/Projects/tmp/industrial-fault-tree-ai/docs/business-flow.md`

## 当前已知限制

- `ChatToModifyGraph` 当前是“一次请求 -> 聚合 AI 回复 -> 成功后落消息”，还不是真正的流式 SSE UI
- AI 校验当前输出的是建议 Markdown，不是严格图结构校验器
- 工作台底部 JSON 区当前作为调试模式保留，默认隐藏
