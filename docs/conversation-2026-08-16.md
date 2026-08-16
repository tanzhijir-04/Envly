# Envly 对话记录（2026-08-16）

本文件记录 Envly 从想法到设计确认的完整对话过程，供后续实现与追溯使用。

## 背景

用户希望做一个跨平台（MAC / WIN / LINUX）的一键开发环境配置桌面小工具，参考其已有项目 [WinDevReady](https://github.com/tanzhijir-04/WinDevReady)（Go + Fyne，Windows 专用，安装 Node / Git / Python / AI CLI / 编辑器 / 终端增强，含网络检测与镜像切换）。

## 需求确认（一问一答记录）

1. **目标用户**：A + B
   - 帮朋友 / 新手一键配好环境（操作者可以是懂行的人）
   - 开发者自己在多台机器上快速复现同一套环境
2. **“系统环境”范围**：A —— 开发相关的环境配置（PATH、Shell 配置、镜像源、代理、包管理器）；不含系统体验类设置（外观、输入法等）
3. **v1 平台策略**：C —— 先做 Windows（延续 WinDevReady 经验），再铺 macOS / Linux
4. **工具范围**：A/B/C/D/E 全选，并补充：
   - 新增：TypeScript、C/C++、PHP、uv、Jupyter Notebook、Pake、draw.io、Sublime Text
   - 移除：VS Code 扩展（Shell 插件保留到第二期）
5. **C/C++ 工具链（Windows）**：C —— MinGW 默认勾选，MSVC Build Tools 作为可选项
6. **Linux 发行版范围**：A —— Ubuntu / Debian 系优先（apt），Fedora（dnf）、Arch（pacman）按优先级后续加
7. **“一键”形态**：A —— 预设模板 + 手动勾选（模板：AI 开发环境 / 前端开发 / 后端开发 / 极简起步）
8. **技术栈方案对比**：候选为 Wails v2（Go + WebView）、Tauri 2（Rust）、Electron；用户选择用 **Pake** 打包
9. **界面布局**：B —— 单页向导（选模板 → 确认清单 → 一键配置 → 实时日志）
10. **仓库名**：Envly（用户要求原创性，否决 DevKit 等常见名）
11. **语言**：中英双语（UI 与工具文案）
12. **网络区域**：大陆外不强制改镜像源；提供自动检测 / 手动覆盖，镜像改动可一键还原

## 关键决策摘要

| 维度 | 决策 |
| --- | --- |
| 架构 | Pake 壳 + Go 引擎（双进程）；备选 Wails 单文件 |
| UI 风格 | Apple Liquid Glass（Web 技术实现） |
| v1 平台 | Windows（winget / npm / pip / download / rustup） |
| v2 平台 | macOS（brew）→ Linux（apt 优先） |
| 语言 | 中英双语，引擎推送结构化消息 key，UI 本地化渲染 |
| 区域策略 | 自动检测连通性；大陆启用镜像，海外跳过镜像仅代理检测 |
| 测试 | Go 单测 + 前端组件测试 + Windows 集成冒烟 + 三平台 CI |
| 里程碑 | M1 引擎骨架 → M2 Windows 全流程 → M3 收尾打包 → M4 跨平台 |

## 设计文档

完整设计见 [docs/superpowers/specs/2026-08-16-envly-design.md](superpowers/specs/2026-08-16-envly-design.md)。

## 后续步骤

- 撰写实施计划（按 M1–M4 里程碑拆解任务）
- 桌面仓库已初始化并提交设计文档
