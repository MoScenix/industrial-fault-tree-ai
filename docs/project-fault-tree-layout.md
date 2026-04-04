# 项目故障树目录设计

当前实现已经统一为“数据库存项目基础信息，文件系统存图版本、暂存图、建议和项目文档”。

项目目录根路径来自 `graph.root_dir`，默认值是 `/graph`。

每个项目目录使用项目创建时生成的 `uuid`：

```text
/graph/<project_uuid>/
  current
  documents/
    <pdf_id>.pdf
  versions/
    v001/
      tree.json
    v002/
      tree.json
  tmp/
    v001/
      tree.json
    v002/
      tree.json
  suggestions/
    v001.md
    v002.md
```

## 目录职责

### `current`

当前版本指针，文件内容直接保存版本号，例如：

`v002`

### `versions/<version>/tree.json`

目录版本文件。

### `tmp/<version>/tree.json`

指定版本对应的暂存编辑内容。

适用场景：

- AI 修改图
- 工程师继续编辑
- 前端优先展示暂存版本

### `suggestions/<version>.md`

指定版本的建议文件。

说明：

- AI 校验只写建议
- 不会直接改图

### `documents/<pdf_id>.pdf`

项目归档文档。

说明：

- BFF 上传项目文档时会把 PDF 复制到这里
- 文档微服务实际解析时读取全局 `/document/<pdf_id>.pdf`

## 当前约束

- 数据库不存版本列表
- 数据库不存当前建议
- 不再使用 `project.json`
- 版本列表直接扫描 `versions/`
- 当前版本通过 `current` 文件读取
- `tmp` 按版本分目录
- 建议按版本分文件
