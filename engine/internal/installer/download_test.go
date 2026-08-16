package installer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadToWritesBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("payload"))
	}))
	defer ts.Close()
	dst := filepath.Join(t.TempDir(), "installer.exe")
	if err := downloadTo(context.Background(), ts.URL, dst); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "payload" {
		t.Fatalf("got %q", string(b))
	}
}

func TestDownloadToFailsOnNonOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer ts.Close()
	if err := downloadTo(context.Background(), ts.URL, filepath.Join(t.TempDir(), "x.exe")); err == nil {
		t.Fatal("expected error on 404")
	}
}

func TestParseArgsSplitsQuotedStrings(t *testing.T) {
	got := parseArgs(`/VERYSILENT /ARG "a b" /X`)
	want := []string{"/VERYSILENT", "/ARG", "a b", "/X"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
