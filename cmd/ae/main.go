package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shanepadgett/agent-extensions/internal/config"
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
	Long:  "Manage installation of skills for AI coding agents",
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

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage local and global cache files",
}

var cacheRepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Rebuild cache files and metadata",
	RunE:  runCacheRepair,
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
	flagTools      []string
	flagScope      string
	flagYes        bool
	flagDoctorFix  bool
	flagCacheScope string
)

func init() {
	installCmd.Flags().StringSliceVarP(&flagTools, "tools", "t", nil, "Tools to install to (comma-separated)")
	installCmd.Flags().StringVarP(&flagScope, "scope", "s", "", "Installation scope (global, local, both)")
	installCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation")

	uninstallCmd.Flags().StringSliceVarP(&flagTools, "tools", "t", nil, "Tools to uninstall from (comma-separated)")
	uninstallCmd.Flags().StringVarP(&flagScope, "scope", "s", "", "Uninstallation scope (global, local, both)")
	uninstallCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation")

	doctorCmd.Flags().BoolVar(&flagDoctorFix, "fix", false, "Repair cache issues when found")

	cacheRepairCmd.Flags().StringVarP(&flagCacheScope, "scope", "s", "both", "Cache scope (global, local, both)")
	cacheCmd.AddCommand(cacheRepairCmd)

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cacheCmd)
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
		confirmMsg := fmt.Sprintf("Install %d skills to %d tools (%s)?",
			totalSkills, len(selectedTools), scope)
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
	inst.SetSource(version)
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
			lines = append(lines, fmt.Sprintf("✓ %s: %d skills", tool.Name, result.Skills))
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
		confirmed, err := u.Confirm(fmt.Sprintf("Uninstall %d skills from %d tools (%s)?",
			totalSkills, len(selectedTools), scope))
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

			if result.Skills > 0 {
				tool, _ := reg.GetTool(toolName)
				lines = append(lines, fmt.Sprintf("✓ %s: removed %d skills", tool.Name, result.Skills))
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
	skills := reg.GetAllSkills()

	u.Header("\nAvailable Content")
	fmt.Printf("  Skills: %d\n", len(skills))

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

		status.global = toolHasInstalledContent(tool, globalPath, projectRoot, true, nil, skills)
		status.local = toolHasInstalledContent(tool, localPath, projectRoot, false, nil, skills)

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

func resolveScopedSkillPath(targetBase, projectRoot, conventionPath string, isGlobal, hasScopeOverride bool) string {
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

	inst := installer.New(reg, getProjectRoot())
	inst.SetSource(version)

	// Check cache directories and integrity
	u.Header("\nCache:")
	cacheResults := make([]*installer.CacheCheckResult, 0, 2)
	for _, scope := range []installer.Scope{installer.ScopeGlobal, installer.ScopeLocal} {
		result, err := inst.CheckCache(scope, version)
		if err != nil {
			u.Error(fmt.Sprintf("%s cache check failed: %v", scopeTitle(scope), err))
			continue
		}
		cacheResults = append(cacheResults, result)

		if !result.Exists {
			u.Info(fmt.Sprintf("%s cache: not created yet (%s)", scopeTitle(scope), result.Path))
			continue
		}

		if len(result.Issues) == 0 && !result.SourceMismatch {
			u.Success(fmt.Sprintf("%s cache: healthy (source %s, 0 mismatches)", scopeTitle(scope), displaySource(result.Source)))
			continue
		}

		u.Warn(fmt.Sprintf("%s cache: issues found", scopeTitle(scope)))
		if result.SourceMismatch {
			u.Warn(fmt.Sprintf("source mismatch: installed by %s, running %s", displaySource(result.Source), version))
		}
		for _, issue := range result.Issues {
			switch issue.Kind {
			case "hash mismatch", "metadata hash mismatch":
				u.Warn(fmt.Sprintf("%s: %s", issue.Kind, issue.Path))
			case "missing file":
				u.Warn(fmt.Sprintf("missing file: %s", issue.Path))
			case "metadata missing":
				u.Warn("metadata missing")
			default:
				u.Warn(fmt.Sprintf("%s: %s", issue.Kind, issue.Path))
			}
		}
		u.Info(fmt.Sprintf("Run: ae cache repair --scope %s", scope))
	}

	if flagDoctorFix {
		u.Header("\nCache Repair:")
		for _, result := range cacheResults {
			if !result.Exists {
				continue
			}
			if len(result.Issues) == 0 && !result.SourceMismatch {
				continue
			}

			scope := result.Scope
			err := u.Spin(fmt.Sprintf("Repairing %s cache", strings.ToLower(scopeTitle(scope))), func() error {
				_, repairErr := inst.RepairCache(scope)
				return repairErr
			})
			if err != nil {
				u.Error(fmt.Sprintf("%s cache repair failed: %v", scopeTitle(scope), err))
				continue
			}
			u.Success(fmt.Sprintf("%s cache repaired", scopeTitle(scope)))
		}

		u.Info("Re-running cache health checks...")
		allHealthy := true
		for _, scope := range []installer.Scope{installer.ScopeGlobal, installer.ScopeLocal} {
			result, err := inst.CheckCache(scope, version)
			if err != nil {
				allHealthy = false
				u.Error(fmt.Sprintf("%s cache re-check failed: %v", scopeTitle(scope), err))
				continue
			}
			if result.Exists && len(result.Issues) == 0 && !result.SourceMismatch {
				u.Success(fmt.Sprintf("%s cache: healthy", scopeTitle(scope)))
				continue
			}
			allHealthy = false
			if !result.Exists {
				u.Info(fmt.Sprintf("%s cache: not created yet", scopeTitle(scope)))
			} else {
				u.Warn(fmt.Sprintf("%s cache: still has issues", scopeTitle(scope)))
			}
		}
		if allHealthy {
			u.Success("All caches healthy")
		}

		u.Header("\nInstall Repair:")
		for _, scope := range []installer.Scope{installer.ScopeGlobal, installer.ScopeLocal} {
			installed := installedToolsForScope(reg, getProjectRoot(), scope)
			if len(installed) == 0 {
				u.Info(fmt.Sprintf("%s installs: none detected", scopeTitle(scope)))
				continue
			}

			refreshErrors := make([]string, 0)
			err := u.Spin(fmt.Sprintf("Refreshing %s installs", strings.ToLower(scopeTitle(scope))), func() error {
				for _, toolKey := range installed {
					if _, installErr := inst.Install(toolKey, scope); installErr != nil {
						refreshErrors = append(refreshErrors, fmt.Sprintf("%s: %v", toolKey, installErr))
					}
				}
				return nil
			})
			if err != nil {
				u.Error(fmt.Sprintf("%s install refresh failed: %v", scopeTitle(scope), err))
				continue
			}

			if len(refreshErrors) == 0 {
				u.Success(fmt.Sprintf("%s installs refreshed (%d tool(s))", scopeTitle(scope), len(installed)))
				continue
			}

			u.Warn(fmt.Sprintf("%s installs refreshed with issues", scopeTitle(scope)))
			for _, issue := range refreshErrors {
				u.Warn(issue)
			}
		}
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

func runCacheRepair(cmd *cobra.Command, args []string) error {
	u := ui.New()

	reg, err := getRegistry()
	if err != nil {
		return fmt.Errorf("loading registry: %w", err)
	}

	scope, err := parseScope(flagCacheScope)
	if err != nil {
		return err
	}

	inst := installer.New(reg, getProjectRoot())
	inst.SetSource(version)

	scopes := []installer.Scope{scope}
	if scope == installer.ScopeBoth {
		scopes = []installer.Scope{installer.ScopeGlobal, installer.ScopeLocal}
	}

	for _, s := range scopes {
		err := u.Spin(fmt.Sprintf("Rebuilding %s cache", strings.ToLower(scopeTitle(s))), func() error {
			_, repairErr := inst.RepairCache(s)
			return repairErr
		})
		if err != nil {
			return err
		}

		result, err := inst.CheckCache(s, version)
		if err != nil {
			return err
		}
		if len(result.Issues) == 0 && !result.SourceMismatch {
			u.Success(fmt.Sprintf("%s cache repaired (%s)", scopeTitle(s), result.Path))
		} else {
			u.Warn(fmt.Sprintf("%s cache repaired with warnings", scopeTitle(s)))
		}
	}

	return nil
}

func parseScope(value string) (installer.Scope, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "global":
		return installer.ScopeGlobal, nil
	case "local":
		return installer.ScopeLocal, nil
	case "both", "":
		return installer.ScopeBoth, nil
	default:
		return "", fmt.Errorf("unknown scope: %s", value)
	}
}

func scopeTitle(scope installer.Scope) string {
	if scope == installer.ScopeGlobal {
		return "Global"
	}
	return "Local"
}

func displaySource(source string) string {
	if strings.TrimSpace(source) == "" {
		return "unknown"
	}
	return source
}

func installedToolsForScope(reg *registry.Registry, projectRoot string, scope installer.Scope) []string {
	tools := reg.GetToolNames()
	sort.Strings(tools)
	commands := reg.GetAllCommands()
	skills := reg.GetAllSkills()

	installed := make([]string, 0)
	isGlobal := scope == installer.ScopeGlobal

	for _, toolKey := range tools {
		tool, _ := reg.GetTool(toolKey)
		targetBase := tool.ResolveLocalPath(projectRoot)
		if isGlobal {
			targetBase = tool.ResolveGlobalPath()
		}

		if toolHasInstalledContent(tool, targetBase, projectRoot, isGlobal, commands, skills) {
			installed = append(installed, toolKey)
		}
	}

	return installed
}

func toolHasInstalledContent(tool config.Tool, targetBase, projectRoot string, isGlobal bool, commands, skills []string) bool {
	for _, c := range commands {
		path := filepath.Join(targetBase, tool.Conventions.CommandPath(c, isGlobal))
		if _, err := os.Lstat(path); err == nil {
			return true
		}
	}

	for _, s := range skills {
		pattern := tool.Conventions.ScopedSkillPath(s, isGlobal)
		hasScopeOverride := (isGlobal && tool.Conventions.GlobalSkills != "") || (!isGlobal && tool.Conventions.LocalSkills != "")
		fullPath := resolveScopedSkillPath(targetBase, projectRoot, pattern, isGlobal, hasScopeOverride)
		checkPath := fullPath
		if filepath.Ext(pattern) != ".md" {
			checkPath = filepath.Dir(fullPath)
		}
		if _, err := os.Lstat(checkPath); err == nil {
			return true
		}
	}

	return false
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
	skills := reg.GetAllSkills()

	// Find what's currently installed and refresh symlinks
	inst := installer.New(reg, projectRoot)
	inst.SetSource(version)

	fmt.Println()
	u.Info("Refreshing installed extensions...")

	refreshedCount := 0
	err = u.Spin("Refreshing extensions", func() error {
		for _, toolKey := range tools {
			tool, _ := reg.GetTool(toolKey)
			globalPath := tool.ResolveGlobalPath()
			localPath := tool.ResolveLocalPath(projectRoot)

			globalInstalled := toolHasInstalledContent(tool, globalPath, projectRoot, true, commands, skills)
			localInstalled := toolHasInstalledContent(tool, localPath, projectRoot, false, commands, skills)

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
