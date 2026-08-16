// Package installer 按安装方式分发真实安装。
package installer

import (
	"context"
	"fmt"

	"github.com/tanzhijir-04/Envly/engine/internal/config"
	"github.com/tanzhijir-04/Envly/engine/internal/runner"
)

type Installer struct {
	run runner.Runner
}

func New(run runner.Runner) *Installer {
	return &Installer{run: run}
}

func (i *Installer) Install(ctx context.Context, spec config.InstallSpec) error {
	switch spec.Method {
	case "winget":
		_, err := i.run.Run(ctx, "winget", "install", "--id", spec.Package, "--accept-package-agreements", "--accept-source-agreements", "--silent")
		return err
	case "npm":
		_, err := i.run.Run(ctx, "npm", "install", "-g", spec.Package)
		return err
	case "pip":
		_, err := i.run.Run(ctx, "pip", "install", spec.Package)
		return err
	case "rustup":
		_, err := i.run.Run(ctx, "rustup", "default", "stable")
		return err
	case "download":
		return i.installDownload(ctx, spec)
	default:
		return fmt.Errorf("unsupported method %q", spec.Method)
	}
}

func (i *Installer) installDownload(_ context.Context, _ config.InstallSpec) error {
	return fmt.Errorf("download not implemented yet")
}
