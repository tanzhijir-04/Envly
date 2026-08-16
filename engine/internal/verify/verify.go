// Package verify 负责执行验证命令并解析版本。
package verify

import (
	"context"
	"strings"

	"github.com/tanzhijir-04/Envly/engine/internal/runner"
)

type Verifier struct {
	run runner.Runner
}

func New(run runner.Runner) *Verifier {
	return &Verifier{run: run}
}

// Check 执行 cmdline（如 "node -v"），返回版本与是否成功。
func (v *Verifier) Check(ctx context.Context, cmdline string) (string, bool) {
	parts := strings.Fields(cmdline)
	if len(parts) == 0 {
		return "", false
	}
	out, err := v.run.Run(ctx, parts[0], parts[1:]...)
	if err != nil {
		return "", false
	}
	return ParseVersion(out), true
}

// ParseVersion 提取首行版本号，去掉常见前缀。
func ParseVersion(output string) string {
	line := strings.TrimSpace(output)
	if i := strings.IndexByte(line, '\n'); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	for _, prefix := range []string{"git version ", "version ", "Version ", "v"} {
		line = strings.TrimPrefix(line, prefix)
	}
	return line
}
