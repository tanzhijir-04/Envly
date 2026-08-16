package env

import (
	"reflect"
	"strings"
	"testing"
)

func TestMergePathDedupesAndPreservesOrder(t *testing.T) {
	got := MergePath([]string{"C:\\a", "C:\\b"}, []string{"C:\\b", "C:\\c"})
	want := []string{"C:\\a", "C:\\b", "C:\\c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNpmRegistrySet(t *testing.T) {
	if NpmRegistrySet("cn") != "https://registry.npmmirror.com" {
		t.Fatal("cn registry mismatch")
	}
	if NpmRegistrySet("global") != "https://registry.npmjs.org" {
		t.Fatal("global registry mismatch")
	}
}

func TestPipIndexURL(t *testing.T) {
	if !strings.Contains(PipIndexURL("cn"), "tuna") {
		t.Fatal("cn pip index mismatch")
	}
	if PipIndexURL("global") != "https://pypi.org/simple" {
		t.Fatal("global pip index mismatch")
	}
}

func TestNpmSetRegistryCmd(t *testing.T) {
	got := NpmSetRegistryCmd("https://x")
	want := []string{"npm", "config", "set", "registry", "https://x"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestUserPathScriptContainsAdd(t *testing.T) {
	script := UserPathScript("%APPDATA%\\npm")
	if !strings.Contains(script, "%APPDATA%\\npm") {
		t.Fatalf("script missing addition: %s", script)
	}
}

func TestProfileBlockHasMarkers(t *testing.T) {
	block := ProfileBlock()
	if !strings.Contains(block, "# >>> Envly >>>") || !strings.Contains(block, "oh-my-posh") {
		t.Fatalf("bad profile block: %s", block)
	}
}
