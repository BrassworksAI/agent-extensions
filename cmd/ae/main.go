package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shanepadgett/agent-extensions/internal/embedded"
	"github.com/shanepadgett/agent-extensions/internal/installer"
	"github.com/shanepadgett/agent-extensions/internal/registry"
	"github.com/shanepadgett/agent-extensions/internal/ui"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ae",
	Short: "Agent Extensions CLI",
	Long:  "Manage installation of commands and skills for AI coding agents",
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install extensions to agent tools",
	RunE:  runInstall,
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall extensions from agent tools",
	RunE:  runUninstall,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available extensions and tools",
	RunE:  runList,
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check installation health and diagnose issues",
	RunE:  runDoctor,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update repository and refresh installed extensions",
	RunE:  runUpdate,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ae version %s\n", version)
	},
}

var (
	flagTools []string
	flagScope string
	flagYes   bool
)

func init() {
	installCmd.Flags().StringSliceVarP(&flagTools, "tools", "t", nil, "Tools to install to (comma-separated)")
	installCmd.Flags().StringVarP(&flagScope, "scope", "s", "", "Installation scope (global, local, both)")
	installCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation")

	uninstallCmd.Flags().StringSliceVarP(&flagTools, "tools", "t", nil, "Tools to uninstall from (comma-separated)")
	uninstallCmd.Flags().StringVarP(&flagScope, "scope", "s", "", "Uninstallation scope (global, local, both)")
	uninstallCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
}

func getRegistry() (*registry.Registry, error) {
	fsys, err := embedded.FS()
	if err != nil {
		return nil, fmt.Errorf("loading embedded content: %w", err)
	}
	return registry.New(fsys)
}

func getProjectRoot() string {
	cwd, _ := os.Getwd()
	return cwd
}

func runInstall(cmd *cobra.Command, args []string) error {
	u := ui.New()

	u.Title()

	reg, err := getRegistry()
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	var selectedTools []string
	var scope string

	totalCommands := len(reg.GetAllCommands())
	totalSkills := len(reg.GetAllSkills())

	if len(flagTools) > 0 {
		selectedTools = flagTools
	} else {
		tools := reg.GetToolNames()
		sort.Strings(tools)
		selectedTools, err = u.ChooseMulti("Select tools to install to:", tools)
		if err != nil {
			return err
		}
	}

	if flagScope != "" {
		scope = flagScope
	} else {
		scope, err = u.Choose("Select scope:", []string{"global", "local", "both"})
		if err != nil {
			return err
		}
	}

	var installScope installer.Scope
	switch scope {
	case "global":
		installScope = installer.ScopeGlobal
	case "local":
		installScope = installer.ScopeLocal
	case "both":
		installScope = installer.ScopeBoth
	default:
		return fmt.Errorf("unknown scope: %s", scope)
	}

	// Confirm
	if !flagYes {
		confirmMsg := fmt.Sprintf("Install %d commands and %d skills to %d tools (%s)?",
			totalCommands, totalSkills, len(selectedTools), scope)
		confirmed, err := u.Confirm(confirmMsg)
		if err != nil {
			return err
		}
		if !confirmed {
			u.Warn("Installation cancelled")
			return nil
		}
	}

	// Install
	inst := installer.New(reg, getProjectRoot())
	var lines []string

	err = u.Spin("Installing extensions", func() error {
		for _, toolName := range selectedTools {
			result, err := inst.Install(toolName, installScope)
			if err != nil {
				lines = append(lines, fmt.Sprintf("✗ %s: %v", toolName, err))
				continue
			}

			if len(result.Errors) > 0 {
				for _, e := range result.Errors {
					lines = append(lines, fmt.Sprintf("! %s: %v", toolName, e))
				}
			}

			tool, _ := reg.GetTool(toolName)
			lines = append(lines, fmt.Sprintf("✓ %s: %d commands, %d skills", tool.Name, result.Commands, result.Skills))
		}
		return nil
	})
	if err != nil {
		return err
	}

	u.Summary(strings.Join(lines, "\n"))

	return nil
}

func runUninstall(cmd *cobra.Command, args []string) error {
	u := ui.New()

	u.Title()

	reg, err := getRegistry()
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	var selectedTools []string
	var scope string

	totalCommands := len(reg.GetAllCommands())
	totalSkills := len(reg.GetAllSkills())

	if len(flagTools) > 0 {
		selectedTools = flagTools
	} else {
		tools := reg.GetToolNames()
		sort.Strings(tools)
		selectedTools, err = u.ChooseMulti("Select tools to uninstall from:", tools)
		if err != nil {
			return err
		}
	}

	if flagScope != "" {
		scope = flagScope
	} else {
		scope, err = u.Choose("Select scope:", []string{"global", "local", "both"})
		if err != nil {
			return err
		}
	}

	var uninstallScope installer.Scope
	switch scope {
	case "global":
		uninstallScope = installer.ScopeGlobal
	case "local":
		uninstallScope = installer.ScopeLocal
	case "both":
		uninstallScope = installer.ScopeBoth
	default:
		return fmt.Errorf("unknown scope: %s", scope)
	}

	// Confirm
	if !flagYes {
		confirmed, err := u.Confirm(fmt.Sprintf("Uninstall %d commands and %d skills from %d tools (%s)?",
			totalCommands, totalSkills, len(selectedTools), scope))
		if err != nil {
			return err
		}
		if !confirmed {
			u.Warn("Uninstallation cancelled")
			return nil
		}
	}

	// Uninstall
	inst := installer.New(reg, getProjectRoot())
	var lines []string

	err = u.Spin("Removing extensions", func() error {
		for _, toolName := range selectedTools {
			result, err := inst.Uninstall(toolName, uninstallScope)
			if err != nil {
				lines = append(lines, fmt.Sprintf("✗ %s: %v", toolName, err))
				continue
			}

			if result.Commands > 0 || result.Skills > 0 {
				tool, _ := reg.GetTool(toolName)
				lines = append(lines, fmt.Sprintf("✓ %s: removed %d commands, %d skills", tool.Name, result.Commands, result.Skills))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	u.Summary(strings.Join(lines, "\n"))

	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	u := ui.New()

	reg, err := getRegistry()
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	projectRoot := getProjectRoot()
	tools := reg.GetToolNames()
	sort.Strings(tools)
	commands := reg.GetAllCommands()
	skills := reg.GetAllSkills()

	u.Header("\nAvailable Content")
	fmt.Printf("  Commands: %d\n", len(commands))
	fmt.Printf("  Skills:   %d\n", len(skills))

	// Check installation status per tool
	type installStatus struct {
		global bool
		local  bool
	}

	u.Header("\nInstallation Status")
	fmt.Println("G=global  L=local  GL=both")
	fmt.Println()

	for _, toolKey := range tools {
		tool, _ := reg.GetTool(toolKey)
		globalPath := tool.ResolveGlobalPath()
		localPath := tool.ResolveLocalPath(projectRoot)

		status := installStatus{}

		// Check if any command is installed
		for _, c := range commands {
			if _, err := os.Stat(filepath.Join(globalPath, tool.Conventions.CommandPath(c, true))); err == nil {
				status.global = true
			}
			if _, err := os.Stat(filepath.Join(localPath, tool.Conventions.CommandPath(c, false))); err == nil {
				status.local = true
			}
		}

		// Check if any skill is installed
		for _, s := range skills {
			skillPath := tool.Conventions.SkillPath(s)
			globalSkillPath := filepath.Join(globalPath, skillPath)
			localSkillPath := filepath.Join(localPath, skillPath)

			if filepath.Ext(skillPath) != ".md" {
				globalSkillPath = filepath.Dir(globalSkillPath)
				localSkillPath = filepath.Dir(localSkillPath)
			}

			if _, err := os.Stat(globalSkillPath); err == nil {
				status.global = true
			}
			if _, err := os.Stat(localSkillPath); err == nil {
				status.local = true
			}
		}

		statusStr := "  "
		if status.global && status.local {
			statusStr = "GL"
		} else if status.global {
			statusStr = "G "
		} else if status.local {
			statusStr = " L"
		}

		fmt.Printf("  [%s] %s\n", statusStr, tool.Name)
	}

	fmt.Println()
	return nil
}

func runDoctor(cmd *cobra.Command, args []string) error {
	u := ui.New()

	u.Header("\n  Agent Extensions Doctor\n")

	reg, err := getRegistry()
	if err != nil {
		u.Error(fmt.Sprintf("Config: %v", err))
		return nil
	}
	u.Success("Config: tools.yaml loaded (embedded)")

	// Check each tool's global path
	u.Header("\nTool Paths:")
	tools := reg.GetToolNames()
	sort.Strings(tools)

	for _, name := range tools {
		tool, _ := reg.GetTool(name)
		globalPath := tool.ResolveGlobalPath()

		if _, err := os.Stat(globalPath); err == nil {
			u.Success(fmt.Sprintf("%s: %s exists", tool.Name, tool.GlobalPath))
		} else {
			u.Warn(fmt.Sprintf("%s: %s not found (tool may not be installed)", tool.Name, tool.GlobalPath))
		}
	}

	// Check cache directories
	u.Header("\nCache:")
	home, _ := os.UserHomeDir()
	globalCache := filepath.Join(home, ".agents", "ae")
	localCache := filepath.Join(getProjectRoot(), ".agents", "ae")

	if _, err := os.Stat(globalCache); err == nil {
		u.Success(fmt.Sprintf("Global cache: %s", globalCache))
	} else {
		u.Info(fmt.Sprintf("Global cache: not created yet (%s)", globalCache))
	}

	if _, err := os.Stat(localCache); err == nil {
		u.Success(fmt.Sprintf("Local cache: %s", localCache))
	} else {
		u.Info(fmt.Sprintf("Local cache: not created yet (%s)", localCache))
	}

	// Check for broken symlinks
	u.Header("\nSymlink Health:")
	brokenLinks := 0

	for _, name := range tools {
		tool, _ := reg.GetTool(name)
		globalPath := tool.ResolveGlobalPath()

		// Check commands dir
		cmdDir := filepath.Join(globalPath, filepath.Dir(tool.Conventions.CommandPath("test", true)))
		if entries, err := os.ReadDir(cmdDir); err == nil {
			for _, entry := range entries {
				fullPath := filepath.Join(cmdDir, entry.Name())
				if info, err := os.Lstat(fullPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
					if _, err := os.Stat(fullPath); err != nil {
						u.Warn(fmt.Sprintf("Broken symlink: %s", fullPath))
						brokenLinks++
					}
				}
			}
		}
	}

	if brokenLinks == 0 {
		u.Success("No broken symlinks found")
	} else {
		u.Warn(fmt.Sprintf("Found %d broken symlinks (run 'ae install' to fix)", brokenLinks))
	}

	fmt.Println()
	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	u := ui.New()

	u.Title()

	u.Info(fmt.Sprintf("ae version %s", version))
	u.Info("Extensions are embedded in the binary. To get new extensions, update the ae binary itself.")
	fmt.Println()

	reg, err := getRegistry()
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	projectRoot := getProjectRoot()
	tools := reg.GetToolNames()
	commands := reg.GetAllCommands()

	// Find what's currently installed and refresh symlinks
	inst := installer.New(reg, projectRoot)

	fmt.Println()
	u.Info("Refreshing installed extensions...")

	refreshedCount := 0
	err = u.Spin("Refreshing extensions", func() error {
		for _, toolKey := range tools {
			tool, _ := reg.GetTool(toolKey)
			globalPath := tool.ResolveGlobalPath()
			localPath := tool.ResolveLocalPath(projectRoot)

			// Check if installed globally (check first command)
			globalInstalled := false
			for _, c := range commands {
				if _, err := os.Stat(filepath.Join(globalPath, tool.Conventions.CommandPath(c, true))); err == nil {
					globalInstalled = true
					break
				}
			}

			// Check if installed locally
			localInstalled := false
			for _, c := range commands {
				if _, err := os.Stat(filepath.Join(localPath, tool.Conventions.CommandPath(c, false))); err == nil {
					localInstalled = true
					break
				}
			}

			// Refresh installations
			if globalInstalled {
				if _, err := inst.Install(toolKey, installer.ScopeGlobal); err == nil {
					refreshedCount++
				}
			}
			if localInstalled {
				if _, err := inst.Install(toolKey, installer.ScopeLocal); err == nil {
					refreshedCount++
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println()
	if refreshedCount > 0 {
		u.Success(fmt.Sprintf("Refreshed %d installation(s)", refreshedCount))
	} else {
		u.Info("No extensions currently installed")
	}

	return nil
}
