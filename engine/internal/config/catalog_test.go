package config

import "testing"

func TestToolByID(t *testing.T) {
	if _, ok := ToolByID("nodejs"); !ok {
		t.Fatal("expected nodejs to exist")
	}
	if _, ok := ToolByID("does-not-exist"); ok {
		t.Fatal("expected unknown tool to be missing")
	}
}

func TestToolsByGroupOnlyReturnsMatchingGroup(t *testing.T) {
	tools := ToolsByGroup("ai")
	for _, tool := range tools {
		if tool.GroupID != "ai" {
			t.Fatalf("tool %s has group %s, want ai", tool.ID, tool.GroupID)
		}
	}
	if len(tools) == 0 {
		t.Fatal("expected at least one AI tool")
	}
}

func TestTemplateToolIDsDropsUnknown(t *testing.T) {
	tmpl := Template{ID: "t", ToolIDs: []string{"nodejs", "ghost"}}
	ids := TemplateToolIDs(tmpl)
	if len(ids) != 1 || ids[0] != "nodejs" {
		t.Fatalf("expected only nodejs, got %v", ids)
	}
}

func TestCatalogHasFullWindowsSet(t *testing.T) {
	required := []string{
		"nodejs", "git", "python", "typescript", "mingw", "msvc", "php", "uv",
		"jupyter", "go", "rust", "java", "claude-code", "codex-cli", "gemini-cli",
		"pake", "vscode", "cursor", "sublime", "drawio", "jetbrains",
		"windows-terminal", "powershell7", "oh-my-posh", "cascadia-font",
		"env-mirrors", "env-pip", "env-proxy", "env-shell", "env-path",
	}
	for _, id := range required {
		tool, ok := ToolByID(id)
		if !ok {
			t.Fatalf("missing tool %s", id)
		}
		if _, ok := tool.Install["win"]; !ok {
			t.Fatalf("tool %s missing win install spec", id)
		}
	}
}

func TestTemplatesAllReferenceExistingTools(t *testing.T) {
	for _, tmpl := range Templates {
		if len(TemplateToolIDs(tmpl)) != len(tmpl.ToolIDs) {
			t.Fatalf("template %s references unknown tool", tmpl.ID)
		}
	}
}
