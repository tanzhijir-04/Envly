package executor

import (
	"context"
	"testing"
	"time"
)

func TestSimulatedEmitsStartAndSuccessForEachTool(t *testing.T) {
	exec := Simulated{Delay: time.Millisecond}
	var got []Progress
	err := exec.Run(context.Background(), []string{"nodejs", "git"}, func(p Progress) {
		got = append(got, p)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 6 { // 2 tools x (running + running + success)
		t.Fatalf("expected 6 progress events, got %d", len(got))
	}
	if got[0].ToolID != "nodejs" || got[0].Status != "running" {
		t.Fatalf("unexpected first event: %+v", got[0])
	}
	if got[len(got)-1].Status != "success" {
		t.Fatalf("expected final success, got %+v", got[len(got)-1])
	}
}

func TestSimulatedStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	exec := Simulated{Delay: time.Hour}
	err := exec.Run(ctx, []string{"nodejs"}, func(Progress) {})
	if err == nil {
		t.Fatal("expected error after cancel")
	}
}
