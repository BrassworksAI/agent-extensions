package installer

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shanepadgett/agent-extensions/internal/config"
	"github.com/shanepadgett/agent-extensions/internal/registry"
)

type Scope string

const (
	ScopeGlobal Scope = "global"
	ScopeLocal  Scope = "local"
	ScopeBoth   Scope = "both"
)

type Installer struct {
	Registry    *registry.Registry
	ProjectRoot string
	Source      string
}

type InstallResult struct {
	Tool     string
	Scope    Scope
	Commands int
	Skills   int
	Errors   []error
}

const cacheReadme = `# Agent Extensions Cache

This folder is managed by ae and contains cached skills.
Agent tools symlink to these files.

Do not delete this folder.
`

func New(reg *registry.Registry, projectRoot string) *Installer {
	return &Installer{
		Registry:    reg,
		ProjectRoot: projectRoot,
		Source:      "dev",
	}
}

func (i *Installer) SetSource(source string) {
	if source == "" {
		return
	}
	i.Source = source
}

func (i *Installer) cacheDir(scope Scope) string {
	if scope == ScopeGlobal {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".agents", "ae")
	}
	return filepath.Join(i.ProjectRoot, ".agents", "ae")
}

func (i *Installer) cacheRoot(scope Scope) string {
	if scope == ScopeGlobal {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".agents")
	}
	return filepath.Join(i.ProjectRoot, ".agents")
}

func (i *Installer) Install(toolName string, scope Scope) (*InstallResult, error) {
	tool, ok := i.Registry.GetTool(toolName)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}

	result := &InstallResult{
		Tool:  toolName,
		Scope: scope,
	}

	commands := i.Registry.GetAllCommands()
	skills := i.Registry.GetAllSkills()
	requiresJSRuntime := i.anySkillUsesScripts(skills)

	scopes := []Scope{scope}
	if scope == ScopeBoth {
		scopes = []Scope{ScopeGlobal, ScopeLocal}
	}

	for _, s := range scopes {
		cache := i.cacheDir(s)
		scopeErrStart := len(result.Errors)
		if err := writeFile(filepath.Join(cache, "README.md"), []byte(cacheReadme)); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("cache readme: %w", err))
		}

		var target string
		if s == ScopeGlobal {
			target = tool.ResolveGlobalPath()
		} else {
			target = tool.ResolveLocalPath(i.ProjectRoot)
		}

		// Install commands
		for _, cmd := range commands {
			if err := i.installCommand(cmd, cache, target, tool.Conventions, s == ScopeGlobal); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("command %s: %w", cmd, err))
			} else {
				result.Commands++
			}
		}

		// Install skills
		for _, skill := range skills {
			if err := i.installSkill(skill, cache, target, tool.Conventions, s == ScopeGlobal); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("skill %s: %w", skill, err))
			} else {
				result.Skills++
			}
		}

		if len(result.Errors) == scopeErrStart {
			if err := i.writeCacheMetadata(cache); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("cache metadata: %w", err))
			}
		}
	}

	if requiresJSRuntime && findJSRuntime() == "" {
		result.Errors = append(result.Errors, fmt.Errorf("no JavaScript runtime found (node, bun, deno); skill scripts may fail"))
	}

	return result, nil
}

func (i *Installer) installCommand(name, cacheDir, targetBase string, conv config.Conventions, isGlobal bool) error {
	if !i.Registry.CommandExists(name) {
		return fmt.Errorf("command not found: %s", name)
	}

	// Read from embedded FS
	data, err := i.Registry.ReadCommand(name)
	if err != nil {
		return fmt.Errorf("reading command: %w", err)
	}

	// Write to cache
	cacheDest := filepath.Join(cacheDir, "commands", name+".md")
	if err := writeFile(cacheDest, data); err != nil {
		return fmt.Errorf("writing to cache: %w", err)
	}

	// Symlink from tool location to cache
	destPath := conv.CommandPath(name, isGlobal)
	dest := filepath.Join(targetBase, destPath)

	return createSymlink(cacheDest, dest, !isGlobal)
}

func (i *Installer) installSkill(name, cacheDir, targetBase string, conv config.Conventions, isGlobal bool) error {
	if !i.Registry.SkillExists(name) {
		return fmt.Errorf("skill not found: %s", name)
	}

	// Copy skill directory from embedded FS to cache
	cacheDest := filepath.Join(cacheDir, "skills", name)
	if err := i.copySkillToCache(name, cacheDest); err != nil {
		return fmt.Errorf("copying to cache: %w", err)
	}

	destPath := conv.ScopedSkillPath(name, isGlobal)
	hasScopeOverride := (isGlobal && conv.GlobalSkills != "") || (!isGlobal && conv.LocalSkills != "")
	dest := resolveTargetPath(targetBase, destPath, i.ProjectRoot, isGlobal, hasScopeOverride)

	// Directory-based skill conventions should link the entire skill folder.
	// Single-file conventions still link just SKILL.md.
	isSingleFile := filepath.Ext(destPath) == ".md" && filepath.Base(destPath) == name+".md"
	if isSingleFile {
		skillFile := filepath.Join(cacheDest, "SKILL.md")
		return createSymlink(skillFile, dest, !isGlobal)
	}

	return createSymlink(cacheDest, filepath.Dir(dest), !isGlobal)
}

func (i *Installer) copySkillToCache(skillName, cacheDest string) error {
	if err := os.RemoveAll(cacheDest); err != nil {
		return err
	}

	srcPath := i.Registry.SkillSourcePath(skillName)
	return fs.WalkDir(i.Registry.FS, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(srcPath, path)
		destPath := filepath.Join(cacheDest, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := fs.ReadFile(i.Registry.FS, path)
		if err != nil {
			return err
		}

		return writeFile(destPath, data)
	})
}

func (i *Installer) anySkillUsesScripts(skills []string) bool {
	for _, skill := range skills {
		if i.skillHasScripts(skill) {
			return true
		}
	}
	return false
}

func (i *Installer) skillHasScripts(skill string) bool {
	entries, err := i.Registry.ListSkillFiles(skill)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == "scripts" {
			return true
		}
	}
	return false
}

func findJSRuntime() string {
	for _, runtime := range []string{"node", "bun", "deno"} {
		if _, err := exec.LookPath(runtime); err == nil {
			return runtime
		}
	}
	return ""
}

func createSymlink(src, dest string, useRelative bool) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("creating parent dir: %w", err)
	}

	if _, err := os.Lstat(dest); err == nil {
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("removing existing: %w", err)
		}
	}

	linkTarget := src
	if useRelative {
		rel, err := filepath.Rel(filepath.Dir(dest), src)
		if err != nil {
			return fmt.Errorf("resolving relative symlink: %w", err)
		}
		linkTarget = rel
	}

	if err := os.Symlink(linkTarget, dest); err != nil {
		return fmt.Errorf("creating symlink: %w", err)
	}

	return nil
}

func writeFile(dest string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	return os.WriteFile(dest, data, 0644)
}

func (i *Installer) Uninstall(toolName string, scope Scope) (*InstallResult, error) {
	tool, ok := i.Registry.GetTool(toolName)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}

	result := &InstallResult{
		Tool:  toolName,
		Scope: scope,
	}

	commands := i.Registry.GetAllCommands()
	skills := i.Registry.GetAllSkills()

	scopes := []Scope{scope}
	if scope == ScopeBoth {
		scopes = []Scope{ScopeGlobal, ScopeLocal}
	}

	for _, s := range scopes {
		var target string
		if s == ScopeGlobal {
			target = tool.ResolveGlobalPath()
		} else {
			target = tool.ResolveLocalPath(i.ProjectRoot)
		}

		cache := i.cacheDir(s)

		// Uninstall commands
		for _, cmd := range commands {
			destPath := tool.Conventions.CommandPath(cmd, s == ScopeGlobal)
			dest := filepath.Join(target, destPath)
			if _, err := os.Lstat(dest); err == nil {
				if err := os.RemoveAll(dest); err == nil {
					result.Commands++
					cleanEmptyParents(filepath.Dir(dest), target)
				}
			}
			// Also remove from cache
			cachePath := filepath.Join(cache, "commands", cmd+".md")
			os.RemoveAll(cachePath)
		}

		// Uninstall skills
		for _, skill := range skills {
			destPath := tool.Conventions.ScopedSkillPath(skill, s == ScopeGlobal)
			hasScopeOverride := (s == ScopeGlobal && tool.Conventions.GlobalSkills != "") || (s != ScopeGlobal && tool.Conventions.LocalSkills != "")
			dest := resolveTargetPath(target, destPath, i.ProjectRoot, s == ScopeGlobal, hasScopeOverride)

			// Check if this is a single-file skill pattern (e.g., skills/{name}.md)
			// vs directory-based (e.g., skills/{name}/SKILL.md)
			isSingleFile := filepath.Ext(destPath) == ".md" && filepath.Base(destPath) == skill+".md"

			if isSingleFile {
				// Single-file skill - remove the .md file
				if _, err := os.Lstat(dest); err == nil {
					if err := os.RemoveAll(dest); err == nil {
						result.Skills++
						cleanEmptyParents(filepath.Dir(dest), target)
					}
				}
			} else {
				// Directory-based skill - remove the skill directory (contains symlinked files)
				// e.g., for skills/{name}/SKILL.md, remove skills/{name}
				skillDir := filepath.Dir(dest)
				if _, err := os.Lstat(skillDir); err == nil {
					if err := os.RemoveAll(skillDir); err == nil {
						result.Skills++
						cleanEmptyParents(filepath.Dir(skillDir), target)
					}
				}
			}
			// Also remove from cache
			cachePath := filepath.Join(cache, "skills", skill)
			os.RemoveAll(cachePath)
		}

		// Clean up cache directory and remove .agents only if empty
		cleanEmptyParents(cache, i.cacheRoot(s))
	}

	return result, nil
}

// cleanEmptyParents removes empty directories from dir up to and including stopAt
func cleanEmptyParents(dir, stopAt string) {
	for {
		if dir == "/" || dir == "." {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
		// Stop after removing stopAt to avoid climbing beyond tool root
		if dir == stopAt {
			return
		}
		dir = filepath.Dir(dir)
	}
}

func resolveTargetPath(targetBase, conventionPath, projectRoot string, isGlobal, hasScopeOverride bool) string {
	path := config.ExpandUserPath(conventionPath)
	if filepath.IsAbs(path) {
		return path
	}

	if hasScopeOverride {
		if isGlobal {
			home, _ := os.UserHomeDir()
			return filepath.Join(home, strings.TrimPrefix(path, "./"))
		}
		return filepath.Join(projectRoot, path)
	}

	return filepath.Join(targetBase, path)
}
