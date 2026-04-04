# 业务流程说明

## 1. 登录

- 用户通过 BFF 登录
- BFF session 中保存：
  - `user_id`
  - `user_role`
- 普通用户只能操作自己的项目
- `admin` 可操作全部项目

## 2. 新建项目

- 前端调用 `AddGraph`
- Graph 微服务：
  - 创建数据库记录
  - 生成项目 `uuid`
  - 初始化 `/graph/<uuid>`
  - 初始化 `versions/v001`
  - 初始化 `current`

## 3. 进入工作台

- 前端先调用 `GetGraphVOById`
- 再调用 `GetWorkingGraph`
- 如果指定版本存在 `tmp/<version>/tree.json`
  - 优先展示暂存图
- 否则展示 `versions/<version>/tree.json`

## 4. 准备暂存

- 前端调用 `StartEdit`
- Graph 微服务：
  - 如果 `tmp/<version>` 已存在则直接复用
  - 否则从 `versions/<version>` 复制到 `tmp/<version>`

## 5. AI 对话改图

- 前端调用 BFF 的 `ChatToModifyGraph`
- BFF：
  - 读取最近消息上下文
  - 调 AI 微服务 `Chat`
  - 聚合本次回复
  - 成功后写入两条消息：
    - user
    - assistant
- AI 修改模式只允许读写 `tmp`

## 6. 工程师手动改图

- 当前编辑内容先在前端内存中变化
- 保存时：
  - 如果当前展示的是 AI 已写好的 `tmp`，可直接 `useTmp=true`
  - 如果工程师改动仍在前端内存，则 `useTmp=false` 并上传 `content`
- BFF 会先把 `content` 写入 `tmp/<fromVersion>/tree.json`
- 再调用 Graph 微服务执行 `from -> to` 保存

## 7. 校验与建议

- 前端调用读取当前建议接口
- BFF 返回 `suggestions/<version>.md`
- AI `Validate` 当前只生成建议内容，不改图

## 8. 保存版本

- Graph 微服务接收 `Save(fromVersion, toVersion)`
- 读取 `tmp/<fromVersion>`
- 将内容复制到 `versions/<toVersion>`
- 更新 `current`
- 删除 `tmp/<fromVersion>`

## 9. 上传个人文档

- 前端从导航栏上传
- BFF 校验登录态
- 文件写入全局 `/document/<pdf_id>.pdf`
- 再调用 Document 微服务 `ParsePersonalPDF`

## 10. 上传项目文档

- 前端在工作台上传
- BFF 校验 owner/admin
- 文件同时写入：
  - `/document/<pdf_id>.pdf`
  - `/graph/<uuid>/documents/<pdf_id>.pdf`
- 再调用 Document 微服务 `ParseProjectPDF`

## 11. 管理员提示词管理

- 只有 `admin` 可见
- 当前只维护两种模式：
  - `MODIFY_MODE`
  - `LOG_MODE`
- 修改后直接写回 AI 服务本地 prompt 文件
