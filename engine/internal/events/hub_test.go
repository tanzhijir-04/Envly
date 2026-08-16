package events

import (
	"testing"
	"time"
)

func TestPublishBroadcastsToSubscribers(t *testing.T) {
	hub := NewHub()
	ch, unsubscribe := hub.Subscribe()
	defer unsubscribe()
	hub.Publish(Event{Type: "running", RunID: "r1", ToolID: "nodejs"})
	select {
	case e := <-ch:
		if e.RunID != "r1" || e.ToolID != "nodejs" {
			t.Fatalf("unexpected event: %+v", e)
		}
		if e.TS == 0 {
			t.Fatal("expected TS to be set")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestSubscribeReplaysHistory(t *testing.T) {
	hub := NewHub()
	hub.Publish(Event{Type: "running", RunID: "r1", ToolID: "git"})
	hub.Publish(Event{Type: "success", RunID: "r1", ToolID: "git"})
	ch, unsubscribe := hub.Subscribe()
	defer unsubscribe()
	for i := 0; i < 2; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for replayed history")
		}
	}
}
