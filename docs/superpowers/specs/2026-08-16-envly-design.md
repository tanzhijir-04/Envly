# Envly — 一键开发环境配置工具 设计文档

- 日期：2026-08-16
- 状态：待评审
- 仓库：`Envly`（Windows 优先，macOS / Linux 架构预留）

## 1. 项目概述

Envly 是一款跨平台桌面小工具，帮助开发者一键配置开发环境与系统环境（开发相关部分）。UI 采用 Apple Liquid Glass 风格，支持中英双语。Windows 为 v1 目标平台，macOS / Linux 为 v2（架构上预留）。

### 目标用户

- 帮朋友 / 新手一键配好环境（操作者不熟悉命令行也能用）
- 开发者自己在多台机器上快速复现同一套环境

### 核心体验

单页向导：选择模板 → 确认 / 修改勾选清单 → 一键配置（实时日志）→ 验证报告。

## 2. 范围与非目标

### 范围内

- 常用开发工具安装：运行时、AI CLI、编辑器 / 工具、终端增强
- 系统环境配置：包管理器、PATH、Shell 配置、镜像源、代理、GitHub 加速
- 预设模板 + 手动勾选
- 中英双语 UI
- 网络区域自适应（中国大陆 / 海外）

### 非目标（v1 不做）

- 系统体验类设置（外观、输入法、电源等）
- VS Code 扩展安装
- Shell 插件（第二期）
- 应用自动更新机制

## 3. 平台与工具清单

### v1：Windows

安装方式：winget / npm / pip / download / rustup

| 分类 | 工具 |
| --- | --- |
| 基础运行时 | Node.js LTS、Git、Python 3、TypeScript、C/C++（MinGW 默认勾选，MSVC 可选）、PHP、uv、Jupyter Notebook、Go（可选）、Rust（可选）、Java Temurin（可选） |
| AI CLI | Claude Code、Codex CLI、Gemini CLI、Pake（依赖 Rust） |
| 编辑器 / 工具 | VS Code、Cursor、Sublime Text、draw.io、JetBrains Toolbox（可选） |
| 终端增强 | Windows Terminal、PowerShell 7、Oh My Posh、Cascadia Code NF |
| 系统环境 | winget 包管理器、用户 PATH 合并、PowerShell Profile、npm / pip 镜像、GitHub 加速 / 代理检测 |

### v2：macOS

安装方式：brew / brew --cask / npm / pip / rustup，Xcode Command Line Tools。

- 运行时：Node.js、Git、Python 3、TypeScript、C/C++（Xcode CLT + brew gcc）、PHP、uv、Jupyter、Go / Rust / Java（可选）
- AI CLI：同 Windows
- 编辑器 / 工具：VS Code、Cursor、Sublime Text、draw.io、JetBrains（可选）
- 终端增强：iTerm2 / Warp、Oh My Zsh / fish、Nerd Font
- 系统环境：Homebrew、.zshrc、镜像源、代理检测

### v2：Linux

安装方式：apt（Ubuntu / Debian 优先）→ dnf（Fedora）/ pacman（Arch）按优先级后续加。

- 运行时：nvm / 官方源 Node.js、Git、Python 3 + pyenv、TypeScript、gcc / clang + build-essential、PHP、uv、Jupyter
- AI CLI：同 Windows
- 编辑器 / 工具：微软源 / snap 的 VS Code、deb / rpm 的 Cursor、Sublime Text、draw.io
- 终端增强：kitty / Alacritty、zsh / fish + Starship
- 系统环境：发行版包管理器、.bashrc / .zshrc、镜像源、代理检测

### 模板

- AI 开发环境（默认推荐）
- 前端开发
- 后端开发
- 极简起步

每个模板是工具 ID 列表；用户可勾选增删；可选工具（Go / Rust / Java / JetBrains / MSVC 等）默认不勾选。

## 4. 架构

Pake 壳 + Go 引擎（双进程）。

```
┌─ Envly-UI.exe（Pake 壳，WebView 窗口）
│   加载 http://127.0.0.1:17521
└─ Envly.exe（Go 引擎，本地 HTTP 服务）
    ├─ 静态资源（前端 SPA）
    ├─ /api/*（REST + SSE 日志流）
    └─ 安装 / 环境 / 网络 / 验证模块
```

- 启动器：双击 Envly.exe → 启动 HTTP 服务 → 拉起 Envly-UI.exe → 心跳监测窗口，窗口关闭后引擎退出
- 端口冲突：检测 17521 被占用则顺延 +1，通过启动参数同步给 UI
- 打包：构建机安装 Node 18+、Rust 工具链、pake-cli（`npm i -g pake-cli`），执行 pake 命令产出 msi（Windows）/ dmg（macOS）/ deb（Linux）
- 备选：若后续希望单文件分发，可迁移到 Wails（Go + WebView），UI 与引擎代码不变

### API 设计（草案）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/health | 健康检查 |
| GET | /api/catalog | 工具清单（平台过滤、i18n 文案） |
| GET | /api/templates | 模板列表 |
| POST | /api/plan | 生成执行计划（标记已装 / 待装 / 不支持） |
| POST | /api/run | 开始执行 |
| GET | /api/events | SSE 推送结构化日志 |
| POST | /api/cancel | 取消执行 |
| GET | /api/report | 最近一次执行报告 |
| GET/POST | /api/settings | 语言、网络区域、镜像开关 |

## 5. 组件划分

- `engine/config`：数据驱动的工具清单与模板（每条记录含 id、中英文名称 / 描述、分组、各平台安装方式、验证命令）
- `engine/installer`：安装执行器（winget / npm / pip / download / brew / apt / rustup），从 WinDevReady 迁移
- `engine/env`：PATH 合并、PowerShell Profile 写入、npm / pip 镜像配置、代理配置
- `engine/network`：连通性检测、区域判定、GitHub 加速
- `engine/verify`：安装后验证命令执行
- `engine/store`：安装记录（JSON，含“由 Envly 安装”标记，供卸载清理）
- `ui`：单页向导前端（模板 → 清单 → 执行 → 日志 → 报告），中英双语，Apple Liquid Glass 风格

## 6. 数据流

1. 启动 → health check → 读取设置（语言、网络区域）
2. 网络区域检测：探测 github.com / registry.npmmirror.com 等端点连通性
3. 前端拉取 catalog + templates，渲染向导
4. 用户选模板 / 勾选 → POST /api/plan → 引擎返回计划（每项状态：已装 / 待装 / 不支持）
5. 用户点“一键配置”→ POST /api/run → 引擎顺序执行：
   - 已装 → 跳过并记录版本
   - 网络类失败 → 按策略重试（最多 2 次）
   - 单项失败 → 标记失败，继续下一项
6. SSE 推送事件（工具 / 状态 / 消息 key + 参数），UI 本地化渲染实时日志
7. 执行结束 → 逐项验证 → 写入 store → 生成报告 → UI 展示报告卡片

## 7. 区域与镜像策略

- 默认自动检测：检测 GitHub / npm 官方与国内镜像的连通性
  - 大陆网络（官方不可达、镜像可达）→ 默认启用国内镜像（npm 用 npmmirror，pip 用清华源，GitHub 走加速）
  - 海外网络 → **不修改镜像源**，仅保留代理检测（读取系统代理并配置 git / npm）
- 手动覆盖：设置页提供“网络区域：自动 / 中国大陆 / 海外”
- 还原能力：切换镜像前记录原始 registry / pip.conf，设置或卸载时可一键还原
- 所有镜像 / 代理改动写入 store，可审计、可回滚

## 8. 错误处理

- 单项失败不中断整体流程
- 网络类错误：重试 2 次（退避 3s / 8s），仍失败则标记并继续
- 下载失败：按 DownloadURLs 列表回退（官方 / 镜像交替）
- 引擎端口占用：顺延端口
- UI 与引擎失联：UI 显示重连提示；引擎心跳检测窗口存活，窗口关闭则引擎退出
- 安装记录损坏：启动时备份为 .bak，不阻塞主流程

## 9. 测试

- Go 单元测试：catalog 解析、模板展开、PATH 合并、区域判定、验证命令输出解析
- 前端组件测试：模板选择、勾选状态、日志渲染、双语切换
- 集成测试：Windows runner 真实安装最小集（Node.js LTS、Git、镜像切换），安装后验证
- CI：GitHub Actions 三平台构建 + Windows 冒烟测试

## 10. 里程碑

- M1 引擎骨架：HTTP 服务、catalog / templates / plan API、前端骨架（向导三屏 + 日志流）
- M2 Windows 全流程：installer / env / network / verify 模块接入，真实安装可用
- M3 收尾：验证报告、双语完善、设置页、Pake 打包与启动器、CI 构建
- M4 跨平台：macOS（brew）→ Linux（apt）适配，按 v2 清单

## 11. 仓库布局

```
Envly/
├─ engine/            # Go 引擎（cmd + internal）
│  ├─ cmd/envly/      # 引擎入口
│  └─ internal/       # config/installer/env/network/verify/store/api
├─ ui/                # 前端源码（HTML/CSS/TS，构建后由引擎内嵌）
├─ pake/              # Pake 打包脚本与图标资源
├─ scripts/           # 启动器 / 构建脚本
├─ docs/superpowers/specs/
├─ .github/workflows/ # 三平台 CI
└─ README.md
```

## 12. 已确认决策记录

- 布局：B 单页向导
- 打包：Pake（双进程）；备选 Wails
- 仓库名：Envly
- 语言：中英双语
- 区域：海外不强制改镜像
- C/C++：MinGW 默认、MSVC 可选
- Linux：apt（Ubuntu / Debian）优先
- 交互：模板 + 手动勾选
