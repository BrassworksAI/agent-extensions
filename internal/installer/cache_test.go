package installer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInstall_WritesCacheMetadata(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)
	inst.SetSource("v1.2.3")

	_, err := inst.Install("dir-based", ScopeLocal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	cacheDir := filepath.Join(projectRoot, ".agents", "ae")
	metadataPath := filepath.Join(cacheDir, CacheMetadataFile)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatalf("failed reading metadata: %v", err)
	}

	var metadata CacheMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		t.Fatalf("failed unmarshaling metadata: %v", err)
	}

	if metadata.Source != "v1.2.3" {
		t.Fatalf("metadata source = %q, want %q", metadata.Source, "v1.2.3")
	}

	if len(metadata.Files) == 0 {
		t.Fatal("metadata files should not be empty")
	}

	foundCommand := false
	for _, file := range metadata.Files {
		if file.Path == "commands/cmd-one.md" {
			foundCommand = true
			break
		}
	}
	if !foundCommand {
		t.Fatal("expected commands/cmd-one.md in metadata files")
	}
}

func TestCheckCache_DetectsTamperedFile(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)
	inst.SetSource("v1.2.3")

	_, err := inst.Install("dir-based", ScopeLocal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	cacheDir := filepath.Join(projectRoot, ".agents", "ae")
	tamperedPath := filepath.Join(cacheDir, "commands", "cmd-one.md")
	if err := os.WriteFile(tamperedPath, []byte("tampered"), 0644); err != nil {
		t.Fatalf("failed to tamper file: %v", err)
	}

	result, err := inst.CheckCache(ScopeLocal, "v1.2.3")
	if err != nil {
		t.Fatalf("CheckCache failed: %v", err)
	}

	if len(result.Issues) == 0 {
		t.Fatal("expected cache issues for tampered file")
	}

	foundHashMismatch := false
	for _, issue := range result.Issues {
		if issue.Kind == "hash mismatch" && issue.Path == "commands/cmd-one.md" {
			foundHashMismatch = true
			break
		}
	}
	if !foundHashMismatch {
		t.Fatal("expected hash mismatch issue for commands/cmd-one.md")
	}
}

func TestCheckCache_DetectsSourceMismatch(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)
	inst.SetSource("v1.0.0")

	_, err := inst.Install("dir-based", ScopeLocal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	result, err := inst.CheckCache(ScopeLocal, "v2.0.0")
	if err != nil {
		t.Fatalf("CheckCache failed: %v", err)
	}

	if !result.SourceMismatch {
		t.Fatal("expected source mismatch")
	}
}

func TestRepairCache_RebuildsAndUpdatesMetadata(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)
	inst.SetSource("v1.0.0")

	_, err := inst.Install("dir-based", ScopeLocal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	cacheDir := filepath.Join(projectRoot, ".agents", "ae")
	tamperedPath := filepath.Join(cacheDir, "commands", "cmd-one.md")
	if err := os.WriteFile(tamperedPath, []byte("tampered"), 0644); err != nil {
		t.Fatalf("failed to tamper file: %v", err)
	}

	inst.SetSource("v2.0.0")
	_, err = inst.RepairCache(ScopeLocal)
	if err != nil {
		t.Fatalf("RepairCache failed: %v", err)
	}

	result, err := inst.CheckCache(ScopeLocal, "v2.0.0")
	if err != nil {
		t.Fatalf("CheckCache failed: %v", err)
	}

	if result.SourceMismatch {
		t.Fatal("source mismatch should be resolved after repair")
	}
	if len(result.Issues) != 0 {
		t.Fatalf("expected no cache issues after repair, got %d", len(result.Issues))
	}

	metadataPath := filepath.Join(cacheDir, CacheMetadataFile)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatalf("failed reading metadata after repair: %v", err)
	}

	var metadata CacheMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		t.Fatalf("failed unmarshaling metadata after repair: %v", err)
	}

	if metadata.Source != "v2.0.0" {
		t.Fatalf("metadata source after repair = %q, want %q", metadata.Source, "v2.0.0")
	}
}
