package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCanonicalSpecsDir_DefaultWhenConfigMissing(t *testing.T) {
	repo := t.TempDir()

	got, err := resolveCanonicalSpecsDir(repo, "")
	if err != nil {
		t.Fatalf("resolveCanonicalSpecsDir: %v", err)
	}
	if got != defaultCanonicalSpecsDir {
		t.Fatalf("expected %q, got %q", defaultCanonicalSpecsDir, got)
	}
}

func TestResolveCanonicalSpecsDir_FromConfig(t *testing.T) {
	repo := t.TempDir()
	config := `{
  "$schema": "` + aeConfigSchemaURL + `",
  "specRoot": "platform/specs"
}
`
	if err := os.WriteFile(filepath.Join(repo, ".ae-config.json"), []byte(config), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	got, err := resolveCanonicalSpecsDir(repo, "")
	if err != nil {
		t.Fatalf("resolveCanonicalSpecsDir: %v", err)
	}
	if got != "platform/specs" {
		t.Fatalf("expected %q, got %q", "platform/specs", got)
	}
}

func TestResolveCanonicalSpecsDir_FlagOverridesConfig(t *testing.T) {
	repo := t.TempDir()
	config := `{
  "specRoot": "platform/specs"
}
`
	if err := os.WriteFile(filepath.Join(repo, ".ae-config.json"), []byte(config), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	got, err := resolveCanonicalSpecsDir(repo, "custom/specs")
	if err != nil {
		t.Fatalf("resolveCanonicalSpecsDir: %v", err)
	}
	if got != "custom/specs" {
		t.Fatalf("expected %q, got %q", "custom/specs", got)
	}
}

func TestResolveCanonicalSpecsDir_RejectsEmptySpecRoot(t *testing.T) {
	repo := t.TempDir()
	config := `{
  "specRoot": ""
}
`
	if err := os.WriteFile(filepath.Join(repo, ".ae-config.json"), []byte(config), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := resolveCanonicalSpecsDir(repo, ""); err == nil {
		t.Fatal("expected error for empty specRoot")
	}
}
