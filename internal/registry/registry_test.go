package registry

import (
	"testing"
	"testing/fstest"
)

func createTestFS() fstest.MapFS {
	return fstest.MapFS{
		"tools.yaml": &fstest.MapFile{Data: []byte(`tools:
  tool-a:
    name: Tool A
    global_path: ~/.tool-a
    local_path: .tool-a
    conventions:
      skills: skills/{name}/SKILL.md
      commands: commands/{name}.md
  tool-b:
    name: Tool B
    global_path: ~/.tool-b
    local_path: .tool-b
    conventions:
      skills: skills/{name}.md
      commands: prompts/{name}.md
`)},
		"repository/commands/cmd-one.md":         &fstest.MapFile{Data: []byte("# Command One\nContent here")},
		"repository/commands/cmd-two.md":         &fstest.MapFile{Data: []byte("# Command Two\nContent here")},
		"repository/commands/cmd-three.md":       &fstest.MapFile{Data: []byte("# Command Three\nContent here")},
		"repository/skills/skill-one/SKILL.md":   &fstest.MapFile{Data: []byte("# Skill One\nContent")},
		"repository/skills/skill-two/SKILL.md":   &fstest.MapFile{Data: []byte("# Skill Two\nContent")},
		"repository/skills/skill-three/SKILL.md": &fstest.MapFile{Data: []byte("# Skill Three\nContent")},
	}
}

func TestNew(t *testing.T) {
	fsys := createTestFS()

	reg, err := New(fsys)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if reg.FS == nil {
		t.Error("FS should not be nil")
	}
	if reg.Tools == nil {
		t.Error("Tools should not be nil")
	}
}

func TestNew_MissingToolsYAML(t *testing.T) {
	fsys := fstest.MapFS{}

	_, err := New(fsys)
	if err == nil {
		t.Error("expected error for missing tools.yaml")
	}
}

func TestRegistry_GetToolNames(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	names := reg.GetToolNames()
	if len(names) != 2 {
		t.Errorf("expected 2 tools, got %d", len(names))
	}

	nameMap := make(map[string]bool)
	for _, n := range names {
		nameMap[n] = true
	}

	if !nameMap["tool-a"] {
		t.Error("tool-a should be in tool names")
	}
	if !nameMap["tool-b"] {
		t.Error("tool-b should be in tool names")
	}
}

func TestRegistry_GetAllCommands(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	commands := reg.GetAllCommands()
	if len(commands) != 3 {
		t.Errorf("expected 3 commands, got %d", len(commands))
	}

	cmdMap := make(map[string]bool)
	for _, c := range commands {
		cmdMap[c] = true
	}

	if !cmdMap["cmd-one"] {
		t.Error("cmd-one should be in commands")
	}
	if !cmdMap["cmd-two"] {
		t.Error("cmd-two should be in commands")
	}
	if !cmdMap["cmd-three"] {
		t.Error("cmd-three should be in commands")
	}
}

func TestRegistry_GetAllSkills(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	skills := reg.GetAllSkills()
	if len(skills) != 3 {
		t.Errorf("expected 3 skills, got %d", len(skills))
	}

	skillMap := make(map[string]bool)
	for _, s := range skills {
		skillMap[s] = true
	}

	if !skillMap["skill-one"] {
		t.Error("skill-one should be in skills")
	}
	if !skillMap["skill-two"] {
		t.Error("skill-two should be in skills")
	}
	if !skillMap["skill-three"] {
		t.Error("skill-three should be in skills")
	}
}

func TestRegistry_GetTool(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	tool, ok := reg.GetTool("tool-a")
	if !ok {
		t.Fatal("tool-a should exist")
	}

	if tool.Name != "Tool A" {
		t.Errorf("expected name 'Tool A', got %q", tool.Name)
	}
	if tool.GlobalPath != "~/.tool-a" {
		t.Errorf("expected global_path '~/.tool-a', got %q", tool.GlobalPath)
	}
}

func TestRegistry_GetTool_NotFound(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	_, ok := reg.GetTool("nonexistent")
	if ok {
		t.Error("nonexistent tool should return false")
	}
}

func TestRegistry_CommandSourcePath(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	path := reg.CommandSourcePath("my-command")
	expected := "repository/commands/my-command.md"
	if path != expected {
		t.Errorf("CommandSourcePath() = %q, want %q", path, expected)
	}
}

func TestRegistry_SkillSourcePath(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	path := reg.SkillSourcePath("my-skill")
	expected := "repository/skills/my-skill"
	if path != expected {
		t.Errorf("SkillSourcePath() = %q, want %q", path, expected)
	}
}

func TestRegistry_CommandExists(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	if !reg.CommandExists("cmd-one") {
		t.Error("cmd-one should exist")
	}
	if !reg.CommandExists("cmd-two") {
		t.Error("cmd-two should exist")
	}
	if reg.CommandExists("nonexistent") {
		t.Error("nonexistent command should not exist")
	}
}

func TestRegistry_SkillExists(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	if !reg.SkillExists("skill-one") {
		t.Error("skill-one should exist")
	}
	if !reg.SkillExists("skill-two") {
		t.Error("skill-two should exist")
	}
	if reg.SkillExists("nonexistent") {
		t.Error("nonexistent skill should not exist")
	}
}

func TestRegistry_ReadCommand(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	data, err := reg.ReadCommand("cmd-one")
	if err != nil {
		t.Fatalf("ReadCommand failed: %v", err)
	}

	expected := "# Command One\nContent here"
	if string(data) != expected {
		t.Errorf("ReadCommand() = %q, want %q", string(data), expected)
	}
}

func TestRegistry_ReadCommand_NotFound(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	_, err := reg.ReadCommand("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}

func TestRegistry_ListSkillFiles(t *testing.T) {
	fsys := createTestFS()
	reg, _ := New(fsys)

	entries, err := reg.ListSkillFiles("skill-one")
	if err != nil {
		t.Fatalf("ListSkillFiles failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 file, got %d", len(entries))
	}

	if entries[0].Name() != "SKILL.md" {
		t.Errorf("expected file 'SKILL.md', got %q", entries[0].Name())
	}
}
