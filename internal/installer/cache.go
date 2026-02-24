package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

const CacheMetadataFile = "cache-metadata.json"

type CacheFileHash struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type CacheMetadata struct {
	Source string          `json:"source"`
	Files  []CacheFileHash `json:"files"`
}

type CacheIssue struct {
	Kind     string
	Path     string
	Expected string
	Actual   string
}

type CacheCheckResult struct {
	Scope          Scope
	Path           string
	Exists         bool
	MetadataFound  bool
	Source         string
	SourceMismatch bool
	Issues         []CacheIssue
}

type CacheRepairResult struct {
	Scope Scope
	Path  string
	Files int
}

var errCacheMetadataMissing = errors.New("cache metadata missing")

func (i *Installer) CacheDir(scope Scope) string {
	return i.cacheDir(scope)
}

func (i *Installer) RepairCache(scope Scope) (*CacheRepairResult, error) {
	cacheDir := i.cacheDir(scope)
	parent := filepath.Dir(cacheDir)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return nil, fmt.Errorf("creating cache parent: %w", err)
	}

	expected, err := i.expectedCacheFiles()
	if err != nil {
		return nil, err
	}

	tmpDir, err := os.MkdirTemp(parent, ".ae-cache-")
	if err != nil {
		return nil, fmt.Errorf("creating temp cache dir: %w", err)
	}
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = os.RemoveAll(tmpDir)
		}
	}()

	for _, f := range expected {
		if err := writeFile(filepath.Join(tmpDir, f.Path), f.Data); err != nil {
			return nil, fmt.Errorf("writing cache file %s: %w", f.Path, err)
		}
	}

	metadata := CacheMetadata{Source: i.Source}
	metadata.Files = make([]CacheFileHash, 0, len(expected))
	for _, f := range expected {
		metadata.Files = append(metadata.Files, CacheFileHash{Path: f.Path, SHA256: f.SHA256})
	}
	sort.Slice(metadata.Files, func(a, b int) bool {
		return metadata.Files[a].Path < metadata.Files[b].Path
	})

	metaData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling cache metadata: %w", err)
	}
	metaData = append(metaData, '\n')
	if err := writeFile(filepath.Join(tmpDir, CacheMetadataFile), metaData); err != nil {
		return nil, fmt.Errorf("writing cache metadata: %w", err)
	}

	backupDir := cacheDir + ".old"
	_ = os.RemoveAll(backupDir)
	if _, err := os.Stat(cacheDir); err == nil {
		if err := os.Rename(cacheDir, backupDir); err != nil {
			return nil, fmt.Errorf("moving existing cache: %w", err)
		}
	}

	if err := os.Rename(tmpDir, cacheDir); err != nil {
		_ = os.Rename(backupDir, cacheDir)
		return nil, fmt.Errorf("activating rebuilt cache: %w", err)
	}
	cleanupTmp = false
	_ = os.RemoveAll(backupDir)

	return &CacheRepairResult{Scope: scope, Path: cacheDir, Files: len(expected)}, nil
}

func (i *Installer) CheckCache(scope Scope, currentSource string) (*CacheCheckResult, error) {
	cacheDir := i.cacheDir(scope)
	result := &CacheCheckResult{Scope: scope, Path: cacheDir}

	if _, err := os.Stat(cacheDir); err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, fmt.Errorf("checking cache dir: %w", err)
	}
	result.Exists = true

	expected, err := i.expectedCacheFiles()
	if err != nil {
		return nil, err
	}
	expectedMap := make(map[string]string, len(expected))
	for _, entry := range expected {
		expectedMap[entry.Path] = entry.SHA256
	}

	metadata, err := readCacheMetadata(cacheDir)
	if err != nil {
		if errors.Is(err, errCacheMetadataMissing) {
			result.Issues = append(result.Issues, CacheIssue{Kind: "metadata missing", Path: CacheMetadataFile})
		} else {
			result.Issues = append(result.Issues, CacheIssue{Kind: "metadata invalid", Path: CacheMetadataFile, Actual: err.Error()})
		}
	} else {
		result.MetadataFound = true
		result.Source = metadata.Source
		result.SourceMismatch = metadata.Source != "" && metadata.Source != currentSource

		installedMap := make(map[string]string, len(metadata.Files))
		for _, entry := range metadata.Files {
			installedMap[entry.Path] = entry.SHA256
		}

		for path, expectedHash := range expectedMap {
			installedHash, ok := installedMap[path]
			if !ok {
				result.Issues = append(result.Issues, CacheIssue{Kind: "missing metadata hash", Path: path, Expected: expectedHash})
				continue
			}
			if installedHash != expectedHash {
				result.Issues = append(result.Issues, CacheIssue{Kind: "metadata hash mismatch", Path: path, Expected: expectedHash, Actual: installedHash})
			}
		}

		for path := range installedMap {
			if _, ok := expectedMap[path]; !ok {
				result.Issues = append(result.Issues, CacheIssue{Kind: "unexpected metadata file", Path: path})
			}
		}
	}

	for _, entry := range expected {
		path := filepath.Join(cacheDir, entry.Path)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				result.Issues = append(result.Issues, CacheIssue{Kind: "missing file", Path: entry.Path, Expected: entry.SHA256})
				continue
			}
			result.Issues = append(result.Issues, CacheIssue{Kind: "file read error", Path: entry.Path, Actual: err.Error()})
			continue
		}
		actualHash := sha256Hex(data)
		if actualHash != entry.SHA256 {
			result.Issues = append(result.Issues, CacheIssue{Kind: "hash mismatch", Path: entry.Path, Expected: entry.SHA256, Actual: actualHash})
		}
	}

	unexpected, err := findUnexpectedCacheFiles(cacheDir, expectedMap)
	if err != nil {
		return nil, err
	}
	for _, path := range unexpected {
		result.Issues = append(result.Issues, CacheIssue{Kind: "unexpected file", Path: path})
	}

	return result, nil
}

func (i *Installer) writeCacheMetadata(cacheDir string) error {
	expected, err := i.expectedCacheFiles()
	if err != nil {
		return err
	}

	files := make([]CacheFileHash, 0, len(expected))
	for _, entry := range expected {
		data, err := os.ReadFile(filepath.Join(cacheDir, entry.Path))
		if err != nil {
			return fmt.Errorf("reading cache file %s: %w", entry.Path, err)
		}
		files = append(files, CacheFileHash{
			Path:   entry.Path,
			SHA256: sha256Hex(data),
		})
	}

	sort.Slice(files, func(a, b int) bool {
		return files[a].Path < files[b].Path
	})

	metadata := CacheMetadata{
		Source: i.Source,
		Files:  files,
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling cache metadata: %w", err)
	}
	data = append(data, '\n')
	return writeFile(filepath.Join(cacheDir, CacheMetadataFile), data)
}

type expectedCacheFile struct {
	Path   string
	Data   []byte
	SHA256 string
}

func (i *Installer) expectedCacheFiles() ([]expectedCacheFile, error) {
	commands := i.Registry.GetAllCommands()
	skills := i.Registry.GetAllSkills()
	sort.Strings(commands)
	sort.Strings(skills)

	files := []expectedCacheFile{{
		Path:   "README.md",
		Data:   []byte(cacheReadme),
		SHA256: sha256Hex([]byte(cacheReadme)),
	}}

	for _, cmd := range commands {
		data, err := i.Registry.ReadCommand(cmd)
		if err != nil {
			return nil, fmt.Errorf("reading command %s: %w", cmd, err)
		}
		path := filepath.Join("commands", cmd+".md")
		files = append(files, expectedCacheFile{Path: filepath.ToSlash(path), Data: data, SHA256: sha256Hex(data)})
	}

	for _, skill := range skills {
		sourceRoot := i.Registry.SkillSourcePath(skill)
		err := fs.WalkDir(i.Registry.FS, sourceRoot, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			data, err := fs.ReadFile(i.Registry.FS, path)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(sourceRoot, path)
			if err != nil {
				return err
			}
			cachePath := filepath.ToSlash(filepath.Join("skills", skill, rel))
			files = append(files, expectedCacheFile{Path: cachePath, Data: data, SHA256: sha256Hex(data)})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("reading skill %s files: %w", skill, err)
		}
	}

	sort.Slice(files, func(a, b int) bool {
		return files[a].Path < files[b].Path
	})

	return files, nil
}

func readCacheMetadata(cacheDir string) (*CacheMetadata, error) {
	path := filepath.Join(cacheDir, CacheMetadataFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errCacheMetadataMissing
		}
		return nil, err
	}

	var metadata CacheMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func findUnexpectedCacheFiles(cacheDir string, expected map[string]string) ([]string, error) {
	unexpected := make([]string, 0)
	roots := []string{"commands", "skills"}

	for _, root := range roots {
		rootPath := filepath.Join(cacheDir, root)
		if _, err := os.Stat(rootPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("checking %s: %w", rootPath, err)
		}

		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(cacheDir, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if _, ok := expected[rel]; !ok {
				unexpected = append(unexpected, rel)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(unexpected)
	return unexpected, nil
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
