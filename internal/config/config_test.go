package config

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestLoadToolsConfig(t *testing.T) {
	content := `tools:
  test-tool:
    name: Test Tool
    global_path: ~/.test-tool
    local_path: .test-tool
    conventions:
      commands: commands/{name}.md
  another-tool:
    name: Another Tool
    global_path: ~/.another
    local_path: .another
    conventions:
      commands: prompts/{name}.md
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tools.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadToolsConfig(configPath)
	if err != nil {
		t.Fatalf("LoadToolsConfig failed: %v", err)
	}

	if len(cfg.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(cfg.Tools))
	}

	tool, ok := cfg.Tools["test-tool"]
	if !ok {
		t.Fatal("test-tool not found")
	}

	if tool.Name != "Test Tool" {
		t.Errorf("expected name 'Test Tool', got %q", tool.Name)
	}
	if tool.GlobalPath != "~/.test-tool" {
		t.Errorf("expected global_path '~/.test-tool', got %q", tool.GlobalPath)
	}
	if tool.LocalPath != ".test-tool" {
		t.Errorf("expected local_path '.test-tool', got %q", tool.LocalPath)
	}
	if tool.Conventions.Commands != "commands/{name}.md" {
		t.Errorf("expected commands convention 'commands/{name}.md', got %q", tool.Conventions.Commands)
	}
}

func TestLoadToolsConfigFromFS(t *testing.T) {
	content := `tools:
  fs-tool:
    name: FS Tool
    global_path: ~/.fs-tool
    local_path: .fs-tool
    conventions:
      commands: commands/{name}.md
`
	fsys := fstest.MapFS{
		"tools.yaml": &fstest.MapFile{Data: []byte(content)},
	}

	cfg, err := LoadToolsConfigFromFS(fsys, "tools.yaml")
	if err != nil {
		t.Fatalf("LoadToolsConfigFromFS failed: %v", err)
	}

	if len(cfg.Tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(cfg.Tools))
	}

	tool, ok := cfg.Tools["fs-tool"]
	if !ok {
		t.Fatal("fs-tool not found")
	}
	if tool.Name != "FS Tool" {
		t.Errorf("expected name 'FS Tool', got %q", tool.Name)
	}
}

func TestLoadToolsConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tools.yaml")
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := LoadToolsConfig(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoadToolsConfig_FileNotFound(t *testing.T) {
	_, err := LoadToolsConfig("/nonexistent/path/tools.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestTool_ResolveGlobalPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name       string
		globalPath string
		want       string
	}{
		{
			name:       "tilde expansion",
			globalPath: "~/.test-tool",
			want:       filepath.Join(home, ".test-tool"),
		},
		{
			name:       "absolute path",
			globalPath: "/opt/test-tool",
			want:       "/opt/test-tool",
		},
		{
			name:       "nested tilde path",
			globalPath: "~/.config/test-tool",
			want:       filepath.Join(home, ".config/test-tool"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := Tool{GlobalPath: tt.globalPath}
			got := tool.ResolveGlobalPath()
			if got != tt.want {
				t.Errorf("ResolveGlobalPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTool_ResolveLocalPath(t *testing.T) {
	tests := []struct {
		name        string
		localPath   string
		projectRoot string
		want        string
	}{
		{
			name:        "simple local path",
			localPath:   ".test-tool",
			projectRoot: "/projects/myproject",
			want:        "/projects/myproject/.test-tool",
		},
		{
			name:        "nested local path",
			localPath:   ".config/test-tool",
			projectRoot: "/home/user/project",
			want:        "/home/user/project/.config/test-tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := Tool{LocalPath: tt.localPath}
			got := tool.ResolveLocalPath(tt.projectRoot)
			if got != tt.want {
				t.Errorf("ResolveLocalPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConventions_SkillPath(t *testing.T) {
	conv := Conventions{}
	got := conv.SkillPath("my-skill")
	want := "skills/my-skill/SKILL.md"
	if got != want {
		t.Errorf("SkillPath() = %q, want %q", got, want)
	}
}

func TestConventions_CommandPath(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		commandName string
		want        string
	}{
		{
			name:        "standard commands directory",
			pattern:     "commands/{name}.md",
			commandName: "my-command",
			want:        "commands/my-command.md",
		},
		{
			name:        "prompts directory",
			pattern:     "prompts/{name}.md",
			commandName: "my-command",
			want:        "prompts/my-command.md",
		},
		{
			name:        "workflows directory",
			pattern:     "workflows/{name}.md",
			commandName: "test",
			want:        "workflows/test.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := Conventions{Commands: tt.pattern}
			got := conv.CommandPath(tt.commandName)
			if got != tt.want {
				t.Errorf("CommandPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
