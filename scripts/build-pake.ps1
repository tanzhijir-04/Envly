param(
  [string]$UiUrl = "http://127.0.0.1:17521",
  [string]$Name = "Envly"
)
$ErrorActionPreference = "Stop"

# 1. Rust 工具链检查
if (-not (Get-Command rustc -ErrorAction SilentlyContinue)) {
  throw "Rust toolchain required. Install: winget install Rustlang.Rustup; rustup default stable"
}

# 2. pake-cli 检查
if (-not (Get-Command pake -ErrorAction SilentlyContinue)) {
  npm install -g pake-cli
}

# 3. 引擎必须已在 $UiUrl 运行
try {
  Invoke-WebRequest "$UiUrl/api/health" -UseBasicParsing -TimeoutSec 3 | Out-Null
} catch {
  throw "Engine not running at $UiUrl — start it first: go run ./cmd/envly"
}

# 4. 打包（Windows 产出 msi；--keep-binary 同时保留 exe）
$icon = Join-Path $PSScriptRoot "..\assets\icon.png"
pake $UiUrl --name $Name --icon $icon --width 1200 --height 780 --installer-language zh-CN --keep-binary --hide-window-decorations
