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
