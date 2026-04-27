# Clear Bill

Clear Bill 是一个面向账单、客户、对账与权限治理场景的管理系统仓库。当前仓库已经包含可运行的后端服务、管理端前端、基础构建流程，以及一套围绕租户、用户、角色权限和账单领域的示例接口。

如果你是第一次接手这个项目，建议先看这份 README，再进入 `mgr/server` 和 `mgr/webui` 分别查看后端与前端细节。

## 项目概览

- 后端：基于 `Go + Gin + GORM + MySQL`
- 前端：基于 `React 18 + Vite + Ant Design 6`
- 组织方式：单仓库维护管理端前后端与辅助脚本
- 当前定位：可继续扩展业务能力的基础版本，已经具备开发、构建和发布打包入口

## 当前能力

后端已提供以下模块的接口或骨架能力：

- 系统健康检查
- 登录、登出、当前用户、修改密码、创建 API Token
- 租户 CRUD
- 用户 CRUD 与重置密码
- 角色 CRUD、权限列表、角色权限分配
- 仪表盘概览、账单列表、客户列表、对账任务列表

前端当前已提供：

- 登录页
- 管理台布局骨架
- 首页、租户、角色权限、凭证、审计日志等导航结构
- 基于 `/api` 的本地开发代理配置

## 目录结构

```text
.
├─ mgr/
│  ├─ server/   # Go 后端服务
│  └─ webui/    # React 管理端
├─ scripts/     # 运维与辅助脚本
├─ tools/       # 诊断、检查与打包辅助工具
├─ Makefile     # 仓库级统一入口
├─ build.sh     # 构建包装脚本
└─ deploy.sh    # 发布脚本占位
```

## 环境要求

建议本地准备以下环境：

- `Go 1.22+`
- `Node.js 20+`
- `pnpm 10.17.1`
- `MySQL 8.x`
- `make`
- Bash 兼容 Shell

说明：

- `mgr/server/go.mod` 当前声明的 Go 版本为 `1.22.0`
- 前端 `package.json` 当前固定使用 `pnpm@10.17.1`
- 根目录下的 `build.sh`、`deploy.sh` 与若干 Makefile 目标使用了 POSIX 风格命令，Windows 环境下建议通过 Git Bash、WSL 或等效环境执行

## 快速开始

### 1. 准备后端配置

后端默认读取：

- `mgr/server/configs/config.toml`

其中 `[GORM]` 段用于配置 MySQL 连接。启动前请按本地环境修改数据库地址、库名、用户名和密码。

支持的环境变量覆盖：

- `CLEAR_BILL_CONFIG`：自定义配置文件路径
- `CLEAR_BILL_HTTP_ADDR`：覆盖监听地址，例如 `127.0.0.1:8080`
- `CLEAR_BILL_WEB_ROOT`：覆盖静态站点目录

### 2. 启动后端

```bash
cd mgr/server
go run ./cmd/clearbill
```

也可以从仓库根目录启动：

```bash
make server-start
```

后端默认监听 `0.0.0.0:8080`，并提供健康检查接口：

```text
GET /health
GET /api/v1/health
```

### 3. 启动前端

```bash
cd mgr/webui
pnpm install
pnpm dev
```

开发环境下，前端会将 `/api` 请求代理到：

```text
http://127.0.0.1:8080
```

如需调整代理目标，可设置这些环境变量：

- `VITE_API_PROXY_TARGET`
- `VITE_API_PROXY_CHANGE_ORIGIN`
- `VITE_API_PROXY_SECURE`
- `VITE_API_PROXY_WS`

## 常用命令

### 仓库根目录

```bash
make dist
make release
make test
make lint
```

说明：

- `make dist`：统一构建子项目
- `make release`：整理发布目录到 `release/clear-bill-<version>`
- `make test`：触发各子项目测试命令
- `make lint`：触发各子项目静态检查命令

### 后端

```bash
cd mgr/server
make build
make test
make lint
make wire
make swagger
```

说明：

- `make build`：输出二进制到 `mgr/server/dist/`
- `make wire`：重新生成依赖注入代码
- `make swagger`：重新生成 `api/docs/` 下的 Swagger 文档

### 前端

```bash
cd mgr/webui
pnpm dev
pnpm build
pnpm lint
pnpm api
```

说明：

- `pnpm api`：基于后端 Swagger 文档重新生成前端请求代码

## 构建与发布

### 构建

统一构建：

```bash
make dist
```

构建结果主要位于：

- `mgr/server/dist/`
- `mgr/webui/dist/`

### 发布打包

```bash
make release
```

该命令会生成：

```text
release/clear-bill-<version>/
```

目录内容包括：

- 后端产物：`server/`
- 前端静态资源：`server/website/`
- 后端配置：`server/configs/`
- 辅助脚本：`scripts/`
- 工具脚本：`tools/`

## 开发建议

- 先启动后端，再启动前端，便于直接联调 `/api`
- 调整接口后，优先执行 `cd mgr/server && make swagger`，再执行 `cd mgr/webui && pnpm api`
- 当前 `deploy.sh` 仍是占位脚本，只提示先完成构建，正式部署流程还需要补齐
- `scripts/` 与 `tools/` 当前保持轻量，后续可以继续沉淀环境初始化、巡检、打包校验等能力

## 相关入口

- 后端入口：`mgr/server/cmd/clearbill/main.go`
- 路由注册：`mgr/server/api/router/`
- 前端路由：`mgr/webui/src/routes.tsx`
- 前端代理配置：`mgr/webui/config/proxy.ts`

