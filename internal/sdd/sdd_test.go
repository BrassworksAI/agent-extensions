package sdd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateLifecycle(t *testing.T) {
	dir := t.TempDir()

	// Create new state
	state := NewState("test-feature", LaneFull)
	if state.Change.Name != "test-feature" {
		t.Errorf("expected name test-feature, got %s", state.Change.Name)
	}
	if state.Change.Lane != LaneFull {
		t.Errorf("expected lane full, got %s", state.Change.Lane)
	}
	if state.Phase.Current != "proposal" {
		t.Errorf("expected phase proposal, got %s", state.Phase.Current)
	}
	if state.Phase.Status != StatusInProgress {
		t.Errorf("expected status in_progress, got %s", state.Phase.Status)
	}

	// Save
	if err := state.Save(dir); err != nil {
		t.Fatalf("failed to save state: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(dir, "state.toml")); err != nil {
		t.Fatalf("state.toml not created: %v", err)
	}

	// Load
	loaded, err := LoadState(dir)
	if err != nil {
		t.Fatalf("failed to load state: %v", err)
	}
	if loaded.Change.Name != "test-feature" {
		t.Errorf("loaded name mismatch: %s", loaded.Change.Name)
	}
	if loaded.Phase.Current != "proposal" {
		t.Errorf("loaded phase mismatch: %s", loaded.Phase.Current)
	}
}

func TestPhaseTransitions(t *testing.T) {
	state := NewState("test", LaneFull)

	// Initial phase
	if state.Phase.Current != "proposal" {
		t.Errorf("expected proposal, got %s", state.Phase.Current)
	}

	// Next phase
	next, ok := state.NextPhase()
	if !ok || next != "specs" {
		t.Errorf("expected specs, got %s", next)
	}

	// Advance
	state.SetPhase("specs")
	if state.Phase.Current != "specs" {
		t.Errorf("expected specs after set, got %s", state.Phase.Current)
	}

	// Complete phase
	state.CompletePhase()
	if state.Phase.Status != StatusComplete {
		t.Errorf("expected complete status, got %s", state.Phase.Status)
	}

	// Invalid phase
	if err := state.SetPhase("invalid"); err == nil {
		t.Error("expected error for invalid phase")
	}
}

func TestPendingItems(t *testing.T) {
	state := NewState("test", LaneFull)

	state.AddPending("item 1")
	state.AddPending("item 2")

	if len(state.Pending.Items) != 2 {
		t.Errorf("expected 2 pending items, got %d", len(state.Pending.Items))
	}

	if err := state.ClearPending(0); err != nil {
		t.Errorf("failed to clear pending: %v", err)
	}

	if len(state.Pending.Items) != 1 {
		t.Errorf("expected 1 pending item after clear, got %d", len(state.Pending.Items))
	}

	if state.Pending.Items[0] != "item 2" {
		t.Errorf("expected item 2 to remain, got %s", state.Pending.Items[0])
	}

	// Invalid index
	if err := state.ClearPending(99); err == nil {
		t.Error("expected error for invalid index")
	}
}

func TestNotes(t *testing.T) {
	state := NewState("test", LaneFull)

	state.SetNotes("test notes")
	if state.Notes.Content != "test notes" {
		t.Errorf("expected 'test notes', got %s", state.Notes.Content)
	}

	state.CompletePhase()
	if state.Notes.Content != "" {
		t.Error("expected notes to be cleared on phase complete")
	}
}

func TestLanePhases(t *testing.T) {
	tests := []struct {
		lane  Lane
		first string
		count int
	}{
		{LaneFull, "proposal", 8},
		{LaneVibe, "context", 5},
		{LaneBug, "triage", 5},
	}

	for _, tt := range tests {
		phases := PhasesForLane(tt.lane)
		if phases[0] != tt.first {
			t.Errorf("lane %s: expected first phase %s, got %s", tt.lane, tt.first, phases[0])
		}
		if len(phases) != tt.count {
			t.Errorf("lane %s: expected %d phases, got %d", tt.lane, tt.count, len(phases))
		}
	}
}

func TestTasksLifecycle(t *testing.T) {
	dir := t.TempDir()

	// Create new tasks
	tasks := NewTasks()

	// Add tasks
	task1 := &Task{
		Title:        "DB Models",
		Description:  "Create database models",
		Requirements: []string{"Req 1", "Req 2"},
	}
	if err := tasks.Add("db-models", task1); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	task2 := &Task{
		Title:       "Auth Endpoint",
		Description: "Build auth endpoint",
	}
	if err := tasks.Add("auth-endpoint", task2); err != nil {
		t.Fatalf("failed to add task: %v", err)
	}

	// Duplicate add should fail
	if err := tasks.Add("db-models", task1); err == nil {
		t.Error("expected error for duplicate task")
	}

	// Save
	if err := tasks.Save(dir); err != nil {
		t.Fatalf("failed to save tasks: %v", err)
	}

	// Load
	loaded, err := LoadTasks(dir)
	if err != nil {
		t.Fatalf("failed to load tasks: %v", err)
	}
	if len(loaded.Task) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loaded.Task))
	}

	// Verify task content
	dbTask, ok := loaded.Get("db-models")
	if !ok {
		t.Fatal("db-models task not found")
	}
	if dbTask.Title != "DB Models" {
		t.Errorf("expected title 'DB Models', got %s", dbTask.Title)
	}
	if dbTask.Name != "db-models" {
		t.Errorf("expected name 'db-models', got %s", dbTask.Name)
	}
	if dbTask.Status != TaskPending {
		t.Errorf("expected pending status, got %s", dbTask.Status)
	}
}

func TestTaskStatusTransitions(t *testing.T) {
	tasks := NewTasks()
	tasks.Add("task-1", &Task{Title: "Task 1"})
	tasks.Add("task-2", &Task{Title: "Task 2"})

	// Start task
	if err := tasks.Start("task-1"); err != nil {
		t.Fatalf("failed to start task: %v", err)
	}

	task, _ := tasks.Get("task-1")
	if task.Status != TaskInProgress {
		t.Errorf("expected in_progress, got %s", task.Status)
	}

	// Current task
	name, current := tasks.CurrentTask()
	if name != "task-1" {
		t.Errorf("expected task-1 as current, got %s", name)
	}
	if current == nil {
		t.Error("current task should not be nil")
	}

	// Complete task
	if err := tasks.Complete("task-1"); err != nil {
		t.Fatalf("failed to complete task: %v", err)
	}

	task, _ = tasks.Get("task-1")
	if task.Status != TaskComplete {
		t.Errorf("expected complete, got %s", task.Status)
	}

	// Stats
	total, complete, inProgress, pending := tasks.Stats()
	if total != 2 || complete != 1 || inProgress != 0 || pending != 1 {
		t.Errorf("unexpected stats: total=%d complete=%d inProgress=%d pending=%d",
			total, complete, inProgress, pending)
	}

	// All complete check
	if tasks.AllComplete() {
		t.Error("AllComplete should be false")
	}

	tasks.Complete("task-2")
	if !tasks.AllComplete() {
		t.Error("AllComplete should be true")
	}

	// Invalid task
	if err := tasks.Start("nonexistent"); err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestLoadTasksNonexistent(t *testing.T) {
	dir := t.TempDir()

	// Should return empty tasks, not error
	tasks, err := LoadTasks(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks.Task) != 0 {
		t.Errorf("expected empty tasks, got %d", len(tasks.Task))
	}
}

func TestTasksList(t *testing.T) {
	tasks := NewTasks()
	tasks.Add("zebra", &Task{Title: "Z"})
	tasks.Add("alpha", &Task{Title: "A"})
	tasks.Add("beta", &Task{Title: "B"})

	list := tasks.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list))
	}

	// Should preserve insertion order
	if list[0] != "zebra" || list[1] != "alpha" || list[2] != "beta" {
		t.Errorf("list not in insertion order: %v", list)
	}
}
