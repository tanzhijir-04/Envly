package runner

import (
	"context"
	"testing"
)

func TestOSRunsCommandAndReturnsOutput(t *testing.T) {
	out, err := OS{}.Run(context.Background(), "go", "version")
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Fatal("expected non-empty output from go version")
	}
}
