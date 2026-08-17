<p align="right"><b>English</b> · <b><a href="README.md">简体中文</a></b></p>

<p align="center">
  <img src="assets/Envly_Icon_Corrected.svg" width="120" alt="Envly">
</p>

<h1 align="center">Envly</h1>

<p align="center">
  <strong>One-click developer environment setup tool</strong><br>
  Windows · macOS · Linux
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white">
  <img alt="Pake" src="https://img.shields.io/badge/Pake-Tauri-24C8DB?logo=rust&logoColor=white">
  <img alt="License" src="https://img.shields.io/badge/License-MIT-brightgreen">
</p>

Envly sets up your development environment and system configuration in one click: pick a template, tweak the checklist, run, and watch the live log and report — all inside a single-page wizard styled with Apple Liquid Glass. Mirror sources are skipped automatically outside mainland China, and the UI is bilingual.

## Download

Grab the installer for your platform from [Releases](https://github.com/tanzhijir-04/Envly/releases):

- **Windows**: `Envly_1.0.0_x64_zh-CN.msi` (installer) or `Envly.exe` (portable)
- **macOS**: `Envly.dmg`
- **Linux**: `envly.deb` (Debian/Ubuntu) or `envly.AppImage`

> The macOS / Linux builds use a "Pake shell + standalone engine" two-process design. Start the matching `envly-engine-*` binary before launching the app (a one-click launcher is planned).

## Features

- **Single-page wizard**: template → checklist → one-click run → live log
- **Real installs**: winget / npm / pip / rustup / download; installed items are skipped
- **Environment config**: npm/pip mirrors, GitHub acceleration/proxy, PowerShell profile, PATH — auditable and restorable
- **Region-aware**: auto-detect network region; mirrors are skipped overseas
- **Execution report**: install records and environment changes at a glance
- **Desktop experience**: frameless window, custom window controls, rounded corners, Liquid Glass motion
- **Bilingual**: instant zh-CN / en switching

## Development quick start

Prerequisites: Go 1.23+, Node 18+.

```bash
cd engine
go run ./cmd/envly -web-dir ../ui
```

Then open http://127.0.0.1:17521

Simulation mode (no real installs):

```bash
go run ./cmd/envly -web-dir ../ui -simulate
```

## Tech stack

| Layer | Technology |
| --- | --- |
| Engine | Go 1.23 (stdlib only: net/http, SSE, Win32 syscall) |
| Frontend | Vanilla HTML/CSS/JS (no framework), Apple Liquid Glass design |
| Motion | GSAP (non-linear easing per Liquid Glass motion spec) |
| Desktop shell | Pake (Tauri / WebView2) |
| Build & release | GitHub Actions (native Windows / macOS / Linux runners) |
| Tests | Go test · Vitest |

## Directory structure

```
Envly/
├─ engine/                    # Go engine: local HTTP service + business logic
│  ├─ cmd/envly/              # Entry point (server, window-control daemon)
│  └─ internal/
│     ├─ api/                 # REST API + SSE live logs + static hosting
│     ├─ config/              # Data-driven tool catalog and templates
│     ├─ env/                 # Mirrors/proxy/profile/PATH config & restore
│     ├─ events/              # SSE event hub
│     ├─ executor/            # Executor (simulated / real installs)
│     ├─ installer/           # Install dispatch (winget/npm/pip/rustup/download)
│     ├─ network/             # Network region detection
│     ├─ runner/              # Command runner abstraction
│     ├─ state/               # Settings persistence
│     ├─ store/               # Install records and env-op logs
│     ├─ verify/              # Post-install verification & version parsing
│     └─ windowctrl/          # Windows window control (frameless/corners/min-max)
├─ ui/                        # Frontend single-page app
│  ├─ src/                    # main/i18n/api/plan/status/motion modules
│  ├─ vendor/gsap.min.js      # GSAP (vendored, offline-capable)
│  └─ index.html · styles.css
├─ scripts/                   # Icon generation, Pake packaging scripts
├─ assets/                    # Icon assets
├─ docs/                      # Design docs / plans / notes
├─ .github/workflows/         # CI and cross-platform release workflow
└─ README.md · README.en.md
```

## Tests

```bash
cd engine && go test ./...
cd ui && npm test
```

## Packaging & releases

Local Windows packaging:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/gen-icon.ps1
# start the engine first, then:
powershell -ExecutionPolicy Bypass -File scripts/build-pake.ps1
```

Cross-platform packaging: push a `v*` tag and the [release workflow](.github/workflows/release.yml) builds and attaches Windows / macOS / Linux artifacts to the matching release.

## Roadmap

- M1 ✅ Engine skeleton + single-page wizard
- M2 ✅ Windows real installs and environment config
- M3 ✅ Report / restore / packaging / CI / Release
- M4 ⏳ Deeper macOS (brew) and Linux (apt) support, one-click launcher

## License

MIT
