package verify

import (
	"context"
	"testing"
)

type fakeRunner struct {
	out string
	err error
}

func (f fakeRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	return f.out, f.err
}

func TestParseVersion(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"v22.14.0", "22.14.0"},
		{"22.14.0\nmore", "22.14.0"},
		{"version 3.12.8", "3.12.8"},
		{"  1.2.3  ", "1.2.3"},
	}
	for _, c := range cases {
		if got := ParseVersion(c.in); got != c.want {
			t.Fatalf("ParseVersion(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCheckReturnsInstalledAndVersion(t *testing.T) {
	v := New(fakeRunner{out: "v22.14.0"})
	version, ok := v.Check(context.Background(), "node -v")
	if !ok {
		t.Fatal("expected installed")
	}
	if version != "22.14.0" {
		t.Fatalf("got version %q", version)
	}
}

func TestCheckFalseOnErrorOrEmptyCmd(t *testing.T) {
	v := New(fakeRunner{err: context.DeadlineExceeded})
	if _, ok := v.Check(context.Background(), "node -v"); ok {
		t.Fatal("expected not installed on error")
	}
	if _, ok := v.Check(context.Background(), ""); ok {
		t.Fatal("expected not installed on empty command")
	}
}
