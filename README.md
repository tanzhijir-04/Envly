# Envly

一键开发环境配置桌面工具（Windows → macOS → Linux）

One-click developer environment setup tool (Windows → macOS → Linux)

Envly 帮你用一次点击完成开发环境安装与系统环境配置：选模板、勾选清单、一键配置、实时日志、执行报告。海外网络自动跳过镜像源，中英双语界面。

Envly sets up your dev environment in one click: pick a template, tweak the checklist, run, and read the live log and report. Mirror sources are skipped automatically outside mainland China. UI is bilingual (中文 / English).

## 功能 / Features

- 单页向导：模板 → 勾选 → 一键配置 → 实时日志 / Single-page wizard: template → checklist → one-click run → live log
- 真实安装：winget / npm / pip / rustup / download / Real installs via winget, npm, pip, rustup, download
- 环境配置：npm/pip 镜像、GitHub 加速/代理、PowerShell Profile、PATH / Environment: npm/pip mirrors, proxy, PowerShell profile, PATH
- 区域自适应：海外不强制改镜像 / Region-aware: mirrors skipped overseas
- 执行报告与一键还原 / Execution report and one-click restore
- 中英双语 / Bilingual (zh-CN / en)

## 快速开始 / Quick start

```bash
cd engine
go run ./cmd/envly -web-dir ../ui
```

打开 http://127.0.0.1:17521

开发模式（模拟执行，不装软件）：

```bash
go run ./cmd/envly -web-dir ../ui -simulate
```

## 测试 / Tests

```bash
cd engine && go test ./...
cd ui && npm test
```

## 打包 / Packaging

```powershell
powershell -ExecutionPolicy Bypass -File scripts/gen-icon.ps1
# 先启动引擎，再执行：
powershell -ExecutionPolicy Bypass -File scripts/build-pake.ps1
```

构建机需要 Node 18+、Rust 工具链、`pake-cli`。CI 见 `.github/workflows/build.yml`。

## 架构 / Architecture

Pake 壳 + Go 引擎（双进程）：引擎提供本地 HTTP 服务（127.0.0.1:17521）与 /api，前端为 Apple Liquid Glass 风格的单页应用。详见 [docs/superpowers/specs/2026-08-16-envly-design.md](docs/superpowers/specs/2026-08-16-envly-design.md)。

## 路线图 / Roadmap

- M1 ✅ 引擎骨架 + 单页向导
- M2 ✅ Windows 真实安装与环境配置
- M3 🔨 报告 / 还原 / 打包 / CI
- M4 ⏳ macOS（brew）→ Linux（apt）

## License

MIT
