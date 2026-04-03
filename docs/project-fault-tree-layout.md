# 项目故障树目录设计

每个项目一个目录，故障树按版本独立存放：

```text
projects/
  <project_id>/
    project.json
    fault_tree/
      current
      versions/
        v001/
          tree.json
        v002/
          tree.json
      tmpversion/
        tree.json
      last_suggestion.json
```

## 目录职责

### `project.json`

项目基础信息，至少包含：

- 项目名
- 设备名
- 顶事件
- 当前正式版本号

### `fault_tree/current`

当前正式版本指针。

文件内容直接保存版本号，例如：

`v002`

### `fault_tree/versions/v001/tree.json`

正式版本图。

### `fault_tree/tmpversion/tree.json`

AI 最近一次生成的中间版本，用于：

- 审核
- 比对
- 回撤

### `fault_tree/last_suggestion.json`

只保留最后一次建议，不保留建议历史，也不混入其他运行信息。

## 设计约束

- 正式版本和 `tmpversion` 使用完全同一份 `tree.json` 结构
- AI 写图只能写 `tmpversion`
- 正式版本的晋升由程序完成，不由 AI 直接改写
- 前端可以直接读取正式版本，也可以读取临时版本进行审核对比
