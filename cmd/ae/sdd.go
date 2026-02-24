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

var sddPhaseCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete current phase (optionally advance next)",
	RunE:  runSddPhaseComplete,
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

var sddTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task workflow management",
}

var sddTaskListCmd = &cobra.Command{
	Use:   "list [name]",
	Short: "List tasks with statuses",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSddTaskList,
}

var sddTaskCurrentCmd = &cobra.Command{
	Use:   "current [name]",
	Short: "Show current in-progress task details",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSddTaskCurrent,
}

var sddTaskNextCmd = &cobra.Command{
	Use:   "next [name]",
	Short: "Show next pending task details",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSddTaskNext,
}

var sddTaskStartCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start the next pending task",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSddTaskStart,
}

var sddTaskCompleteCmd = &cobra.Command{
	Use:   "complete [name]",
	Short: "Complete current task (optionally start next)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSddTaskComplete,
}

var (
	flagLane              string
	flagPhaseCompleteNext bool
	flagTaskCompleteNext  bool
)

func init() {
	sddInitCmd.Flags().StringVarP(&flagLane, "lane", "l", "full", "Lane type (full, vibe, bug)")
	sddPhaseCompleteCmd.Flags().BoolVar(&flagPhaseCompleteNext, "next", false, "Advance to next phase after completing current")
	sddTaskCompleteCmd.Flags().BoolVar(&flagTaskCompleteNext, "next", false, "Start next pending task after completing current")

	sddPhaseCmd.AddCommand(sddPhaseNextCmd)
	sddPhaseCmd.AddCommand(sddPhaseCompleteCmd)
	sddPhaseCmd.AddCommand(sddPhaseSetCmd)

	sddPendingCmd.AddCommand(sddPendingAddCmd)
	sddPendingCmd.AddCommand(sddPendingClearCmd)

	sddNotesCmd.AddCommand(sddNotesSetCmd)

	sddTaskCmd.AddCommand(sddTaskListCmd)
	sddTaskCmd.AddCommand(sddTaskCurrentCmd)
	sddTaskCmd.AddCommand(sddTaskNextCmd)
	sddTaskCmd.AddCommand(sddTaskStartCmd)
	sddTaskCmd.AddCommand(sddTaskCompleteCmd)

	sddCmd.AddCommand(sddInitCmd)
	sddCmd.AddCommand(sddStatusCmd)
	sddCmd.AddCommand(sddPhaseCmd)
	sddCmd.AddCommand(sddPendingCmd)
	sddCmd.AddCommand(sddNotesCmd)
	sddCmd.AddCommand(sddTaskCmd)

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

	createdTitle := sdd.HumanizeName(name)
	fmt.Printf("✓ Created change set: %s (%s lane)\n", createdTitle, sdd.LaneLabel(lane))
	fmt.Printf("  name: %s\n", name)
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

	if state.Phase.Status != sdd.StatusComplete {
		return fmt.Errorf("cannot advance phase %q while status is %q; run `ae sdd phase complete` after finishing work", state.Phase.Current, state.Phase.Status)
	}

	nextPhase, err := resolveNextPhase(changeDir, state)
	if err != nil {
		return err
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

func runSddPhaseComplete(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSet("")
	if err != nil {
		return err
	}

	state, err := sdd.LoadState(changeDir)
	if err != nil {
		return err
	}

	if state.Phase.Status == sdd.StatusComplete && !flagPhaseCompleteNext {
		fmt.Printf("Phase %q is already complete\n", state.Phase.Current)
		return nil
	}

	if err := validatePhaseCompletion(changeDir, state); err != nil {
		return err
	}

	state.CompletePhase()

	if !flagPhaseCompleteNext {
		if err := state.Save(changeDir); err != nil {
			return err
		}
		fmt.Printf("✓ Completed phase: %s\n", state.Phase.Current)
		return nil
	}

	nextPhase, err := resolveNextPhase(changeDir, state)
	if err != nil {
		if err := state.Save(changeDir); err != nil {
			return err
		}
		fmt.Printf("✓ Completed phase: %s\n", state.Phase.Current)
		fmt.Printf("→ No next phase: %v\n", err)
		return nil
	}

	completedPhase := state.Phase.Current
	if err := state.SetPhase(nextPhase); err != nil {
		return err
	}
	if err := state.Save(changeDir); err != nil {
		return err
	}

	fmt.Printf("✓ Completed phase: %s\n", completedPhase)
	fmt.Printf("✓ Advanced to phase: %s\n", nextPhase)
	return nil
}

func validatePhaseCompletion(changeDir string, state *sdd.State) error {
	if state.Phase.Status == sdd.StatusComplete {
		return nil
	}

	if state.Phase.Current != "implement" || state.Change.Lane != sdd.LaneFull {
		return nil
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	total, _, _, _ := tasks.Stats()
	if total == 0 {
		return fmt.Errorf("cannot complete implement phase: no tasks defined")
	}

	if name, currentTask, err := tasks.CurrentTaskStrict(); err != nil {
		return err
	} else if currentTask != nil {
		if name == "" {
			name = "current task"
		}
		return fmt.Errorf("cannot complete implement phase: complete %q first", name)
	}

	if !tasks.AllComplete() {
		return fmt.Errorf("cannot complete implement phase: complete all tasks first")
	}

	return nil
}

func resolveNextPhase(changeDir string, state *sdd.State) (string, error) {
	var nextPhase string
	if state.Phase.Current == "implement" {
		tasks, err := sdd.LoadTasks(changeDir)
		if err != nil {
			return "", err
		}

		total, _, _, _ := tasks.Stats()
		if total == 0 {
			return "", fmt.Errorf("cannot advance from implement: no tasks defined")
		}

		if name, currentTask, err := tasks.CurrentTaskStrict(); err != nil {
			return "", err
		} else if currentTask != nil {
			if name == "" {
				name = "current task"
			}
			return "", fmt.Errorf("cannot advance from implement: complete %q first", name)
		}

		if tasks.AllComplete() {
			candidate, ok := state.NextPhase()
			if !ok {
				return "", fmt.Errorf("already at final phase")
			}
			nextPhase = candidate
		} else {
			nextPhase = "plan"
		}
	} else {
		candidate, ok := state.NextPhase()
		if !ok {
			return "", fmt.Errorf("already at final phase")
		}
		nextPhase = candidate
	}

	return nextPhase, nil
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

func runSddTaskList(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSetFromArgs(args)
	if err != nil {
		return err
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	if len(tasks.Task) == 0 {
		fmt.Println("No tasks defined in tasks.toml")
		return nil
	}

	fmt.Printf("Tasks for %s\n\n", filepath.Base(changeDir))
	for _, task := range tasks.Task {
		if task == nil {
			continue
		}
		name := task.Name
		if name == "" {
			name = task.Title
		}
		fmt.Printf("%s %s\n", taskSymbol(task.Status), name)
	}

	return nil
}

func runSddTaskCurrent(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSetFromArgs(args)
	if err != nil {
		return err
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	name, task, err := tasks.CurrentTaskStrict()
	if err != nil {
		return err
	}
	if task == nil {
		fmt.Println("No task is currently in progress")
		return nil
	}

	printTaskDetails("Current task", name, task)
	return nil
}

func runSddTaskNext(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSetFromArgs(args)
	if err != nil {
		return err
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	name, task := tasks.NextPendingTask()
	if task == nil {
		fmt.Println("No pending tasks")
		return nil
	}

	printTaskDetails("Next task", name, task)
	return nil
}

func runSddTaskStart(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSetFromArgs(args)
	if err != nil {
		return err
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	name, task, err := tasks.StartNextTask()
	if err != nil {
		return err
	}

	if err := tasks.Save(changeDir); err != nil {
		return err
	}

	fmt.Printf("✓ Started task: %s\n", name)
	printTaskDetails("In progress", name, task)
	return nil
}

func runSddTaskComplete(cmd *cobra.Command, args []string) error {
	changeDir, err := resolveChangeSetFromArgs(args)
	if err != nil {
		return err
	}

	tasks, err := sdd.LoadTasks(changeDir)
	if err != nil {
		return err
	}

	completedName, _, err := tasks.CompleteCurrentTask()
	if err != nil {
		return err
	}

	if !flagTaskCompleteNext {
		if err := tasks.Save(changeDir); err != nil {
			return err
		}
		fmt.Printf("✓ Completed task: %s\n", completedName)
		return nil
	}

	nextName, nextTask, startErr := tasks.StartNextTask()
	if err := tasks.Save(changeDir); err != nil {
		return err
	}

	fmt.Printf("✓ Completed task: %s\n", completedName)
	if startErr != nil {
		fmt.Printf("→ No next task started: %v\n", startErr)
		return nil
	}

	fmt.Printf("✓ Started next task: %s\n", nextName)
	printTaskDetails("In progress", nextName, nextTask)
	return nil
}

func resolveChangeSetFromArgs(args []string) (string, error) {
	if len(args) > 0 {
		return resolveChangeSet(args[0])
	}
	return resolveChangeSet("")
}

func taskSymbol(status sdd.TaskStatus) string {
	switch status {
	case sdd.TaskComplete:
		return "✓"
	case sdd.TaskInProgress:
		return "◐"
	default:
		return "○"
	}
}

func printTaskDetails(label, name string, task *sdd.Task) {
	fmt.Printf("%s\n", label)
	fmt.Printf("  Name: %s\n", name)
	if task.Title != "" {
		fmt.Printf("  Title: %s\n", task.Title)
	}
	fmt.Printf("  Status: %s\n", task.Status)
	if task.Description != "" {
		fmt.Printf("  Description: %s\n", task.Description)
	}
	if len(task.SpecRequirements) > 0 {
		fmt.Println("  Spec Requirements:")
		for _, specReq := range task.SpecRequirements {
			specName := specReq.Spec
			if specName == "" {
				specName = "(unspecified spec)"
			}
			fmt.Printf("    %s\n", specName)
			for _, req := range specReq.Requirements {
				fmt.Printf("      - %s\n", req)
			}
		}
	} else if len(task.Requirements) > 0 {
		fmt.Println("  Requirements:")
		for _, req := range task.Requirements {
			fmt.Printf("    - %s\n", req)
		}
	}
}
