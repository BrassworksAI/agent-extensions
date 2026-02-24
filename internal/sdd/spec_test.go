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
