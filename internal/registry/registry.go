package registry

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/shanepadgett/agent-extensions/internal/config"
)

type Registry struct {
	FS    fs.FS
	Tools *config.ToolsConfig
}

func New(fsys fs.FS) (*Registry, error) {
	tools, err := config.LoadToolsConfigFromFS(fsys, "tools.yaml")
	if err != nil {
		return nil, err
	}

	return &Registry{
		FS:    fsys,
		Tools: tools,
	}, nil
}

func (r *Registry) GetToolNames() []string {
	names := make([]string, 0, len(r.Tools.Tools))
	for name := range r.Tools.Tools {
		names = append(names, name)
	}
	return names
}

func (r *Registry) GetTool(name string) (config.Tool, bool) {
	tool, ok := r.Tools.Tools[name]
	return tool, ok
}

func (r *Registry) GetAllCommands() []string {
	var commands []string
	entries, err := fs.ReadDir(r.FS, "repository/commands")
	if err != nil {
		return commands
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			commands = append(commands, name)
		}
	}
	return commands
}

func (r *Registry) GetAllSkills() []string {
	var skills []string
	entries, err := fs.ReadDir(r.FS, "repository/skills")
	if err != nil {
		return skills
	}
	for _, entry := range entries {
		if entry.IsDir() {
			skills = append(skills, entry.Name())
		}
	}
	return skills
}

func (r *Registry) CommandSourcePath(commandName string) string {
	return filepath.Join("repository", "commands", commandName+".md")
}

func (r *Registry) SkillSourcePath(skillName string) string {
	return filepath.Join("repository", "skills", skillName)
}

func (r *Registry) CommandExists(name string) bool {
	_, err := fs.Stat(r.FS, r.CommandSourcePath(name))
	return err == nil
}

func (r *Registry) SkillExists(name string) bool {
	_, err := fs.Stat(r.FS, r.SkillSourcePath(name))
	return err == nil
}

func (r *Registry) ReadCommand(name string) ([]byte, error) {
	return fs.ReadFile(r.FS, r.CommandSourcePath(name))
}

func (r *Registry) ListSkillFiles(skillName string) ([]fs.DirEntry, error) {
	return fs.ReadDir(r.FS, r.SkillSourcePath(skillName))
}
