package sdd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeChangeSpecs_NewAndDelta(t *testing.T) {
	repo := t.TempDir()
	change := "auth-refresh"
	changeSpecsDir := filepath.Join(repo, "changes", change, "specs", "auth")
	if err := os.MkdirAll(changeSpecsDir, 0755); err != nil {
		t.Fatalf("mkdir change specs: %v", err)
	}
	canonicalDir := filepath.Join(repo, "specs", "auth")
	if err := os.MkdirAll(canonicalDir, 0755); err != nil {
		t.Fatalf("mkdir canonical: %v", err)
	}

	canonicalLogin := `# Login

## Overview

Login capability.

## Requirements

### Entry

- WHEN a user navigates to login the system SHALL render the login page.

### Exit

- WHEN sign-out completes the system SHALL redirect to home.
`
	if err := os.WriteFile(filepath.Join(canonicalDir, "login.md"), []byte(canonicalLogin), 0644); err != nil {
		t.Fatalf("write canonical login: %v", err)
	}

	delta := `---
kind: delta
---
# Login

## Requirements

### MODIFIED

#### Entry

**Before:**
- WHEN a user navigates to login the system SHALL render the login page.

**After:**
- WHEN a user navigates to login the system SHALL render the sign-in page.

### ADDED

#### Exit

- WHEN sign-in succeeds the system SHALL redirect to dashboard.
`
	if err := os.WriteFile(filepath.Join(changeSpecsDir, "login.md"), []byte(delta), 0644); err != nil {
		t.Fatalf("write delta: %v", err)
	}

	newSpec := `---
kind: new
---
# Logout

## Overview

Logout capability.

## Requirements

### Entry

- WHEN a user initiates sign-out the system SHALL terminate the session.
`
	if err := os.WriteFile(filepath.Join(changeSpecsDir, "logout.md"), []byte(newSpec), 0644); err != nil {
		t.Fatalf("write new spec: %v", err)
	}

	summary, err := MergeChangeSpecs(MergeOptions{
		RepoRootAbs: repo,
		ChangeName:  change,
		DryRun:      false,
		SpecsDir:    "specs",
	})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}

	if summary.Counts.Created != 1 || summary.Counts.Modified != 1 {
		t.Fatalf("unexpected counts: %+v", summary.Counts)
	}

	mergedLoginBytes, err := os.ReadFile(filepath.Join(canonicalDir, "login.md"))
	if err != nil {
		t.Fatalf("read merged login: %v", err)
	}
	mergedLogin := string(mergedLoginBytes)
	if !strings.Contains(mergedLogin, "render the sign-in page") {
		t.Fatal("expected modified login requirement")
	}
	if !strings.Contains(mergedLogin, "redirect to dashboard") {
		t.Fatal("expected added login requirement")
	}

	if _, err := os.Stat(filepath.Join(canonicalDir, "logout.md")); err != nil {
		t.Fatalf("expected created spec file: %v", err)
	}
}

func TestMergeChangeSpecs_DryRun(t *testing.T) {
	repo := t.TempDir()
	change := "auth-refresh"
	changeSpecsDir := filepath.Join(repo, "changes", change, "specs")
	if err := os.MkdirAll(changeSpecsDir, 0755); err != nil {
		t.Fatalf("mkdir change specs: %v", err)
	}

	newSpec := `---
kind: new
---
# Session

## Overview

Session capability.

## Requirements

- The system SHALL maintain session state.
`
	if err := os.WriteFile(filepath.Join(changeSpecsDir, "session.md"), []byte(newSpec), 0644); err != nil {
		t.Fatalf("write new spec: %v", err)
	}

	summary, err := MergeChangeSpecs(MergeOptions{
		RepoRootAbs: repo,
		ChangeName:  change,
		DryRun:      true,
		SpecsDir:    "docs/specs",
	})
	if err != nil {
		t.Fatalf("dry-run merge: %v", err)
	}
	if summary.Counts.Created != 1 {
		t.Fatalf("expected one created in dry run, got %+v", summary.Counts)
	}
	if _, err := os.Stat(filepath.Join(repo, "docs", "specs", "session.md")); !os.IsNotExist(err) {
		t.Fatalf("expected no file written in dry run, stat err: %v", err)
	}
}
