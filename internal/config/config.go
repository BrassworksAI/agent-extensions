package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Tool struct {
	Name        string      `yaml:"name"`
	GlobalPath  string      `yaml:"global_path"`
	LocalPath   string      `yaml:"local_path"`
	Conventions Conventions `yaml:"conventions"`
}

type Conventions struct {
	Commands       string `yaml:"commands"`
	GlobalCommands string `yaml:"global_commands,omitempty"`
	LocalCommands  string `yaml:"local_commands,omitempty"`
	Skills         string `yaml:"skills,omitempty"`
	GlobalSkills   string `yaml:"global_skills,omitempty"`
	LocalSkills    string `yaml:"local_skills,omitempty"`
}

type ToolsConfig struct {
	Tools map[string]Tool `yaml:"tools"`
}

func LoadToolsConfigFromFS(fsys fs.FS, path string) (*ToolsConfig, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("reading tools config: %w", err)
	}

	var cfg ToolsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing tools config: %w", err)
	}

	return &cfg, nil
}

func (t *Tool) ResolveGlobalPath() string {
	return ExpandUserPath(t.GlobalPath)
}

func (t *Tool) ResolveLocalPath(projectRoot string) string {
	return filepath.Join(projectRoot, t.LocalPath)
}

func (c *Conventions) SkillPath(name string) string {
	return resolveConventionPath(c.Skills, name, "skills/{name}/SKILL.md", filepath.Join(name, "SKILL.md"))
}

func (c *Conventions) ScopedSkillPath(name string, isGlobal bool) string {
	var pattern string
	if isGlobal && c.GlobalSkills != "" {
		pattern = c.GlobalSkills
	} else if !isGlobal && c.LocalSkills != "" {
		pattern = c.LocalSkills
	} else {
		pattern = c.SkillPath(name)
		return pattern
	}
	return resolveConventionPath(pattern, name, "skills/{name}/SKILL.md", filepath.Join(name, "SKILL.md"))
}

func (c *Conventions) CommandPath(name string, isGlobal bool) string {
	var pattern string
	if isGlobal && c.GlobalCommands != "" {
		pattern = c.GlobalCommands
	} else if !isGlobal && c.LocalCommands != "" {
		pattern = c.LocalCommands
	} else {
		pattern = c.Commands
	}
	return resolveConventionPath(pattern, name, "commands/{name}.md", name+".md")
}

func ExpandUserPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func resolveConventionPath(pattern, name, defaultPattern, inferredSuffix string) string {
	if pattern == "" {
		pattern = defaultPattern
	}

	if strings.Contains(pattern, "{name}") {
		return strings.ReplaceAll(pattern, "{name}", name)
	}

	if filepath.Ext(pattern) == ".md" {
		return pattern
	}

	return filepath.Join(pattern, inferredSuffix)
}
