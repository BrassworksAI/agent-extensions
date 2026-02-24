package sdd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateChangeSpecMarkdown_NewSpecOK(t *testing.T) {
	md := `---
kind: new
---
# Login

## Overview

Login capability.

## Requirements

### Entry

- WHEN a user navigates to login the system SHALL render the login page.
`

	result := ValidateChangeSpecMarkdown(md)
	if !result.OK {
		t.Fatalf("expected validation to pass, got issues: %+v", result.Issues)
	}
	if result.Kind != "new" {
		t.Fatalf("expected kind=new, got %s", result.Kind)
	}
}

func TestValidateChangeSpecMarkdown_DeltaMalformedModified(t *testing.T) {
	md := `---
kind: delta
---
# Login

## Requirements

### MODIFIED

#### Entry

**Before:**
- WHEN a user navigates to login the system SHALL render login.
`

	result := ValidateChangeSpecMarkdown(md)
	if result.OK {
		t.Fatal("expected validation to fail")
	}
	found := false
	for _, is := range result.Issues {
		if is.Code == "delta.modified_malformed" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected delta.modified_malformed issue, got %+v", result.Issues)
	}
}

func TestValidateAllChangeSpecs(t *testing.T) {
	repo := t.TempDir()
	changeDir := filepath.Join(repo, "changes", "auth-refresh")
	if err := os.MkdirAll(filepath.Join(changeDir, "specs", "auth"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	okFile := filepath.Join(changeDir, "specs", "auth", "login.md")
	badFile := filepath.Join(changeDir, "specs", "auth", "logout.md")

	if err := os.WriteFile(okFile, []byte(`---
kind: new
---
# Login

## Overview

text

## Requirements

- The system SHALL do thing.
`), 0644); err != nil {
		t.Fatalf("write ok file: %v", err)
	}

	if err := os.WriteFile(badFile, []byte(`# Missing frontmatter`), 0644); err != nil {
		t.Fatalf("write bad file: %v", err)
	}

	results, err := ValidateAllChangeSpecs(repo, changeDir)
	if err != nil {
		t.Fatalf("validate all: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	var seenOK, seenBad bool
	for _, r := range results {
		if r.Path == "changes/auth-refresh/specs/auth/login.md" {
			seenOK = true
			if !r.OK {
				t.Fatalf("expected login.md to be valid, got %+v", r.Issues)
			}
		}
		if r.Path == "changes/auth-refresh/specs/auth/logout.md" {
			seenBad = true
			if r.OK {
				t.Fatal("expected logout.md to be invalid")
			}
		}
	}
	if !seenOK || !seenBad {
		t.Fatalf("missing expected result paths: %+v", results)
	}
}

func TestListCanonicalSpecs(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "specs", "auth"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "specs", "auth", "login.md"), []byte("# Login"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "specs", "auth", "notes.txt"), []byte("ignore"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	paths, err := ListCanonicalSpecs(repo, "specs")
	if err != nil {
		t.Fatalf("list canonical specs: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 markdown spec, got %d", len(paths))
	}
	if paths[0] != "specs/auth/login.md" {
		t.Fatalf("unexpected path: %s", paths[0])
	}
}

func TestListCanonicalSpecs_CustomDir(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "docs", "specs"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "docs", "specs", "session.md"), []byte("# Session"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	paths, err := ListCanonicalSpecs(repo, "docs/specs")
	if err != nil {
		t.Fatalf("list canonical specs: %v", err)
	}
	if len(paths) != 1 || paths[0] != "docs/specs/session.md" {
		t.Fatalf("unexpected paths: %+v", paths)
	}
}

func TestListCanonicalSpecs_DefaultDirUsesDocsSpecs(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "docs", "specs", "auth"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "docs", "specs", "auth", "login.md"), []byte("# Login"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	paths, err := ListCanonicalSpecs(repo, "")
	if err != nil {
		t.Fatalf("list canonical specs: %v", err)
	}
	if len(paths) != 1 || paths[0] != "docs/specs/auth/login.md" {
		t.Fatalf("unexpected paths: %+v", paths)
	}
}

func TestListCanonicalSpecs_RejectsUnsafeDir(t *testing.T) {
	repo := t.TempDir()
	if _, err := ListCanonicalSpecs(repo, "../specs"); err == nil {
		t.Fatal("expected unsafe path error")
	}
}
