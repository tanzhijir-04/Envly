// Package runner 抽象系统命令执行，便于测试注入 fake。
package runner

import (
	"context"
	"os/exec"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
}

// OS 是真实实现：直接执行系统命令。
type OS struct{}

func (OS) Run(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
