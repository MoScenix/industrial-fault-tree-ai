# Frontend 运行说明

## 环境变量

前端本地默认读取：

- `app/frontend/.env`

当前默认值：

```env
VITE_API_BASE_URL=http://localhost:8082/api
VITE_DEPLOY_DOMAIN=http://localhost
```

## 页面结构

- `/`
  - 首页导航
- `/graph/manage`
  - 我的项目
- `/graph/workspace/:id`
  - 故障树工作台
- `/admin/graphManage`
  - 管理员项目总览
- `/admin/promptManage`
  - 提示词管理
- `/admin/userManage`
  - 用户管理

## 工作台

当前工作台已接入 `Vue Flow`，并包含：

- 画布节点/边展示
- 对话侧栏
- 版本侧栏
- 建议侧栏
- 导入图
- 上传项目文档
- 导出当前图
- 保存当前内容

## 启动

```bash
pnpm install
pnpm dev
```

如果前端安装依赖遇到权限或网络问题，需要允许终端在当前环境下联网安装依赖。
