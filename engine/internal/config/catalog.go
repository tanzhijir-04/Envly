package config

// Catalog 是 M1 骨架数据集：结构完整，工具列表在 M2 扩充到设计文档全量。
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
			"win":   {Method: "winget", Package: "Git.Git", VerifyCmd: "git --version"},
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
	{ID: "vscode", NameZh: "VS Code", NameEn: "VS Code", DescZh: "代码编辑器", DescEn: "Code editor", GroupID: "editor",
		Install: map[string]InstallSpec{
			"win":   {Method: "winget", Package: "Microsoft.VisualStudioCode", VerifyCmd: "code --version"},
			"mac":   {Method: "brew", Package: "visual-studio-code", VerifyCmd: "code --version"},
			"linux": {Method: "apt", Package: "code", VerifyCmd: "code --version"},
		}},
	{ID: "oh-my-posh", NameZh: "Oh My Posh", NameEn: "Oh My Posh", DescZh: "终端提示符美化", DescEn: "Terminal prompt theme", GroupID: "terminal",
		Install: map[string]InstallSpec{
			"win":   {Method: "winget", Package: "JanDeDobbeleer.OhMyPosh", VerifyCmd: "oh-my-posh --version"},
			"mac":   {Method: "brew", Package: "oh-my-posh", VerifyCmd: "oh-my-posh --version"},
			"linux": {Method: "apt", Package: "oh-my-posh", VerifyCmd: "oh-my-posh --version"},
		}},
	{ID: "env-mirrors", NameZh: "npm / pip 镜像源", NameEn: "npm / pip Mirrors", DescZh: "按网络区域自动切换", DescEn: "Region-aware mirror switching", GroupID: "env",
		Install: map[string]InstallSpec{
			"win":   {Method: "config", Package: "mirrors", VerifyCmd: "npm config get registry"},
			"mac":   {Method: "config", Package: "mirrors", VerifyCmd: "npm config get registry"},
			"linux": {Method: "config", Package: "mirrors", VerifyCmd: "npm config get registry"},
		}},
	{ID: "env-proxy", NameZh: "GitHub 加速 / 代理", NameEn: "GitHub Acceleration / Proxy", DescZh: "自动检测并配置", DescEn: "Auto-detect and configure", GroupID: "env",
		Install: map[string]InstallSpec{
			"win":   {Method: "config", Package: "proxy", VerifyCmd: "git config --get http.proxy"},
			"mac":   {Method: "config", Package: "proxy", VerifyCmd: "git config --get http.proxy"},
			"linux": {Method: "config", Package: "proxy", VerifyCmd: "git config --get http.proxy"},
		}},
}

var Templates = []Template{
	{ID: "ai", NameZh: "AI 开发环境", NameEn: "AI Dev Environment", DescZh: "Node.js、Git、Claude Code、Codex、VS Code", DescEn: "Node.js, Git, Claude Code, Codex, VS Code",
		ToolIDs: []string{"nodejs", "git", "claude-code", "codex-cli", "vscode", "env-mirrors", "env-proxy"}},
	{ID: "frontend", NameZh: "前端开发", NameEn: "Frontend", DescZh: "Node.js、TypeScript、VS Code、Oh My Posh", DescEn: "Node.js, TypeScript, VS Code, Oh My Posh",
		ToolIDs: []string{"nodejs", "typescript", "git", "vscode", "oh-my-posh", "env-mirrors"}},
	{ID: "backend", NameZh: "后端开发", NameEn: "Backend", DescZh: "Python、TypeScript、Git、镜像源", DescEn: "Python, TypeScript, Git, mirrors",
		ToolIDs: []string{"python", "typescript", "git", "env-mirrors"}},
	{ID: "minimal", NameZh: "极简起步", NameEn: "Minimal", DescZh: "Git、Python、代理检测", DescEn: "Git, Python, proxy detection",
		ToolIDs: []string{"git", "python", "env-proxy"}},
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
