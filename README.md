# Clear Bill

`Clear-Bill` 的目录组织参考了 `../lxopeasier`，但前端技术栈已经按本项目要求切换为 `Ant Design 6`。

## 项目结构

- `mgr/server`：后端服务
- `mgr/webui`：前端应用
- `scripts`：运维与辅助脚本
- `tools`：工具脚本与检查工具

## 技术说明

前端保留了参考项目的目录位置，也就是仍然放在 `mgr/webui` 下，但没有继续沿用旧的 Umi 方案，而是重建为更适合当前项目的组合：

- React 18
- Vite
- Ant Design 6

这样可以在保持整体目录架构一致的同时，避免沿用过旧的前端运行时和配置方式。

## 快速开始

### 启动后端

```bash
cd mgr/server
go run ./cmd/clearbill
```

默认监听地址可通过环境变量配置：

- `CLEAR_BILL_HTTP_ADDR`：服务监听地址，默认 `:8080`
- `CLEAR_BILL_WEB_ROOT`：前端静态资源目录，默认 `./website`

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

如果只想构建前端：

```bash
cd mgr/webui
pnpm build
```

如果只想构建后端：

```bash
cd mgr/server
go build ./cmd/clearbill
```

## 发布目录

执行：

```bash
make release
```

会生成类似参考项目的发布目录，并整理以下内容：

- 后端产物到 `release/clear-bill-<version>/server`
- 前端构建结果到 `release/clear-bill-<version>/server/website`
- 配置文件到 `release/clear-bill-<version>/server/configs`
- 脚本到 `release/clear-bill-<version>/scripts`
- 工具到 `release/clear-bill-<version>/tools`

## 当前状态

当前仓库已经完成一套可运行的基础骨架：

- 后端提供了基础示例接口
- 前端提供了基于 `Ant Design 6` 的示例工作台页面
- 顶层 `Makefile` 已接入统一构建与发布入口

后续可以在这个基础上继续补充真实业务接口、页面模块和部署流程。
