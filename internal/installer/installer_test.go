package installer

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/shanepadgett/agent-extensions/internal/registry"
)

// Test tools with different conventions
var testToolsYAML = `tools:
  dir-based:
    name: DirBased
    global_path: {GLOBAL}/dir-based
    local_path: .dir-based
    conventions:
      skills: skills/{name}/SKILL.md
      commands: commands/{name}.md
  prompts-style:
    name: PromptsStyle
    global_path: {GLOBAL}/prompts-style
    local_path: .prompts-style
    conventions:
      skills: skills/{name}/SKILL.md
      commands: prompts/{name}.md
`

func createTestRegistry(t *testing.T, globalBase string) *registry.Registry {
	t.Helper()

	// Replace placeholder with actual temp directory
	toolsYAML := testToolsYAML
	toolsYAML = filepath.ToSlash(toolsYAML)
	// We need absolute path, not tilde
	toolsYAML = `tools:
  dir-based:
    name: DirBased
    global_path: ` + globalBase + `/dir-based
    local_path: .dir-based
    conventions:
      skills: skills/{name}/SKILL.md
      commands: commands/{name}.md
  prompts-style:
    name: PromptsStyle
    global_path: ` + globalBase + `/prompts-style
    local_path: .prompts-style
    conventions:
      skills: skills/{name}/SKILL.md
      commands: prompts/{name}.md
`

	fsys := fstest.MapFS{
		"tools.yaml":                            &fstest.MapFile{Data: []byte(toolsYAML)},
		"repository/commands/cmd-one.md":        &fstest.MapFile{Data: []byte("# Command One\nThis is command one content.")},
		"repository/commands/cmd-two.md":        &fstest.MapFile{Data: []byte("# Command Two\nThis is command two content.")},
		"repository/commands/cmd-three.md":      &fstest.MapFile{Data: []byte("# Command Three\nThis is command three content.")},
		"repository/skills/skill-one/SKILL.md":  &fstest.MapFile{Data: []byte("# Skill One\nSkill one content.")},
		"repository/skills/skill-one/helper.md": &fstest.MapFile{Data: []byte("Helper file for skill one.")},
		"repository/skills/skill-two/SKILL.md":  &fstest.MapFile{Data: []byte("# Skill Two\nSkill two content.")},
	}

	reg, err := registry.New(fsys)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	return reg
}

func TestInstaller_New(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)

	inst := New(reg, projectRoot)
	if inst == nil {
		t.Fatal("New() returned nil")
	}
	if inst.Registry != reg {
		t.Error("Registry not set correctly")
	}
	if inst.ProjectRoot != projectRoot {
		t.Error("ProjectRoot not set correctly")
	}
}

func TestInstaller_cacheDir(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	home, _ := os.UserHomeDir()

	globalCache := inst.cacheDir(ScopeGlobal)
	expectedGlobal := filepath.Join(home, ".agents", "ae")
	if globalCache != expectedGlobal {
		t.Errorf("global cache = %q, want %q", globalCache, expectedGlobal)
	}

	localCache := inst.cacheDir(ScopeLocal)
	expectedLocal := filepath.Join(projectRoot, ".agents", "ae")
	if localCache != expectedLocal {
		t.Errorf("local cache = %q, want %q", localCache, expectedLocal)
	}
}

// TestInstallMatrix tests all tool × scope combinations
func TestInstallMatrix(t *testing.T) {
	tools := []string{"dir-based", "prompts-style"}
	scopes := []Scope{ScopeGlobal, ScopeLocal, ScopeBoth}

	for _, tool := range tools {
		for _, scope := range scopes {
			name := tool + "/" + string(scope)
			t.Run(name, func(t *testing.T) {
				globalDir := t.TempDir()
				projectRoot := t.TempDir()
				reg := createTestRegistry(t, globalDir)
				inst := New(reg, projectRoot)

				result, err := inst.Install(tool, scope)
				if err != nil {
					t.Fatalf("Install failed: %v", err)
				}

				if len(result.Errors) > 0 {
					for _, e := range result.Errors {
						t.Errorf("Install error: %v", e)
					}
				}

				// Verify installation based on scope
				verifyInstallation(t, reg, inst, tool, scope, globalDir, projectRoot)
			})
		}
	}
}

func verifyInstallation(t *testing.T, reg *registry.Registry, inst *Installer, toolName string, scope Scope, globalDir, projectRoot string) {
	t.Helper()

	tool, _ := reg.GetTool(toolName)
	commands := reg.GetAllCommands()
	skills := reg.GetAllSkills()

	checkScopes := []Scope{scope}
	if scope == ScopeBoth {
		checkScopes = []Scope{ScopeGlobal, ScopeLocal}
	}

	for _, s := range checkScopes {
		var targetBase string
		if s == ScopeGlobal {
			targetBase = tool.ResolveGlobalPath()
		} else {
			targetBase = tool.ResolveLocalPath(projectRoot)
		}

		// Verify commands
		for _, cmd := range commands {
			cmdPath := filepath.Join(targetBase, tool.Conventions.CommandPath(cmd))
			verifySymlink(t, cmdPath, cmd+".md")
			verifySymlinkContent(t, cmdPath, "# Command")
		}

		// Verify skills
		for _, skill := range skills {
			skillPath := tool.Conventions.SkillPath(skill)
			fullPath := filepath.Join(targetBase, skillPath)

			// All skills are directory-based - check the directory contains symlinked files
			skillDir := filepath.Dir(fullPath)
			verifyDirWithSymlinkedFiles(t, skillDir)
		}
	}
}

func verifySymlink(t *testing.T, path, expectedTarget string) {
	t.Helper()

	info, err := os.Lstat(path)
	if err != nil {
		t.Errorf("symlink not found at %s: %v", path, err)
		return
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("%s is not a symlink", path)
		return
	}

	// Verify it resolves (not broken)
	_, err = os.Stat(path)
	if err != nil {
		t.Errorf("broken symlink at %s: %v", path, err)
	}
}

func verifySymlinkContent(t *testing.T, path, expectedContains string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("failed to read symlink target at %s: %v", path, err)
		return
	}

	if len(data) == 0 {
		t.Errorf("symlink target at %s is empty", path)
		return
	}

	content := string(data)
	if expectedContains != "" && !contains(content, expectedContains) {
		t.Errorf("content at %s does not contain %q", path, expectedContains)
	}
}

func verifySymlinkIsDir(t *testing.T, path string) {
	t.Helper()

	info, err := os.Lstat(path)
	if err != nil {
		t.Errorf("path not found at %s: %v", path, err)
		return
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("%s is not a symlink", path)
		return
	}

	// Verify it resolves to a directory
	targetInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("broken symlink at %s: %v", path, err)
		return
	}

	if !targetInfo.IsDir() {
		t.Errorf("%s does not resolve to a directory", path)
	}
}

func verifyDirWithSymlinkedFiles(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("directory not found at %s: %v", path, err)
		return
	}

	if !info.IsDir() {
		t.Errorf("%s is not a directory", path)
		return
	}

	// Verify it contains at least one symlinked file (SKILL.md)
	skillFile := filepath.Join(path, "SKILL.md")
	fileInfo, err := os.Lstat(skillFile)
	if err != nil {
		t.Errorf("SKILL.md not found in %s: %v", path, err)
		return
	}

	if fileInfo.Mode()&os.ModeSymlink == 0 {
		t.Errorf("SKILL.md in %s is not a symlink", path)
		return
	}

	// Verify symlink resolves
	_, err = os.Stat(skillFile)
	if err != nil {
		t.Errorf("broken symlink at %s: %v", skillFile, err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestUninstallMatrix tests all tool × scope uninstall combinations
func TestUninstallMatrix(t *testing.T) {
	tools := []string{"dir-based", "prompts-style"}
	scopes := []Scope{ScopeGlobal, ScopeLocal, ScopeBoth}

	for _, tool := range tools {
		for _, scope := range scopes {
			name := tool + "/" + string(scope)
			t.Run(name, func(t *testing.T) {
				globalDir := t.TempDir()
				projectRoot := t.TempDir()
				reg := createTestRegistry(t, globalDir)
				inst := New(reg, projectRoot)

				// First install
				_, err := inst.Install(tool, scope)
				if err != nil {
					t.Fatalf("Install failed: %v", err)
				}

				// Then uninstall
				result, err := inst.Uninstall(tool, scope)
				if err != nil {
					t.Fatalf("Uninstall failed: %v", err)
				}

				// Verify counts
				expectedCommands := len(reg.GetAllCommands())
				expectedSkills := len(reg.GetAllSkills())

				if scope == ScopeBoth {
					expectedCommands *= 2
					expectedSkills *= 2
				}

				if result.Commands != expectedCommands {
					t.Errorf("uninstalled %d commands, expected %d", result.Commands, expectedCommands)
				}
				if result.Skills != expectedSkills {
					t.Errorf("uninstalled %d skills, expected %d", result.Skills, expectedSkills)
				}

				// Verify removal
				verifyUninstallation(t, reg, tool, scope, globalDir, projectRoot)
			})
		}
	}
}

func verifyUninstallation(t *testing.T, reg *registry.Registry, toolName string, scope Scope, globalDir, projectRoot string) {
	t.Helper()

	tool, _ := reg.GetTool(toolName)
	commands := reg.GetAllCommands()
	skills := reg.GetAllSkills()

	checkScopes := []Scope{scope}
	if scope == ScopeBoth {
		checkScopes = []Scope{ScopeGlobal, ScopeLocal}
	}

	for _, s := range checkScopes {
		var targetBase string
		if s == ScopeGlobal {
			targetBase = tool.ResolveGlobalPath()
		} else {
			targetBase = tool.ResolveLocalPath(projectRoot)
		}

		// Verify commands are removed
		for _, cmd := range commands {
			cmdPath := filepath.Join(targetBase, tool.Conventions.CommandPath(cmd))
			if _, err := os.Lstat(cmdPath); err == nil {
				t.Errorf("command still exists at %s", cmdPath)
			}
		}

		// Verify skills are removed
		for _, skill := range skills {
			skillPath := tool.Conventions.SkillPath(skill)
			fullPath := filepath.Join(targetBase, skillPath)

			// All skills are directory-based
			skillDir := filepath.Dir(fullPath)
			if _, err := os.Lstat(skillDir); err == nil {
				t.Errorf("skill directory still exists at %s", skillDir)
			}
		}

		// Verify tool root directory is removed when empty
		if _, err := os.Stat(targetBase); err == nil {
			entries, _ := os.ReadDir(targetBase)
			if len(entries) == 0 {
				t.Errorf("empty tool root directory should be removed: %s", targetBase)
			}
		}
	}
}

func TestUninstall_NotInstalled(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	// Install to only one tool
	_, err := inst.Install("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Uninstall from a tool that was never installed to
	result, err := inst.Uninstall("prompts-style", ScopeGlobal)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Should report zero removals since nothing was installed there
	if result.Commands != 0 {
		t.Errorf("expected 0 commands removed, got %d", result.Commands)
	}
	if result.Skills != 0 {
		t.Errorf("expected 0 skills removed, got %d", result.Skills)
	}

	// Uninstall from the tool that was actually installed to
	result, err = inst.Uninstall("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Should report actual removals
	expectedCommands := len(reg.GetAllCommands())
	expectedSkills := len(reg.GetAllSkills())

	if result.Commands != expectedCommands {
		t.Errorf("expected %d commands removed, got %d", expectedCommands, result.Commands)
	}
	if result.Skills != expectedSkills {
		t.Errorf("expected %d skills removed, got %d", expectedSkills, result.Skills)
	}
}

func TestInstall_UnknownTool(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	_, err := inst.Install("nonexistent-tool", ScopeGlobal)
	if err == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestUninstall_UnknownTool(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	_, err := inst.Uninstall("nonexistent-tool", ScopeGlobal)
	if err == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestInstall_Reinstall(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	// Install once
	result1, err := inst.Install("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}

	// Install again (should overwrite without error)
	result2, err := inst.Install("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("second install failed: %v", err)
	}

	if result1.Commands != result2.Commands {
		t.Errorf("command counts differ: %d vs %d", result1.Commands, result2.Commands)
	}
	if result1.Skills != result2.Skills {
		t.Errorf("skill counts differ: %d vs %d", result1.Skills, result2.Skills)
	}

	// Verify still works
	verifyInstallation(t, reg, inst, "dir-based", ScopeGlobal, globalDir, projectRoot)
}

func TestCacheContentIntegrity(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	_, err := inst.Install("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Verify cache files have correct content
	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".agents", "ae")
	readmePath := filepath.Join(cacheDir, "README.md")
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read cache readme: %v", err)
	}
	if !contains(string(readme), "Agent Extensions Cache") {
		t.Errorf("cache readme content mismatch: %q", string(readme))
	}

	// Check command cache
	cmdCache := filepath.Join(cacheDir, "commands", "cmd-one.md")
	data, err := os.ReadFile(cmdCache)
	if err != nil {
		t.Fatalf("failed to read cached command: %v", err)
	}
	if string(data) != "# Command One\nThis is command one content." {
		t.Errorf("cached command content mismatch: %q", string(data))
	}

	// Check skill cache
	skillCache := filepath.Join(cacheDir, "skills", "skill-one", "SKILL.md")
	data, err = os.ReadFile(skillCache)
	if err != nil {
		t.Fatalf("failed to read cached skill: %v", err)
	}
	if string(data) != "# Skill One\nSkill one content." {
		t.Errorf("cached skill content mismatch: %q", string(data))
	}

	// Verify helper file is also copied
	helperCache := filepath.Join(cacheDir, "skills", "skill-one", "helper.md")
	data, err = os.ReadFile(helperCache)
	if err != nil {
		t.Fatalf("failed to read cached helper: %v", err)
	}
	if string(data) != "Helper file for skill one." {
		t.Errorf("cached helper content mismatch: %q", string(data))
	}
}

func TestEmptyDirectoryCleanup(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	tool, _ := reg.GetTool("dir-based")
	targetBase := tool.ResolveGlobalPath()

	// Install
	_, err := inst.Install("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Verify directories exist
	commandsDir := filepath.Join(targetBase, "commands")
	if _, err := os.Stat(commandsDir); err != nil {
		t.Errorf("commands directory should exist: %v", err)
	}
	skillsDir := filepath.Join(targetBase, "skills")
	if _, err := os.Stat(skillsDir); err != nil {
		t.Errorf("skills directory should exist: %v", err)
	}

	// Uninstall
	_, err = inst.Uninstall("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Verify empty directories are cleaned up
	if _, err := os.Stat(commandsDir); err == nil {
		entries, _ := os.ReadDir(commandsDir)
		if len(entries) == 0 {
			t.Error("empty commands directory should be removed")
		}
	}

	if _, err := os.Stat(skillsDir); err == nil {
		entries, _ := os.ReadDir(skillsDir)
		if len(entries) == 0 {
			t.Error("empty skills directory should be removed")
		}
	}

	// Verify tool root directory is also removed when empty
	if _, err := os.Stat(targetBase); err == nil {
		entries, _ := os.ReadDir(targetBase)
		if len(entries) == 0 {
			t.Error("empty tool root directory should be removed")
		}
	}
}

// Test specific tool conventions
func TestToolConventions_DirBasedSkills(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	_, err := inst.Install("dir-based", ScopeGlobal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	tool, _ := reg.GetTool("dir-based")
	targetBase := tool.ResolveGlobalPath()

	// For dir-based tool, skills should be at skills/{name} (directory with symlinked files)
	skillDir := filepath.Join(targetBase, "skills", "skill-one")
	info, err := os.Stat(skillDir)
	if err != nil {
		t.Fatalf("skill directory not found: %v", err)
	}

	if !info.IsDir() {
		t.Error("skill should be a directory")
	}

	// The directory should contain SKILL.md as a symlink
	skillFile := filepath.Join(skillDir, "SKILL.md")
	fileInfo, err := os.Lstat(skillFile)
	if err != nil {
		t.Errorf("SKILL.md not found in skill directory: %v", err)
	}

	if fileInfo.Mode()&os.ModeSymlink == 0 {
		t.Error("SKILL.md should be a symlink")
	}
}

func TestToolConventions_PromptsStyle(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	_, err := inst.Install("prompts-style", ScopeGlobal)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	tool, _ := reg.GetTool("prompts-style")
	targetBase := tool.ResolveGlobalPath()

	// For prompts-style tool, commands should be at prompts/{name}.md
	cmdFile := filepath.Join(targetBase, "prompts", "cmd-one.md")
	if _, err := os.Stat(cmdFile); err != nil {
		t.Errorf("command file not found at prompts directory: %v", err)
	}
}

// Benchmark tests
func BenchmarkInstall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		globalDir := b.TempDir()
		projectRoot := b.TempDir()

		fsys := fstest.MapFS{
			"tools.yaml": &fstest.MapFile{Data: []byte(`tools:
  bench-tool:
    name: Bench Tool
    global_path: ` + globalDir + `/bench-tool
    local_path: .bench-tool
    conventions:
      skills: skills/{name}/SKILL.md
      commands: commands/{name}.md
`)},
			"repository/commands/cmd-one.md":       &fstest.MapFile{Data: []byte("# Command")},
			"repository/skills/skill-one/SKILL.md": &fstest.MapFile{Data: []byte("# Skill")},
		}

		reg, _ := registry.New(fsys)
		inst := New(reg, projectRoot)
		inst.Install("bench-tool", ScopeGlobal)
	}
}

func BenchmarkUninstall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		globalDir := b.TempDir()
		projectRoot := b.TempDir()

		fsys := fstest.MapFS{
			"tools.yaml": &fstest.MapFile{Data: []byte(`tools:
  bench-tool:
    name: Bench Tool
    global_path: ` + globalDir + `/bench-tool
    local_path: .bench-tool
    conventions:
      skills: skills/{name}/SKILL.md
      commands: commands/{name}.md
`)},
			"repository/commands/cmd-one.md":       &fstest.MapFile{Data: []byte("# Command")},
			"repository/skills/skill-one/SKILL.md": &fstest.MapFile{Data: []byte("# Skill")},
		}

		reg, _ := registry.New(fsys)
		inst := New(reg, projectRoot)
		inst.Install("bench-tool", ScopeGlobal)

		b.ResetTimer()
		inst.Uninstall("bench-tool", ScopeGlobal)
	}
}
