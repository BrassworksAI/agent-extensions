package installer

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/shanepadgett/agent-extensions/internal/registry"
)

func createTestRegistry(t *testing.T, globalBase string) *registry.Registry {
	t.Helper()

	// Build custom YAML for this test
	toolsYAML := `tools:
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
  codex-style:
    name: CodexStyle
    global_path: ` + globalBase + `/.codex
    local_path: .codex
    conventions:
      skills: skills/{name}/SKILL.md
      global_skills: ~/.agents/skills
      local_skills: .agents/skills
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
	tools := []string{"dir-based", "prompts-style", "codex-style"}
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
			cmdPath := filepath.Join(targetBase, tool.Conventions.CommandPath(cmd, s == ScopeGlobal))
			verifySymlink(t, cmdPath, cmd+".md")
			verifySymlinkContent(t, cmdPath, "# Command")
		}

		// Verify skills
		for _, skill := range skills {
			skillPath := tool.Conventions.ScopedSkillPath(skill, s == ScopeGlobal)
			hasScopeOverride := (s == ScopeGlobal && tool.Conventions.GlobalSkills != "") || (s != ScopeGlobal && tool.Conventions.LocalSkills != "")
			fullPath := resolveTargetPath(targetBase, skillPath, projectRoot, s == ScopeGlobal, hasScopeOverride)

			// All skills are directory-based - check the skill directory is symlinked
			skillDir := filepath.Dir(fullPath)
			verifySymlinkedSkillDir(t, skillDir)
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

func verifySymlinkedSkillDir(t *testing.T, path string) {
	t.Helper()

	info, err := os.Lstat(path)
	if err != nil {
		t.Errorf("skill symlink not found at %s: %v", path, err)
		return
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("%s is not a symlink", path)
		return
	}

	resolved, err := os.Stat(path)
	if err != nil {
		t.Errorf("broken skill symlink at %s: %v", path, err)
		return
	}
	if !resolved.IsDir() {
		t.Errorf("symlink target for %s is not a directory", path)
		return
	}

	// Verify linked directory contains SKILL.md file.
	skillFile := filepath.Join(path, "SKILL.md")
	fileInfo, err := os.Stat(skillFile)
	if err != nil {
		t.Errorf("SKILL.md not found in %s: %v", path, err)
		return
	}

	if fileInfo.IsDir() {
		t.Errorf("SKILL.md in %s is unexpectedly a directory", path)
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
	tools := []string{"dir-based", "prompts-style", "codex-style"}
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
			cmdPath := filepath.Join(targetBase, tool.Conventions.CommandPath(cmd, s == ScopeGlobal))
			if _, err := os.Lstat(cmdPath); err == nil {
				t.Errorf("command still exists at %s", cmdPath)
			}
		}

		// Verify skills are removed
		for _, skill := range skills {
			skillPath := tool.Conventions.ScopedSkillPath(skill, s == ScopeGlobal)
			hasScopeOverride := (s == ScopeGlobal && tool.Conventions.GlobalSkills != "") || (s != ScopeGlobal && tool.Conventions.LocalSkills != "")
			fullPath := resolveTargetPath(targetBase, skillPath, projectRoot, s == ScopeGlobal, hasScopeOverride)

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
	if _, statErr := os.Stat(commandsDir); statErr != nil {
		t.Errorf("commands directory should exist: %v", statErr)
	}
	skillsDir := filepath.Join(targetBase, "skills")
	if _, statErr := os.Stat(skillsDir); statErr != nil {
		t.Errorf("skills directory should exist: %v", statErr)
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

	// For dir-based tool, skills should be at skills/{name} as a symlink to cached skill directory
	skillDir := filepath.Join(targetBase, "skills", "skill-one")
	info, err := os.Lstat(skillDir)
	if err != nil {
		t.Fatalf("skill symlink not found: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("skill should be a symlink")
	}

	// The linked directory should contain SKILL.md.
	skillFile := filepath.Join(skillDir, "SKILL.md")
	fileInfo, err := os.Stat(skillFile)
	if err != nil {
		t.Errorf("SKILL.md not found in linked skill directory: %v", err)
	}

	if fileInfo.IsDir() {
		t.Error("SKILL.md should be a file")
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

func TestToolConventions_CodexSkillPaths(t *testing.T) {
	globalDir := t.TempDir()
	projectRoot := t.TempDir()
	reg := createTestRegistry(t, globalDir)
	inst := New(reg, projectRoot)

	_, err := inst.Install("codex-style", ScopeBoth)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	tool, _ := reg.GetTool("codex-style")
	globalSkill := resolveTargetPath(
		tool.ResolveGlobalPath(),
		tool.Conventions.ScopedSkillPath("skill-one", true),
		projectRoot,
		true,
		tool.Conventions.GlobalSkills != "",
	)
	localSkill := resolveTargetPath(
		tool.ResolveLocalPath(projectRoot),
		tool.Conventions.ScopedSkillPath("skill-one", false),
		projectRoot,
		false,
		tool.Conventions.LocalSkills != "",
	)

	globalSkillDir := filepath.Dir(globalSkill)
	localSkillDir := filepath.Dir(localSkill)

	if _, err := os.Stat(globalSkillDir); err != nil {
		t.Fatalf("global codex skill directory not found: %v", err)
	}
	if _, err := os.Stat(localSkillDir); err != nil {
		t.Fatalf("local codex skill directory not found: %v", err)
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
		_, _ = inst.Install("bench-tool", ScopeGlobal)
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
		_, _ = inst.Install("bench-tool", ScopeGlobal)

		b.ResetTimer()
		_, _ = inst.Uninstall("bench-tool", ScopeGlobal)
	}
}
