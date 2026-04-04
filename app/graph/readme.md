# Graph

Graph 微服务负责图项目基础管理和版本管理。

## 当前职责

- 项目元信息管理
- 项目目录初始化
- 聊天记录存储与查询
- `tmp` 生命周期管理
- 版本创建、删除、重命名、列出
- 保存 `from_version -> to_version`

## 当前数据边界

数据库只保存：

- 项目基础信息
- 聊天记录

文件系统保存：

- `versions/<version>/tree.json`
- `tmp/<version>/tree.json`
- `suggestions/<version>.md`
- `documents/*.pdf`

## 当前规则

- 不再使用 `project.json`
- 当前版本通过 `current` 文件读取
- 版本列表通过扫描 `versions/` 获取
