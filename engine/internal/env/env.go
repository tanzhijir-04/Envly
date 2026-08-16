// Package env 提供环境配置的纯函数（PATH / Profile / 镜像 / 代理命令生成）。
package env

import "strings"

func MergePath(current, additions []string) []string {
	seen := make(map[string]bool, len(current)+len(additions))
	var out []string
	for _, p := range append(current, additions...) {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		out = append(out, trimmed)
	}
	return out
}

func NpmRegistrySet(region string) string {
	if region == "cn" {
		return "https://registry.npmmirror.com"
	}
	return "https://registry.npmjs.org"
}

func PipIndexURL(region string) string {
	if region == "cn" {
		return "https://pypi.tuna.tsinghua.edu.cn/simple"
	}
	return "https://pypi.org/simple"
}

func NpmSetRegistryCmd(registry string) []string {
	return []string{"npm", "config", "set", "registry", registry}
}

func NpmGetRegistryCmd() []string {
	return []string{"npm", "config", "get", "registry"}
}

func PipSetIndexCmd(url string) []string {
	return []string{"pip", "config", "set", "global.index-url", url}
}

func PipGetIndexCmd() []string {
	return []string{"pip", "config", "get", "global.index-url"}
}

func GitProxySetCmd(proxy string) []string {
	return []string{"git", "config", "--global", "http.proxy", proxy}
}

func GitProxyGetCmd() []string {
	return []string{"git", "config", "--get", "http.proxy"}
}

func UserPathScript(add string) string {
	return "$p=[Environment]::GetEnvironmentVariable('Path','User'); if($p -notlike '*" + add + "*'){[Environment]::SetEnvironmentVariable('Path',($p.TrimEnd(';')+';" + add + "'),'User')}"
}

func ProfileBlock() string {
	return "# >>> Envly >>>\n" +
		"if (Get-Command oh-my-posh -ErrorAction SilentlyContinue) {\n" +
		"  oh-my-posh init pwsh | Invoke-Expression\n" +
		"}\n" +
		"# <<< Envly <<<"
}
