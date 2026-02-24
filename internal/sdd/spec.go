package sdd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ValidationIssue struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationResult struct {
	OK     bool              `json:"ok"`
	Path   string            `json:"path,omitempty"`
	Kind   string            `json:"kind,omitempty"`
	Issues []ValidationIssue `json:"issues"`
}

func issue(code, message string) ValidationIssue {
	return ValidationIssue{Code: code, Message: message}
}

func splitFrontmatter(markdown string) (frontmatterRaw string, body string) {
	if !strings.HasPrefix(markdown, "---\n") {
		return "", markdown
	}

	endIndex := strings.Index(markdown[4:], "\n---\n")
	if endIndex == -1 {
		return "", markdown
	}
	endIndex += 4

	frontmatterRaw = markdown[4 : endIndex+1]
	body = markdown[endIndex+5:]
	return frontmatterRaw, body
}

func parseKind(frontmatterRaw string) string {
	if frontmatterRaw == "" {
		return ""
	}
	kindRe := regexp.MustCompile(`(?m)^kind:\s*(new|delta)\s*$`)
	match := kindRe.FindStringSubmatch(frontmatterRaw)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func hasHeadingPrefix(body, prefix string) bool {
	escaped := regexp.QuoteMeta(prefix)
	re := regexp.MustCompile(`(?m)^` + escaped + `\s+.+$`)
	return re.MatchString(body)
}

func validateNew(body string) []ValidationIssue {
	issues := make([]ValidationIssue, 0)

	if !hasHeadingPrefix(body, "#") {
		issues = append(issues, issue("new.missing_title", "Missing top-level '# <Title>' heading"))
	}
	if !hasHeadingPrefix(body, "## Overview") {
		issues = append(issues, issue("new.missing_overview", "Missing '## Overview' section"))
	}
	if !hasHeadingPrefix(body, "## Requirements") {
		issues = append(issues, issue("new.missing_requirements", "Missing '## Requirements' section"))
	}

	bucketRe := regexp.MustCompile(`(?m)^###\s+(ADDED|MODIFIED|REMOVED)\s*$`)
	if bucketRe.MatchString(body) {
		issues = append(issues, issue("new.delta_buckets", "`kind: new` spec must not include delta buckets (### ADDED/MODIFIED/REMOVED)"))
	}

	return issues
}

func validateDelta(body string) []ValidationIssue {
	issues := make([]ValidationIssue, 0)

	if !hasHeadingPrefix(body, "#") {
		issues = append(issues, issue("delta.missing_title", "Missing top-level '# <Title>' heading"))
	}
	if !hasHeadingPrefix(body, "## Requirements") {
		issues = append(issues, issue("delta.missing_requirements", "Missing '## Requirements' section"))
	}

	requirementsIndex := strings.Index(body, "## Requirements")
	if requirementsIndex != -1 {
		nextSectionIndex := strings.Index(body[requirementsIndex+len("## Requirements"):], "\n## ")
		requirementsSection := ""
		if nextSectionIndex == -1 {
			requirementsSection = body[requirementsIndex:]
		} else {
			requirementsSection = body[requirementsIndex : requirementsIndex+len("## Requirements")+nextSectionIndex+1]
		}

		hasBucketRe := regexp.MustCompile(`(?m)^(###\s+(ADDED|MODIFIED|REMOVED)\s*)$`)
		if !hasBucketRe.MatchString(requirementsSection) {
			issues = append(issues, issue("delta.missing_buckets", "Delta spec must include at least one bucket under '## Requirements'"))
		}

		if strings.Contains(requirementsSection, "### MODIFIED") {
			modifiedStart := strings.Index(requirementsSection, "### MODIFIED")
			afterModified := requirementsSection[modifiedStart+1:]
			nextAddedRel := strings.Index(afterModified, "### ADDED")
			nextRemovedRel := strings.Index(afterModified, "### REMOVED")
			nextBucket := len(requirementsSection)
			if nextAddedRel != -1 {
				nextAddedAbs := modifiedStart + 1 + nextAddedRel
				if nextAddedAbs < nextBucket {
					nextBucket = nextAddedAbs
				}
			}
			if nextRemovedRel != -1 {
				nextRemovedAbs := modifiedStart + 1 + nextRemovedRel
				if nextRemovedAbs < nextBucket {
					nextBucket = nextRemovedAbs
				}
			}

			modifiedBlock := requirementsSection[modifiedStart:nextBucket]
			topics := splitTopicChunks(strings.TrimSpace(regexp.MustCompile(`(?m)^###\s+MODIFIED\s*\n?`).ReplaceAllString(modifiedBlock, "")))
			for _, topic := range topics {
				beforeIndex := strings.Index(topic, "**Before:**")
				afterIndex := strings.Index(topic, "**After:**")
				if beforeIndex == -1 || afterIndex == -1 || afterIndex < beforeIndex {
					issues = append(issues, issue("delta.modified_malformed", "Each MODIFIED topic must contain '**Before:**' followed by '**After:**'"))
					break
				}
			}
		}
	}

	return issues
}

func ValidateChangeSpecMarkdown(markdown string) ValidationResult {
	frontmatterRaw, body := splitFrontmatter(markdown)
	kind := parseKind(frontmatterRaw)

	issues := make([]ValidationIssue, 0)
	if frontmatterRaw == "" {
		issues = append(issues, issue("fm.missing", "Missing YAML frontmatter (--- ... ---)"))
	}
	if kind == "" {
		issues = append(issues, issue("fm.kind_missing", "Missing required 'kind: new|delta' in frontmatter"))
	}

	if kind == "new" {
		issues = append(issues, validateNew(body)...)
	}
	if kind == "delta" {
		issues = append(issues, validateDelta(body)...)
	}

	return ValidationResult{OK: len(issues) == 0, Kind: kind, Issues: issues}
}

func ValidateChangeSpecPath(repoRootAbs, relPath string) (ValidationResult, error) {
	if filepath.IsAbs(relPath) {
		return ValidationResult{}, fmt.Errorf("refusing absolute path; pass a repo-relative path")
	}
	if strings.Contains(relPath, "..") {
		return ValidationResult{}, fmt.Errorf("refusing path traversal; '..' is not allowed")
	}

	fileAbs := filepath.Join(repoRootAbs, relPath)
	markdown, err := os.ReadFile(fileAbs)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("reading %s: %w", relPath, err)
	}

	result := ValidateChangeSpecMarkdown(string(markdown))
	result.Path = filepath.ToSlash(relPath)
	return result, nil
}

func listMarkdownFiles(rootAbs string) ([]string, error) {
	results := make([]string, 0)
	err := filepath.WalkDir(rootAbs, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".md") {
			results = append(results, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(results)
	return results, nil
}

func ValidateAllChangeSpecs(repoRootAbs, changeDirAbs string) ([]ValidationResult, error) {
	changeSpecsDirAbs := filepath.Join(changeDirAbs, "specs")
	if _, err := os.Stat(changeSpecsDirAbs); err != nil {
		return nil, fmt.Errorf("missing change specs directory: %s", filepath.ToSlash(strings.TrimPrefix(changeSpecsDirAbs, repoRootAbs+string(filepath.Separator))))
	}

	files, err := listMarkdownFiles(changeSpecsDirAbs)
	if err != nil {
		return nil, fmt.Errorf("listing change specs: %w", err)
	}

	results := make([]ValidationResult, 0, len(files))
	for _, fileAbs := range files {
		rel, err := filepath.Rel(repoRootAbs, fileAbs)
		if err != nil {
			return nil, fmt.Errorf("computing path for %s: %w", fileAbs, err)
		}
		result, err := ValidateChangeSpecPath(repoRootAbs, filepath.ToSlash(rel))
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func ListCanonicalSpecs(repoRootAbs, specsDir string) ([]string, error) {
	specsDir = strings.TrimSpace(specsDir)
	if specsDir == "" {
		specsDir = "specs"
	}
	if filepath.IsAbs(specsDir) {
		return nil, fmt.Errorf("refusing absolute --specs-dir path: %s", specsDir)
	}
	if strings.Contains(specsDir, "..") {
		return nil, fmt.Errorf("refusing unsafe --specs-dir path: %s", specsDir)
	}

	specsDirAbs := filepath.Join(repoRootAbs, specsDir)
	if _, err := os.Stat(specsDirAbs); err != nil {
		return nil, fmt.Errorf("missing canonical specs directory: %s", filepath.ToSlash(specsDir))
	}

	files, err := listMarkdownFiles(specsDirAbs)
	if err != nil {
		return nil, fmt.Errorf("listing canonical specs: %w", err)
	}

	paths := make([]string, 0, len(files))
	for _, fileAbs := range files {
		rel, err := filepath.Rel(repoRootAbs, fileAbs)
		if err != nil {
			return nil, fmt.Errorf("computing path for %s: %w", fileAbs, err)
		}
		paths = append(paths, filepath.ToSlash(rel))
	}

	return paths, nil
}

func splitTopicChunks(body string) []string {
	parts := regexp.MustCompile(`(?m)^####\s+`).Split(body, -1)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
