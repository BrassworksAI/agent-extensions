package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/shanepadgett/agent-extensions/internal/sdd"
	"github.com/spf13/cobra"
)

var sddCmd = &cobra.Command{
	Use:   "sdd",
	Short: "SDD change set management",
}

var sddInitCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Initialize a new change set",
	Args:  cobra.ExactArgs(1),
	RunE:  runSddInit,
}

var sddStatusCmd = &cobra.Command{
	Use:   "status [name]",
	Short: "Show change set status",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSddStatus,
}

var sddPhaseCmd = &cobra.Command{
	Use:   "phase",
	Short: "Phase management",
}

var sddPhaseNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Advance to next phase",
	RunE:  runSddPhaseNext,
}

var sddPhaseSetCmd = &cobra.Command{
	Use:   "set <phase>",
	Short: "Set current phase",
	Args:  cobra.ExactArgs(1),
	RunE:  runSddPhaseSet,
}

var sddPendingCmd = &cobra.Command{
	Use:   "pending",
	Short: "Pending items management",
}

var sddPendingAddCmd = &cobra.Command{
	Use:   "add <item>",
	Short: "Add a pending item",
	Args:  cobra.ExactArgs(1),
	RunE:  runSddPendingAdd,
}

var sddPendingClearCmd = &cobra.Command{
	Use:   "clear <index>",
	Short: "Clear a pending item by index",
	Args:  cobra.ExactArgs(1),
	RunE:  runSddPendingClear,
}

var sddNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Notes management",
}

var sddNotesSetCmd = &cobra.Command{
	Use:   "set <content>",
	Short: "Set notes content",
	Args:  cobra.ExactArgs(1),
	RunE:  runSddNotesSet,
}

var (
	flagLane string
)

func init() {
	sddInitCmd.Flags().StringVarP(&flagLane, "lane", "l", "full", "Lane type (full, vibe, bug)")

	sddPhaseCmd.AddCommand(sddPhaseNextCmd)
	sddPhaseCmd.AddCommand(sddPhaseSetCmd)

	sddPendingCmd.AddCommand(sddPendingAddCmd)
	sddPendingCmd.AddCommand(sddPendingClearCmd)

	sddNotesCmd.AddCommand(sddNotesSetCmd)

	sddCmd.AddCommand(sddInitCmd)
	sddCmd.AddCommand(sddStatusCmd)
	sddCmd.AddCommand(sddPhaseCmd)
	sddCmd.AddCommand(sddPendingCmd)
	sddCmd.AddCommand(sddNotesCmd)

	rootCmd.AddCommand(sddCmd)
}

func getChangesDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "changes")
}

func resolveChangeSet(name string) (string, error) {
	changesDir := getChangesDir()

	if name != "" {
		dir := filepath.Join(changesDir, name)
		if _, err := os.Stat(dir); err != nil {
			return "", fmt.Errorf("change set %q not found", name)
		}
		return dir, nil
	}

	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return "", fmt.Errorf("no changes directory found")
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			dirs = append(dirs, e.Name())
		}
	}

	if len(dirs) == 0 {
		return "", fmt.Errorf("no change sets found")
	}
	if len(dirs) == 1 {
		return filepath.Join(changesDir, dirs[0]), nil
	}

	return "", fmt.Errorf("multiple change sets found, please specify one: %v", dirs)
}

func isKebabCase(s string) bool {
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`, s)
	return matched
}

func runSddInit(cmd *cobra.Command, args []string) error {
	name := args[0]

	if !isKebabCase(name) {
		return fmt.Errorf("name must be kebab-case (lowercase, hyphens only)")
	}

	lane := sdd.Lane(flagLane)
	if lane != sdd.LaneFull && lane != sdd.LaneVibe && lane != sdd.LaneBug {
		return fmt.Errorf("invalid lane: %s (must be full, vibe, or bug)", flagLane)
	}

	changeDir := filepath.Join(getChangesDir(), name)
	if _, err := os.Stat(changeDir); err == nil {
		return fmt.Errorf("change set %q already exists", name)
	}

	if err := os.MkdirAll(changeDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	state := sdd.NewState(name, lane)
	if err := state.Save(changeDir); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	fmt.Printf("✓ Created change set: %s (%s lane)\n", name, lane)
	fmt.Printf("  → %s\n", changeDir)
	return nil
}

func runSddStatus(cmd *cobra.Command, args []string) error {
	var name string
	if len(args) > 0 {
		name = args[0]
	}

	changeDir, err := resolveChangeSet(name)
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	renderer := sdd.NewStatusRenderer(state, tasks)
	renderer.Render()

	return nil
}

func runSddPhaseNext(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSet("")
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	nextPhase, ok := state.NextPhase()
	if !ok {
		return fmt.Errorf("already at final phase")
	}

	if err := state.SetPhase(nextPhase); err != nil {
		return err
	}
	if err := state.Save(changeDir); err != nil {
		return err
	}

	fmt.Printf("✓ Advanced to phase: %s\n", nextPhase)
	return nil
}

func runSddPhaseSet(cmd *cobra.Command, args []string) error {
	phase := args[0]

	changeDir, err := resolveChangeSet("")
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	if err := state.SetPhase(phase); err != nil {
		return err
	}

	if err := state.Save(changeDir); err != nil {
		return err
	}

	fmt.Printf("✓ Set phase to: %s\n", phase)
	return nil
}

func runSddPendingAdd(cmd *cobra.Command, args []string) error {
	item := args[0]

	changeDir, err := resolveChangeSet("")
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	state.AddPending(item)

	if err := state.Save(changeDir); err != nil {
		return err
	}

	fmt.Printf("✓ Added pending item: %s\n", item)
	return nil
}

func runSddPendingClear(cmd *cobra.Command, args []string) error {
	index, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid index: %s", args[0])
	}

	changeDir, err := resolveChangeSet("")
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	if err := state.ClearPending(index); err != nil {
		return err
	}

	if err := state.Save(changeDir); err != nil {
		return err
	}

	fmt.Println("✓ Cleared pending item")
	return nil
}

func runSddNotesSet(cmd *cobra.Command, args []string) error {
	content := args[0]

	changeDir, err := resolveChangeSet("")
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	state.SetNotes(content)

	if err := state.Save(changeDir); err != nil {
		return err
	}

	fmt.Println("✓ Updated notes")
	return nil
}
