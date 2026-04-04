# 文档管理微服务

文档管理微服务当前负责：

- 解析个人 PDF
- 解析项目 PDF
- 查询单个文档
- 根据问题返回搜索结果

## 当前接口

### `ParsePersonalPDF`

输入：

- `user_id`
- `pdf_id`
- `file_name`
- `display_name`

输出：

- `success`
- `document_id`
- `error_message`

### `ParseProjectPDF`

输入：

- `project_id`
- `pdf_id`
- `file_name`
- `display_name`

输出：

- `success`
- `document_id`
- `error_message`

### `GetDocument`

按 `document_id` 查询文档详情。

### `SearchDocuments`

输入：

- `user_id`
- 可选 `project_id`
- `query`
- `top_k`

输出：

- 命中的片段列表

## 当前实现约束

- BFF 上传后会先把文件写入 `/document/<pdf_id>.pdf`
- 项目文档还会额外归档到项目目录 `documents/`
- 文档解析函数默认从 `/document/<pdf_id>.pdf` 读取内容
- 解析当前按同步方式实现
- 搜索方向固定为：
  - `Milvus`
  - `CloudWeGo Eino`

## 当前数据语义

个人文档：

- `owner_type = PERSONAL`
- `owner_id = user_id`

项目文档：

- `owner_type = PROJECT`
- `owner_id = project_id`

## 搜索语义

- 只传 `user_id`
  - 搜个人文档
- 同时传 `user_id + project_id`
  - 搜指定项目上下文

当前不定义只传 `project_id` 的独立调用方式。
