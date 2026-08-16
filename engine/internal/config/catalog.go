package config

// Catalog 是数据驱动的工具清单：Windows 为 M2 全量，macOS/Linux 在 M4 扩充。
type Group struct {
	ID     string `json:"id"`
	NameZh string `json:"name_zh"`
	NameEn string `json:"name_en"`
}

type InstallSpec struct {
	Method       string   `json:"method"`
	Package      string   `json:"package"`
	VerifyCmd    string   `json:"verify_cmd,omitempty"`
	DownloadURLs []string `json:"download_urls,omitempty"`
	SilentArgs   string   `json:"silent_args,omitempty"`
}

type Tool struct {
	ID       string                 `json:"id"`
	NameZh   string                 `json:"name_zh"`
	NameEn   string                 `json:"name_en"`
	DescZh   string                 `json:"desc_zh"`
	DescEn   string                 `json:"desc_en"`
	GroupID  string                 `json:"group_id"`
	Optional bool                   `json:"optional"`
	Install  map[string]InstallSpec `json:"install"`
}

type Template struct {
	ID      string   `json:"id"`
	NameZh  string   `json:"name_zh"`
	NameEn  string   `json:"name_en"`
	DescZh  string   `json:"desc_zh"`
	DescEn  string   `json:"desc_en"`
	ToolIDs []string `json:"tool_ids"`
}

var Groups = []Group{
	{ID: "runtime", NameZh: "基础运行时", NameEn: "Basic Runtime"},
	{ID: "ai", NameZh: "AI CLI", NameEn: "AI CLI"},
	{ID: "editor", NameZh: "编辑器 / 工具", NameEn: "Editors & Tools"},
	{ID: "terminal", NameZh: "终端增强", NameEn: "Terminal"},
	{ID: "env", NameZh: "系统环境", NameEn: "System"},
}

var Tools = []Tool{
	{ID: "nodejs", NameZh: "Node.js LTS", NameEn: "Node.js LTS", DescZh: "JavaScript 运行时", DescEn: "JavaScript runtime", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win":   {Method: "winget", Package: "OpenJS.NodeJS.LTS", VerifyCmd: "node -v"},
			"mac":   {Method: "brew", Package: "node", VerifyCmd: "node -v"},
			"linux": {Method: "nvm", Package: "node", VerifyCmd: "node -v"},
		}},
	{ID: "git", NameZh: "Git", NameEn: "Git", DescZh: "版本控制", DescEn: "Version control", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win": {Method: "download", Package: "Git", VerifyCmd: "git --version", DownloadURLs: []string{
				"https://github.com/git-for-windows/git/releases/latest/download/Git-2.49.0-64-bit.exe",
				"https://registry.npmmirror.com/-/binary/git-for-windows/2.49.0.windows.1/Git-2.49.0-64-bit.exe",
			}, SilentArgs: "/VERYSILENT /NORESTART /NOCANCEL /SP- /CLOSEAPPLICATIONS /COMPONENTS=icons,ext,ext\\shell\\here,gitlfs,assoc,assoc_sh"},
			"mac":   {Method: "brew", Package: "git", VerifyCmd: "git --version"},
			"linux": {Method: "apt", Package: "git", VerifyCmd: "git --version"},
		}},
	{ID: "python", NameZh: "Python 3", NameEn: "Python 3", DescZh: "Python 运行时", DescEn: "Python runtime", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win":   {Method: "winget", Package: "Python.Python.3.12", VerifyCmd: "python --version"},
			"mac":   {Method: "brew", Package: "python@3.12", VerifyCmd: "python3 --version"},
			"linux": {Method: "apt", Package: "python3", VerifyCmd: "python3 --version"},
		}},
	{ID: "typescript", NameZh: "TypeScript", NameEn: "TypeScript", DescZh: "TS 编译器", DescEn: "TS compiler", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win":   {Method: "npm", Package: "typescript", VerifyCmd: "tsc --version"},
			"mac":   {Method: "npm", Package: "typescript", VerifyCmd: "tsc --version"},
			"linux": {Method: "npm", Package: "typescript", VerifyCmd: "tsc --version"},
		}},
	{ID: "mingw", NameZh: "C/C++ · MinGW", NameEn: "C/C++ · MinGW", DescZh: "轻量 C/C++ 工具链", DescEn: "Lightweight C/C++ toolchain", GroupID: "runtime", Optional: true,
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "BrechtSanders.WinLibs.POSIX", VerifyCmd: "gcc --version"},
		}},
	{ID: "msvc", NameZh: "MSVC Build Tools", NameEn: "MSVC Build Tools", DescZh: "微软 C++ 工具链（可选）", DescEn: "Microsoft C++ toolchain (optional)", GroupID: "runtime", Optional: true,
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "Microsoft.VisualStudio.2022.BuildTools"},
		}},
	{ID: "php", NameZh: "PHP", NameEn: "PHP", DescZh: "PHP 运行时", DescEn: "PHP runtime", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "PHP.PHP", VerifyCmd: "php --version"},
		}},
	{ID: "uv", NameZh: "uv", NameEn: "uv", DescZh: "Python 包管理器", DescEn: "Python package manager", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "AstralSoftware.uv", VerifyCmd: "uv --version"},
		}},
	{ID: "jupyter", NameZh: "Jupyter Notebook", NameEn: "Jupyter Notebook", DescZh: "交互式笔记本", DescEn: "Interactive notebooks", GroupID: "runtime",
		Install: map[string]InstallSpec{
			"win": {Method: "pip", Package: "notebook", VerifyCmd: "jupyter --version"},
		}},
	{ID: "go", NameZh: "Go", NameEn: "Go", DescZh: "Go 语言", DescEn: "Go language", GroupID: "runtime", Optional: true,
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "GoLang.Go", VerifyCmd: "go version"},
		}},
	{ID: "rust", NameZh: "Rust", NameEn: "Rust", DescZh: "Rust 工具链（可选）", DescEn: "Rust toolchain (optional)", GroupID: "runtime", Optional: true,
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "Rustlang.Rustup", VerifyCmd: "rustc --version"},
		}},
	{ID: "java", NameZh: "Java Temurin", NameEn: "Java Temurin", DescZh: "JDK 21（可选）", DescEn: "JDK 21 (optional)", GroupID: "runtime", Optional: true,
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "EclipseAdoptium.Temurin.21.JDK", VerifyCmd: "java -version"},
		}},

	{ID: "claude-code", NameZh: "Claude Code", NameEn: "Claude Code", DescZh: "Anthropic AI 编程 CLI", DescEn: "Anthropic AI coding CLI", GroupID: "ai",
		Install: map[string]InstallSpec{
			"win":   {Method: "npm", Package: "@anthropic-ai/claude-code", VerifyCmd: "claude --version"},
			"mac":   {Method: "npm", Package: "@anthropic-ai/claude-code", VerifyCmd: "claude --version"},
			"linux": {Method: "npm", Package: "@anthropic-ai/claude-code", VerifyCmd: "claude --version"},
		}},
	{ID: "codex-cli", NameZh: "Codex CLI", NameEn: "Codex CLI", DescZh: "OpenAI 官方编程 CLI", DescEn: "OpenAI coding CLI", GroupID: "ai",
		Install: map[string]InstallSpec{
			"win":   {Method: "npm", Package: "@openai/codex", VerifyCmd: "codex --version"},
			"mac":   {Method: "npm", Package: "@openai/codex", VerifyCmd: "codex --version"},
			"linux": {Method: "npm", Package: "@openai/codex", VerifyCmd: "codex --version"},
		}},
	{ID: "gemini-cli", NameZh: "Gemini CLI", NameEn: "Gemini CLI", DescZh: "Google AI 编程 CLI", DescEn: "Google AI coding CLI", GroupID: "ai",
		Install: map[string]InstallSpec{
			"win": {Method: "npm", Package: "@google/gemini-cli", VerifyCmd: "gemini --version"},
		}},
	{ID: "pake", NameZh: "Pake", NameEn: "Pake", DescZh: "网页打包桌面应用（依赖 Rust）", DescEn: "Package web apps as desktop apps (needs Rust)", GroupID: "ai",
		Install: map[string]InstallSpec{
			"win": {Method: "npm", Package: "pake-cli", VerifyCmd: "pake --version"},
		}},

	{ID: "vscode", NameZh: "VS Code", NameEn: "VS Code", DescZh: "代码编辑器", DescEn: "Code editor", GroupID: "editor",
		Install: map[string]InstallSpec{
			"win":   {Method: "winget", Package: "Microsoft.VisualStudioCode", VerifyCmd: "code --version"},
			"mac":   {Method: "brew", Package: "visual-studio-code", VerifyCmd: "code --version"},
			"linux": {Method: "apt", Package: "code", VerifyCmd: "code --version"},
		}},
	{ID: "cursor", NameZh: "Cursor", NameEn: "Cursor", DescZh: "AI 原生代码编辑器", DescEn: "AI-native code editor", GroupID: "editor",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "Cursor.Cursor", VerifyCmd: "cursor --version"},
		}},
	{ID: "sublime", NameZh: "Sublime Text", NameEn: "Sublime Text", DescZh: "轻量编辑器", DescEn: "Lightweight editor", GroupID: "editor",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "SublimeHQ.SublimeText", VerifyCmd: "subl --version"},
		}},
	{ID: "drawio", NameZh: "draw.io", NameEn: "draw.io", DescZh: "流程图 / 图表工具", DescEn: "Diagrams and flowcharts", GroupID: "editor",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "JGraph.Draw", VerifyCmd: "draw.io --version"},
		}},
	{ID: "jetbrains", NameZh: "JetBrains Toolbox", NameEn: "JetBrains Toolbox", DescZh: "JetBrains IDE 管理（可选）", DescEn: "JetBrains IDE manager (optional)", GroupID: "editor", Optional: true,
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "JetBrains.Toolbox"},
		}},

	{ID: "windows-terminal", NameZh: "Windows Terminal", NameEn: "Windows Terminal", DescZh: "现代终端", DescEn: "Modern terminal", GroupID: "terminal",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "Microsoft.WindowsTerminal", VerifyCmd: "wt --version"},
		}},
	{ID: "powershell7", NameZh: "PowerShell 7", NameEn: "PowerShell 7", DescZh: "新一代 Shell", DescEn: "Next-gen shell", GroupID: "terminal",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "Microsoft.PowerShell", VerifyCmd: "pwsh --version"},
		}},
	{ID: "oh-my-posh", NameZh: "Oh My Posh", NameEn: "Oh My Posh", DescZh: "终端提示符美化", DescEn: "Terminal prompt theme", GroupID: "terminal",
		Install: map[string]InstallSpec{
			"win":   {Method: "winget", Package: "JanDeDobbeleer.OhMyPosh", VerifyCmd: "oh-my-posh --version"},
			"mac":   {Method: "brew", Package: "oh-my-posh", VerifyCmd: "oh-my-posh --version"},
			"linux": {Method: "apt", Package: "oh-my-posh", VerifyCmd: "oh-my-posh --version"},
		}},
	{ID: "cascadia-font", NameZh: "Cascadia Code NF", NameEn: "Cascadia Code NF", DescZh: "Nerd Font 字体", DescEn: "Nerd Font", GroupID: "terminal",
		Install: map[string]InstallSpec{
			"win": {Method: "winget", Package: "Microsoft.CascadiaCode"},
		}},

	{ID: "env-mirrors", NameZh: "npm 镜像源", NameEn: "npm Mirror", DescZh: "按网络区域自动切换", DescEn: "Region-aware npm registry", GroupID: "env",
		Install: map[string]InstallSpec{
			"win":   {Method: "config", Package: "mirrors", VerifyCmd: "npm config get registry"},
			"mac":   {Method: "config", Package: "mirrors", VerifyCmd: "npm config get registry"},
			"linux": {Method: "config", Package: "mirrors", VerifyCmd: "npm config get registry"},
		}},
	{ID: "env-pip", NameZh: "pip 镜像源", NameEn: "pip Mirror", DescZh: "按网络区域自动切换", DescEn: "Region-aware pip index", GroupID: "env",
		Install: map[string]InstallSpec{
			"win":   {Method: "config", Package: "mirrors", VerifyCmd: "pip config get global.index-url"},
			"mac":   {Method: "config", Package: "mirrors", VerifyCmd: "pip config get global.index-url"},
			"linux": {Method: "config", Package: "mirrors", VerifyCmd: "pip config get global.index-url"},
		}},
	{ID: "env-proxy", NameZh: "GitHub 加速 / 代理", NameEn: "GitHub Acceleration / Proxy", DescZh: "自动检测并配置", DescEn: "Auto-detect and configure", GroupID: "env",
		Install: map[string]InstallSpec{
			"win":   {Method: "config", Package: "proxy", VerifyCmd: "git config --get http.proxy"},
			"mac":   {Method: "config", Package: "proxy", VerifyCmd: "git config --get http.proxy"},
			"linux": {Method: "config", Package: "proxy", VerifyCmd: "git config --get http.proxy"},
		}},
	{ID: "env-shell", NameZh: "PowerShell Profile", NameEn: "PowerShell Profile", DescZh: "终端提示符配置", DescEn: "Terminal prompt config", GroupID: "env",
		Install: map[string]InstallSpec{
			"win": {Method: "config", Package: "profile"},
		}},
	{ID: "env-path", NameZh: "PATH 管理", NameEn: "PATH", DescZh: "合并用户 PATH", DescEn: "Merge user PATH", GroupID: "env",
		Install: map[string]InstallSpec{
			"win": {Method: "config", Package: "path"},
		}},
}

var Templates = []Template{
	{ID: "ai", NameZh: "AI 开发环境", NameEn: "AI Dev Environment", DescZh: "Node.js、Git、Claude Code、Codex、VS Code", DescEn: "Node.js, Git, Claude Code, Codex, VS Code",
		ToolIDs: []string{"nodejs", "git", "claude-code", "codex-cli", "vscode", "env-mirrors", "env-proxy"}},
	{ID: "frontend", NameZh: "前端开发", NameEn: "Frontend", DescZh: "Node.js、TypeScript、VS Code、Oh My Posh", DescEn: "Node.js, TypeScript, VS Code, Oh My Posh",
		ToolIDs: []string{"nodejs", "typescript", "git", "vscode", "oh-my-posh", "env-mirrors"}},
	{ID: "backend", NameZh: "后端开发", NameEn: "Backend", DescZh: "Python、Go、PHP、uv、Jupyter、draw.io", DescEn: "Python, Go, PHP, uv, Jupyter, draw.io",
		ToolIDs: []string{"python", "go", "php", "uv", "jupyter", "git", "drawio", "env-mirrors"}},
	{ID: "minimal", NameZh: "极简起步", NameEn: "Minimal", DescZh: "Git、Python、PowerShell 7、代理检测", DescEn: "Git, Python, PowerShell 7, proxy detection",
		ToolIDs: []string{"git", "python", "powershell7", "env-proxy"}},
}

func ToolByID(id string) (Tool, bool) {
	for _, t := range Tools {
		if t.ID == id {
			return t, true
		}
	}
	return Tool{}, false
}

func ToolsByGroup(groupID string) []Tool {
	var out []Tool
	for _, t := range Tools {
		if t.GroupID == groupID {
			out = append(out, t)
		}
	}
	return out
}

func TemplateToolIDs(t Template) []string {
	var out []string
	for _, id := range t.ToolIDs {
		if _, ok := ToolByID(id); ok {
			out = append(out, id)
		}
	}
	return out
}
