# 文档管理微服务接口

文档管理微服务只负责两类事情：

- 解析指定 PDF 为个人文档或项目文档
- 根据用户问题返回命中文档结果

这轮先定义 RPC/Proto 接口，不讨论具体存储路径和索引落盘方式。底层搜索方案当前固定为：

- `Milvus`
- `CloudWeGo Eino Retriever/Indexer`

当前工程骨架已经用 `cwgo` 生成。`biz/model` 只保留文档结构和仓储接口；搜索层预留为 `Milvus + Eino` 的接口形状，具体数据库实现后续再挂载。

## 服务

- `DocumentService`

## RPC

### `ParsePersonalPDF`

用途：

- 把指定 PDF 解析为某个用户的个人文档

请求字段：

- `user_id`
- `pdf_id`
- `file_name`
- `display_name`

响应字段：

- `success`
- `document_id`
- `error_message`

说明：

- 上游已经完成文件上传
- 这里只接收 PDF 标识和基础文件信息
- 解析接口只返回是否成功和生成出的 `document_id`

### `ParseProjectPDF`

用途：

- 把指定 PDF 解析为某个项目的项目文档

请求字段：

- `project_id`
- `pdf_id`
- `file_name`
- `display_name`

响应字段：

- `success`
- `document_id`
- `error_message`

### `GetDocument`

用途：

- 按 `document_id` 查询单个文档详情

请求字段：

- `document_id`

响应字段：

- `document`

### `ListDocuments`

用途：

- 按归属列出文档

请求字段：

- `owner_type`
- `owner_id`
- `page`
- `page_size`

`owner_type` 固定支持：

- `PERSONAL`
- `PROJECT`

响应字段：

- `documents[]`
- `total`

### `SearchDocuments`

用途：

- 根据传入问题直接返回搜索命中结果

请求字段：

- `user_id`
- `project_id`
- `query`
- `top_k`

语义：

- 只传 `user_id` 时，搜索用户个人文档范围
- 同时传 `user_id` 和 `project_id` 时，搜索该用户在指定项目上下文中的结果
- 当前不定义只传 `project_id` 不传 `user_id` 的调用方式

响应字段：

- `results[]`

每条结果至少包含：

- `document_id`
- `document_name`
- `chunk_id`
- `text`
- `score`

## 核心类型

### `Document`

- `document_id`
- `owner_type`
- `owner_id`
- `pdf_id`
- `file_name`
- `display_name`
- `parse_status`
- `summary`
- `chunks[]`
- `created_at`
- `updated_at`

### `DocumentChunk`

- `chunk_id`
- `text`
- `page`
- `order`

### `SearchResult`

- `document_id`
- `document_name`
- `chunk_id`
- `text`
- `score`

## 约束

- 解析接口按同步完成定义
- 解析成功后，调用方通过 `document_id` 再查详情或继续搜索
- 搜索接口返回命中片段，不只返回文档列表
- 当前接口层不暴露底层检索实现细节
- 当前搜索实现方向固定为 `Milvus + Eino`
- 当前接口层不包含故障树字段
- 当前接口层不包含文件系统路径字段

## 覆盖场景

- `ParsePersonalPDF` 成功返回 `success=true` 和 `document_id`
- `ParseProjectPDF` 成功返回 `success=true` 和 `document_id`
- 解析失败返回 `success=false` 和 `error_message`
- `GetDocument` 返回文档详情和 `chunks`
- `ListDocuments` 支持按 `PERSONAL + user_id` 查询
- `ListDocuments` 支持按 `PROJECT + project_id` 查询
- `SearchDocuments` 在只传 `user_id` 时返回个人文档命中结果
- `SearchDocuments` 在传 `user_id + project_id` 时返回项目上下文命中结果
