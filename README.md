# Clear Bill

Clear Bill 是一个面向账单、客户与对账流程的管理系统，当前仓库已经包含前后端基础骨架、示例接口和可运行的界面。

## 项目结构

- `mgr/server`：后端服务
- `mgr/webui`：前端应用
- `scripts`：运维与辅助脚本
- `tools`：工具脚本与检查工具

## 技术栈

前端：

- React 18
- Vite
- Ant Design 6

后端：

- Go
- Gin

## 快速开始

### 启动后端

```bash
cd mgr/server
go run ./cmd/clearbill
```

后端默认会从以下文件读取静态配置：

- `mgr/server/configs/config.toml`

支持的环境变量覆盖：

- `CLEAR_BILL_HTTP_ADDR`：服务监听地址，默认 `:8080`
- `CLEAR_BILL_WEB_ROOT`：前端静态资源目录，默认 `./website`
- `CLEAR_BILL_CONFIG`：自定义配置文件路径

### 启动前端

```bash
cd mgr/webui
pnpm install
pnpm dev
```

开发环境下，前端会将 `/api` 请求代理到本地 `http://127.0.0.1:8080`。

## 构建

构建整个项目：

```bash
make dist
```

单独构建前端：

```bash
cd mgr/webui
pnpm build
```

单独构建后端：

```bash
cd mgr/server
go build ./cmd/clearbill
```

## 发布目录

执行：

```bash
make release
```

会生成发布目录，并整理以下内容：

- 后端产物到 `release/clear-bill-<version>/server`
- 前端构建结果到 `release/clear-bill-<version>/server/website`
- 配置文件到 `release/clear-bill-<version>/server/configs`
- 脚本到 `release/clear-bill-<version>/scripts`
- 工具到 `release/clear-bill-<version>/tools`

## 当前状态

当前仓库已经完成一套可运行的基础版本：

- 后端提供了健康检查、仪表盘、账单、客户、对账等示例接口
- 后端内部已按 `action / bll / dal` 三层结构组织
- 前端提供了基于 `Ant Design 6` 的示例工作台页面
- 顶层 `Makefile` 已接入统一构建与发布入口

后续可以在这个基础上继续补充真实业务接口、页面模块和部署流程。
