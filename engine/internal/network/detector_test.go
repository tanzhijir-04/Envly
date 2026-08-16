package network

import (
	"context"
	"errors"
	"testing"
)

func TestDetectManualOverride(t *testing.T) {
	d := NewDetector(func(context.Context, string) error { return nil })
	if got := d.Detect(context.Background(), "cn"); got.Region != "cn" {
		t.Fatalf("manual cn: %+v", got)
	}
	if got := d.Detect(context.Background(), "global"); got.Region != "global" {
		t.Fatalf("manual global: %+v", got)
	}
}

func TestDetectAutoCNWhenOfficialDownMirrorUp(t *testing.T) {
	d := NewDetector(func(_ context.Context, url string) error {
		if url == "https://registry.npmjs.org" {
			return errors.New("down")
		}
		return nil
	})
	got := d.Detect(context.Background(), "auto")
	if got.Region != "cn" || !got.MirrorOK || got.OfficialOK {
		t.Fatalf("expected cn, got %+v", got)
	}
}

func TestDetectAutoGlobalWhenOfficialUp(t *testing.T) {
	d := NewDetector(func(context.Context, string) error { return nil })
	got := d.Detect(context.Background(), "auto")
	if got.Region != "global" || !got.OfficialOK {
		t.Fatalf("expected global, got %+v", got)
	}
}
