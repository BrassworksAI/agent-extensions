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
	Commands string `yaml:"commands"`
}

type ToolsConfig struct {
	Tools map[string]Tool `yaml:"tools"`
}

func LoadToolsConfig(path string) (*ToolsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tools config: %w", err)
	}

	var cfg ToolsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing tools config: %w", err)
	}

	return &cfg, nil
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
	path := t.GlobalPath
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}
	return path
}

func (t *Tool) ResolveLocalPath(projectRoot string) string {
	return filepath.Join(projectRoot, t.LocalPath)
}

func (c *Conventions) SkillPath(name string) string {
	return "skills/" + name + "/SKILL.md"
}

func (c *Conventions) CommandPath(name string) string {
	return strings.ReplaceAll(c.Commands, "{name}", name)
}
