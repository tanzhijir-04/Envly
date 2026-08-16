package executor

import (
	"context"
	"fmt"

	"github.com/tanzhijir-04/Envly/engine/internal/config"
	"github.com/tanzhijir-04/Envly/engine/internal/store"
)

// 小接口便于测试注入。
type InstallRunner interface {
	Install(ctx context.Context, spec config.InstallSpec) error
}

type EnvApplier interface {
	ApplyMirrors(ctx context.Context, region string) error
	ApplyPipMirror(ctx context.Context, region string) error
	ApplyProxy(ctx context.Context) error
	ApplyProfile(ctx context.Context) error
	ApplyPath(ctx context.Context, add string) error
}

type Verifier interface {
	Check(ctx context.Context, cmdline string) (string, bool)
}

type Recorder interface {
	AppendRecord(r store.Record) error
}

type Real struct {
	inst     InstallRunner
	env      EnvApplier
	ver      Verifier
	records  Recorder
	region   func() string
	platform func() string
}

func NewReal(inst InstallRunner, env EnvApplier, ver Verifier, records Recorder, region func() string, platform func() string) *Real {
	return &Real{inst: inst, env: env, ver: ver, records: records, region: region, platform: platform}
}

func (r *Real) Run(ctx context.Context, toolIDs []string, emit func(Progress)) error {
	for _, id := range toolIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		tool, ok := config.ToolByID(id)
		if !ok {
			emit(Progress{ToolID: id, Status: "failed", MsgKey: "tool.unknown", Params: map[string]any{"tool": id}})
			continue
		}
		spec, ok := tool.Install[r.platform()]
		if !ok {
			emit(Progress{ToolID: id, Status: "skipped", MsgKey: "tool.unsupported", Params: map[string]any{"tool": tool.NameEn}})
			continue
		}
		if version, installed := r.ver.Check(ctx, spec.VerifyCmd); installed {
			emit(Progress{ToolID: id, Status: "skipped", MsgKey: "tool.installed", Params: map[string]any{"tool": tool.NameEn, "version": version}})
			continue
		}
		if spec.Method == "config" {
			r.runConfig(ctx, id, spec, emit)
			continue
		}
		emit(Progress{ToolID: id, Status: "running", MsgKey: "tool.start", Params: map[string]any{"tool": tool.NameEn}})
		if err := r.inst.Install(ctx, spec); err != nil {
			emit(Progress{ToolID: id, Status: "failed", MsgKey: "tool.failed", Params: map[string]any{"tool": tool.NameEn, "error": err.Error()}})
			continue
		}
		version, ok := r.ver.Check(ctx, spec.VerifyCmd)
		if !ok {
			version = "installed"
		}
		_ = r.records.AppendRecord(store.Record{ToolID: tool.ID, Name: tool.NameEn, Version: version, Method: spec.Method})
		emit(Progress{ToolID: id, Status: "success", MsgKey: "tool.done", Params: map[string]any{"tool": tool.NameEn, "version": version}})
	}
	return nil
}

func (r *Real) runConfig(ctx context.Context, id string, spec config.InstallSpec, emit func(Progress)) {
	region := r.region()
	if (id == "env-mirrors" || id == "env-pip") && region != "cn" {
		emit(Progress{ToolID: id, Status: "skipped", MsgKey: "tool.mirror.skipped", Params: map[string]any{"tool": id, "region": region}})
		return
	}
	var err error
	switch id {
	case "env-mirrors":
		err = r.env.ApplyMirrors(ctx, region)
	case "env-pip":
		err = r.env.ApplyPipMirror(ctx, region)
	case "env-proxy":
		err = r.env.ApplyProxy(ctx)
	case "env-shell":
		err = r.env.ApplyProfile(ctx)
	case "env-path":
		err = r.env.ApplyPath(ctx, `%APPDATA%\npm`)
	default:
		err = fmt.Errorf("unknown config tool %q", id)
	}
	if err != nil {
		emit(Progress{ToolID: id, Status: "failed", MsgKey: "tool.failed", Params: map[string]any{"tool": id, "error": err.Error()}})
		return
	}
	emit(Progress{ToolID: id, Status: "success", MsgKey: "tool.done", Params: map[string]any{"tool": id}})
}
