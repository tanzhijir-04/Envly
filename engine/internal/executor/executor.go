package executor

import (
	"context"
	"time"
)

type Progress struct {
	ToolID string
	Status string // running | success | failed
	MsgKey string
	Params map[string]any
}

type Executor interface {
	Run(ctx context.Context, toolIDs []string, emit func(Progress)) error
}

type Simulated struct {
	Delay time.Duration
}

func (s Simulated) Run(ctx context.Context, toolIDs []string, emit func(Progress)) error {
	for _, id := range toolIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		emit(Progress{ToolID: id, Status: "running", MsgKey: "tool.start", Params: map[string]any{"tool": id}})
		if !sleep(ctx, s.Delay) {
			return ctx.Err()
		}
		emit(Progress{ToolID: id, Status: "running", MsgKey: "tool.progress", Params: map[string]any{"tool": id, "percent": 50}})
		if !sleep(ctx, s.Delay) {
			return ctx.Err()
		}
		emit(Progress{ToolID: id, Status: "success", MsgKey: "tool.done", Params: map[string]any{"tool": id}})
	}
	return nil
}

func sleep(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		d = 5 * time.Millisecond
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
