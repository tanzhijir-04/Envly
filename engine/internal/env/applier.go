package env

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/tanzhijir-04/Envly/engine/internal/runner"
	"github.com/tanzhijir-04/Envly/engine/internal/store"
)

// Applier 通过 runner 执行环境配置命令，并把变更记录到 store 以便回滚。
type Applier struct {
	run   runner.Runner
	store *store.DB
}

func NewApplier(run runner.Runner, db *store.DB) *Applier {
	return &Applier{run: run, store: db}
}

func (a *Applier) ApplyMirrors(ctx context.Context, region string) error {
	before := a.npmRegistry(ctx)
	registry := NpmRegistrySet(region)
	if _, err := a.run.Run(ctx, NpmSetRegistryCmd(registry)[0], NpmSetRegistryCmd(registry)[1:]...); err != nil {
		return err
	}
	after := a.npmRegistry(ctx)
	return a.store.AppendEnvOp(store.EnvOp{Key: "npm_registry", Before: before, After: after})
}

func (a *Applier) ApplyPipMirror(ctx context.Context, region string) error {
	before := a.pipIndex(ctx)
	url := PipIndexURL(region)
	if _, err := a.run.Run(ctx, PipSetIndexCmd(url)[0], PipSetIndexCmd(url)[1:]...); err != nil {
		return err
	}
	after := a.pipIndex(ctx)
	return a.store.AppendEnvOp(store.EnvOp{Key: "pip_index", Before: before, After: after})
}

func (a *Applier) ApplyProxy(ctx context.Context) error {
	proxy := strings.TrimSpace(os.Getenv("HTTP_PROXY"))
	if proxy == "" {
		proxy = strings.TrimSpace(os.Getenv("http_proxy"))
	}
	if proxy == "" {
		return nil
	}
	before := a.gitProxy(ctx)
	if _, err := a.run.Run(ctx, GitProxySetCmd(proxy)[0], GitProxySetCmd(proxy)[1:]...); err != nil {
		return err
	}
	return a.store.AppendEnvOp(store.EnvOp{Key: "git_http_proxy", Before: before, After: proxy})
}

func (a *Applier) ApplyPath(ctx context.Context, add string) error {
	_, err := a.run.Run(ctx, "powershell", "-NoProfile", "-Command", UserPathScript(add))
	return err
}

func (a *Applier) ApplyProfile(ctx context.Context) error {
	profile, err := a.profilePath(ctx)
	if err != nil {
		return err
	}
	b, readErr := os.ReadFile(profile)
	exists := readErr == nil
	content := string(b)
	if strings.Contains(content, "# >>> Envly >>>") {
		return nil
	}
	if exists {
		content += "\n"
	}
	content += ProfileBlock() + "\n"
	if err := os.MkdirAll(filepath.Dir(profile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(profile, []byte(content), 0o644)
}

func (a *Applier) profilePath(ctx context.Context) (string, error) {
	out, err := a.run.Run(ctx, "powershell", "-NoProfile", "-Command", "$PROFILE")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (a *Applier) npmRegistry(ctx context.Context) string {
	cmd := NpmGetRegistryCmd()
	out, _ := a.run.Run(ctx, cmd[0], cmd[1:]...)
	return strings.TrimSpace(out)
}

func (a *Applier) pipIndex(ctx context.Context) string {
	cmd := PipGetIndexCmd()
	out, _ := a.run.Run(ctx, cmd[0], cmd[1:]...)
	return strings.TrimSpace(out)
}

func (a *Applier) gitProxy(ctx context.Context) string {
	cmd := GitProxyGetCmd()
	out, _ := a.run.Run(ctx, cmd[0], cmd[1:]...)
	return strings.TrimSpace(out)
}
