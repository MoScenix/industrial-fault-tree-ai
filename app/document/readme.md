# Document

Document 微服务负责文档解析和文档搜索。

## 当前接口

- `ParsePersonalPDF`
- `ParseProjectPDF`
- `GetDocument`
- `SearchDocuments`

## 当前实现说明

- BFF 上传后先把 PDF 写到 `/document/<pdf_id>.pdf`
- 项目文档还会复制到项目目录 `documents/`
- 解析函数从 `/document/<pdf_id>.pdf` 读取
- 搜索方向固定为 `Milvus + Eino`

## 当前边界

- 不负责图管理
- 不负责项目权限编排
- 不负责前端上传页面逻辑
