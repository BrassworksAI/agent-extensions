package sdd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type MergeOptions struct {
	RepoRootAbs string
	ChangeName  string
	DryRun      bool
	SpecsDir    string
}

type MergeSummary struct {
	Change   string   `json:"change"`
	DryRun   bool     `json:"dryRun"`
	Counts   Counts   `json:"counts"`
	Created  []string `json:"created"`
	Modified []string `json:"modified"`
	Skipped  []string `json:"skipped"`
}

type Counts struct {
	Created  int `json:"created"`
	Modified int `json:"modified"`
	Skipped  int `json:"skipped"`
}

type mergePlanItem struct {
	CanonicalRelPath  string
	ChangeSpecRelPath string
	Kind              string
}

func MergeChangeSpecs(opts MergeOptions) (MergeSummary, error) {
	if opts.ChangeName == "" {
		return MergeSummary{}, fmt.Errorf("missing change name")
	}
	specsDir := strings.TrimSpace(opts.SpecsDir)
	if specsDir == "" {
		specsDir = "docs/specs"
	}
	if filepath.IsAbs(specsDir) {
		return MergeSummary{}, fmt.Errorf("refusing absolute --specs-dir path: %s", specsDir)
	}
	if strings.Contains(specsDir, "..") {
		return MergeSummary{}, fmt.Errorf("refusing unsafe --specs-dir path: %s", specsDir)
	}

	changeDirAbs := filepath.Join(opts.RepoRootAbs, "changes", opts.ChangeName)
	changeSpecsDirAbs := filepath.Join(changeDirAbs, "specs")
	if _, err := os.Stat(changeSpecsDirAbs); err != nil {
		return MergeSummary{}, fmt.Errorf("missing change specs directory: %s", filepath.ToSlash(filepath.Join("changes", opts.ChangeName, "specs")))
	}

	changeSpecAbsPaths, err := listMarkdownFiles(changeSpecsDirAbs)
	if err != nil {
		return MergeSummary{}, fmt.Errorf("listing change specs: %w", err)
	}

	plan := make([]mergePlanItem, 0, len(changeSpecAbsPaths))
	for _, changeSpecAbs := range changeSpecAbsPaths {
		canonicalRelPath, err := computeCanonicalRelPath(opts.ChangeName, changeSpecAbs, opts.RepoRootAbs, specsDir)
		if err != nil {
			return MergeSummary{}, err
		}
		relChange, err := filepath.Rel(opts.RepoRootAbs, changeSpecAbs)
		if err != nil {
			return MergeSummary{}, fmt.Errorf("computing change spec path for %s: %w", changeSpecAbs, err)
		}
		plan = append(plan, mergePlanItem{
			CanonicalRelPath:  canonicalRelPath,
			ChangeSpecRelPath: filepath.ToSlash(relChange),
			Kind:              "new",
		})
	}

	for i := range plan {
		item := &plan[i]
		changeSpecAbs := filepath.Join(opts.RepoRootAbs, filepath.FromSlash(item.ChangeSpecRelPath))
		changeMd, err := os.ReadFile(changeSpecAbs)
		if err != nil {
			return MergeSummary{}, fmt.Errorf("reading %s: %w", item.ChangeSpecRelPath, err)
		}

		validation := ValidateChangeSpecMarkdown(string(changeMd))
		if !validation.OK {
			var b strings.Builder
			b.WriteString(fmt.Sprintf("invalid change-set spec: %s", item.ChangeSpecRelPath))
			for _, is := range validation.Issues {
				b.WriteString("\n- [")
				b.WriteString(is.Code)
				b.WriteString("] ")
				b.WriteString(is.Message)
			}
			return MergeSummary{}, fmt.Errorf("%s", b.String())
		}

		frontmatterRaw, _ := splitFrontmatter(string(changeMd))
		kind := parseKind(frontmatterRaw)
		if kind == "" {
			return MergeSummary{}, fmt.Errorf("missing required 'kind: new|delta' in frontmatter: %s", item.ChangeSpecRelPath)
		}
		item.Kind = kind
	}

	sort.Slice(plan, func(i, j int) bool {
		return plan[i].CanonicalRelPath < plan[j].CanonicalRelPath
	})

	summary := MergeSummary{
		Change:   opts.ChangeName,
		DryRun:   opts.DryRun,
		Created:  []string{},
		Modified: []string{},
		Skipped:  []string{},
	}

	for _, item := range plan {
		changeSpecAbs := filepath.Join(opts.RepoRootAbs, filepath.FromSlash(item.ChangeSpecRelPath))
		canonicalAbs := filepath.Join(opts.RepoRootAbs, filepath.FromSlash(item.CanonicalRelPath))

		changeMdBytes, err := os.ReadFile(changeSpecAbs)
		if err != nil {
			return MergeSummary{}, fmt.Errorf("reading %s: %w", item.ChangeSpecRelPath, err)
		}
		_, body := splitFrontmatter(string(changeMdBytes))

		if item.Kind == "new" {
			nextBody := normalizeNewSpecBody(body)
			if _, err := os.Stat(canonicalAbs); err == nil {
				if !opts.DryRun {
					if err := ensureParentDir(canonicalAbs); err != nil {
						return MergeSummary{}, err
					}
					if err := os.WriteFile(canonicalAbs, []byte(nextBody), 0644); err != nil {
						return MergeSummary{}, fmt.Errorf("writing %s: %w", item.CanonicalRelPath, err)
					}
				}
				summary.Modified = append(summary.Modified, item.CanonicalRelPath)
			} else {
				if !opts.DryRun {
					if err := ensureParentDir(canonicalAbs); err != nil {
						return MergeSummary{}, err
					}
					if err := os.WriteFile(canonicalAbs, []byte(nextBody), 0644); err != nil {
						return MergeSummary{}, fmt.Errorf("writing %s: %w", item.CanonicalRelPath, err)
					}
				}
				summary.Created = append(summary.Created, item.CanonicalRelPath)
			}
			continue
		}

		if _, err := os.Stat(canonicalAbs); err != nil {
			return MergeSummary{}, fmt.Errorf("delta targets missing canonical spec: %s", item.CanonicalRelPath)
		}

		canonicalMdBytes, err := os.ReadFile(canonicalAbs)
		if err != nil {
			return MergeSummary{}, fmt.Errorf("reading %s: %w", item.CanonicalRelPath, err)
		}
		patched, err := patchCanonicalWithDelta(string(canonicalMdBytes), body, item.ChangeSpecRelPath)
		if err != nil {
			return MergeSummary{}, err
		}

		if !opts.DryRun {
			if err := os.WriteFile(canonicalAbs, []byte(patched), 0644); err != nil {
				return MergeSummary{}, fmt.Errorf("writing %s: %w", item.CanonicalRelPath, err)
			}
		}
		summary.Modified = append(summary.Modified, item.CanonicalRelPath)
	}

	sort.Strings(summary.Created)
	sort.Strings(summary.Modified)
	sort.Strings(summary.Skipped)
	summary.Counts = Counts{Created: len(summary.Created), Modified: len(summary.Modified), Skipped: len(summary.Skipped)}

	return summary, nil
}

func computeCanonicalRelPath(changeName, changeSpecAbs, repoRootAbs, specsDir string) (string, error) {
	changeSpecsDirAbs := filepath.Join(repoRootAbs, "changes", changeName, "specs")
	relFromChangeSpecs, err := filepath.Rel(changeSpecsDirAbs, changeSpecAbs)
	if err != nil {
		return "", fmt.Errorf("computing path for %s: %w", changeSpecAbs, err)
	}
	if strings.HasPrefix(relFromChangeSpecs, "..") {
		return "", fmt.Errorf("change spec is outside expected directory: %s", changeSpecAbs)
	}

	return path.Join(filepath.ToSlash(specsDir), filepath.ToSlash(relFromChangeSpecs)), nil
}

func normalizeNewSpecBody(changeSpecBody string) string {
	return regexp.MustCompile(`^\s+`).ReplaceAllString(changeSpecBody, "")
}

func extractDeltaBuckets(deltaBody, fileLabel string) (added string, modified string, removed string, err error) {
	requirementsHeader := "## Requirements"
	requirementsIndex := strings.Index(deltaBody, requirementsHeader)
	if requirementsIndex == -1 {
		return "", "", "", fmt.Errorf("delta spec missing '## Requirements': %s", fileLabel)
	}

	nextSectionIndex := strings.Index(deltaBody[requirementsIndex+len(requirementsHeader):], "\n## ")
	requirementsSection := ""
	if nextSectionIndex == -1 {
		requirementsSection = deltaBody[requirementsIndex:]
	} else {
		requirementsSection = deltaBody[requirementsIndex : requirementsIndex+len(requirementsHeader)+nextSectionIndex+1]
	}

	addedIndex := strings.Index(requirementsSection, "### ADDED")
	modifiedIndex := strings.Index(requirementsSection, "### MODIFIED")
	removedIndex := strings.Index(requirementsSection, "### REMOVED")
	if modifiedIndex == -1 && addedIndex == -1 && removedIndex == -1 {
		return "", "", "", fmt.Errorf("delta spec has no ADDED/MODIFIED/REMOVED buckets under Requirements: %s", fileLabel)
	}

	sliceBucket := func(start int, candidates []int) string {
		if start == -1 {
			return ""
		}
		end := len(requirementsSection)
		for _, n := range candidates {
			if n != -1 && n > start && n < end {
				end = n
			}
		}
		return strings.TrimRight(requirementsSection[start:end], "\n") + "\n"
	}

	return sliceBucket(addedIndex, []int{modifiedIndex, removedIndex}),
		sliceBucket(modifiedIndex, []int{addedIndex, removedIndex}),
		sliceBucket(removedIndex, []int{addedIndex, modifiedIndex}), nil
}

func patchCanonicalWithDelta(canonicalBody, deltaBody, fileLabel string) (string, error) {
	added, modified, removed, err := extractDeltaBuckets(deltaBody, fileLabel)
	if err != nil {
		return "", err
	}

	next := canonicalBody

	removeExactBlocks := func(bucketContent string) (map[string]struct{}, error) {
		touchedTopics := make(map[string]struct{})
		if bucketContent == "" {
			return touchedTopics, nil
		}

		body := strings.TrimSpace(regexp.MustCompile(`(?m)^###\s+REMOVED\s*\n?`).ReplaceAllString(bucketContent, ""))
		topicChunks := splitTopicChunks(body)
		for _, chunk := range topicChunks {
			lines := strings.Split(chunk, "\n")
			if len(lines) == 0 {
				continue
			}
			topic := strings.TrimSpace(lines[0])
			if topic != "" {
				touchedTopics[topic] = struct{}{}
			}

			remainder := strings.TrimSpace(strings.Join(lines[1:], "\n"))
			if remainder == "" {
				continue
			}

			reasonIndex := strings.Index(remainder, "**Reason:**")
			block := strings.TrimSpace(remainder)
			if reasonIndex != -1 {
				block = strings.TrimSpace(remainder[:reasonIndex])
			}
			if block == "" {
				continue
			}

			normalizedBlock := strings.ReplaceAll(block, "\r\n", "\n")
			pos := strings.Index(next, normalizedBlock)
			if pos == -1 {
				return nil, fmt.Errorf("REMOVED block not found in canonical spec for: %s", fileLabel)
			}
			next = next[:pos] + next[pos+len(normalizedBlock):]
		}

		return touchedTopics, nil
	}

	cleanupEmptyTopics := func(topics map[string]struct{}) {
		for topic := range topics {
			topicHeader := "### " + topic
			headerIndex := strings.Index(next, topicHeader)
			if headerIndex == -1 {
				continue
			}

			sectionStart := headerIndex
			contentStart := headerIndex + len(topicHeader)
			followingIndex := strings.Index(next[contentStart:], "\n### ")
			sectionEnd := len(next)
			if followingIndex != -1 {
				sectionEnd = contentStart + followingIndex + 1
			}

			sectionBody := next[contentStart:sectionEnd]
			hasBullet := regexp.MustCompile(`(?m)^\s*-\s+`).MatchString(sectionBody)
			if hasBullet {
				continue
			}

			before := regexp.MustCompile(`\n{3,}$`).ReplaceAllString(next[:sectionStart], "\n\n")
			after := regexp.MustCompile(`^\n{3,}`).ReplaceAllString(next[sectionEnd:], "\n\n")
			next = before + after
		}
	}

	insertAddedBlocks := func(bucketContent string) error {
		if bucketContent == "" {
			return nil
		}

		body := strings.TrimSpace(regexp.MustCompile(`(?m)^###\s+ADDED\s*\n?`).ReplaceAllString(bucketContent, ""))
		topicChunks := splitTopicChunks(body)
		for _, chunk := range topicChunks {
			lines := strings.Split(chunk, "\n")
			if len(lines) == 0 {
				continue
			}
			topic := strings.TrimSpace(lines[0])
			insertion := strings.TrimSpace(strings.Join(lines[1:], "\n"))
			if topic == "" || insertion == "" {
				continue
			}

			topicHeader := "### " + topic
			topicIndex := strings.Index(next, topicHeader)
			if topicIndex == -1 {
				return fmt.Errorf("ADDED topic not found in canonical spec (%s) for: %s", topicHeader, fileLabel)
			}

			searchStart := topicIndex + len(topicHeader)
			nextTopicIndex := strings.Index(next[searchStart:], "\n### ")
			insertAt := len(next)
			if nextTopicIndex != -1 {
				insertAt = searchStart + nextTopicIndex + 1
			}

			toInsert := strings.TrimRight(strings.ReplaceAll(insertion, "\r\n", "\n"), "\n")
			next = next[:insertAt] + "\n" + toInsert + "\n" + next[insertAt:]
		}
		return nil
	}

	removedTopics, err := removeExactBlocks(removed)
	if err != nil {
		return "", err
	}

	if modified != "" {
		modifiedBody := strings.TrimSpace(regexp.MustCompile(`(?m)^###\s+MODIFIED\s*\n?`).ReplaceAllString(modified, ""))
		topicChunks := splitTopicChunks(modifiedBody)
		for _, chunk := range topicChunks {
			beforeLabelIndex := strings.Index(chunk, "**Before:**")
			afterLabelIndex := strings.Index(chunk, "**After:**")
			if beforeLabelIndex == -1 || afterLabelIndex == -1 || afterLabelIndex < beforeLabelIndex {
				return "", fmt.Errorf("invalid MODIFIED topic; expected **Before:** then **After:** in: %s", fileLabel)
			}

			beforeBlock := strings.TrimSpace(chunk[beforeLabelIndex+len("**Before:**") : afterLabelIndex])
			afterBlock := strings.TrimSpace(chunk[afterLabelIndex+len("**After:**"):])
			if beforeBlock == "" || afterBlock == "" {
				return "", fmt.Errorf("empty Before/After block in MODIFIED topic: %s", fileLabel)
			}

			normalizedBefore := strings.ReplaceAll(beforeBlock, "\r\n", "\n")
			normalizedAfter := strings.ReplaceAll(afterBlock, "\r\n", "\n")

			beforePos := strings.Index(next, normalizedBefore)
			if beforePos == -1 {
				return "", fmt.Errorf("Before block not found in canonical spec for: %s", fileLabel)
			}

			next = next[:beforePos] + normalizedAfter + next[beforePos+len(normalizedBefore):]
		}
	}

	if err := insertAddedBlocks(added); err != nil {
		return "", err
	}

	cleanupEmptyTopics(removedTopics)
	return next, nil
}

func ensureParentDir(fileAbs string) error {
	if err := os.MkdirAll(filepath.Dir(fileAbs), 0755); err != nil {
		return fmt.Errorf("creating parent dir for %s: %w", fileAbs, err)
	}
	return nil
}
