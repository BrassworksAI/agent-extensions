package sdd

import "testing"

func TestTasks_StartNextTask(t *testing.T) {
	tasks := &Tasks{Task: []*Task{
		{Name: "one", Status: TaskPending},
		{Name: "two", Status: TaskPending},
	}}

	name, task, err := tasks.StartNextTask()
	if err != nil {
		t.Fatalf("StartNextTask() error = %v", err)
	}
	if name != "one" {
		t.Fatalf("StartNextTask() name = %q, want %q", name, "one")
	}
	if task == nil || task.Status != TaskInProgress {
		t.Fatalf("StartNextTask() did not set task in progress")
	}
}

func TestTasks_StartNextTask_AlreadyInProgress(t *testing.T) {
	tasks := &Tasks{Task: []*Task{
		{Name: "one", Status: TaskInProgress},
		{Name: "two", Status: TaskPending},
	}}

	_, _, err := tasks.StartNextTask()
	if err == nil {
		t.Fatal("expected error when task already in progress")
	}
}

func TestTasks_CompleteCurrentTask(t *testing.T) {
	tasks := &Tasks{Task: []*Task{
		{Name: "one", Status: TaskInProgress},
		{Name: "two", Status: TaskPending},
	}}

	name, task, err := tasks.CompleteCurrentTask()
	if err != nil {
		t.Fatalf("CompleteCurrentTask() error = %v", err)
	}
	if name != "one" {
		t.Fatalf("CompleteCurrentTask() name = %q, want %q", name, "one")
	}
	if task == nil || task.Status != TaskComplete {
		t.Fatalf("CompleteCurrentTask() did not set task complete")
	}
}

func TestTasks_CurrentTaskStrict_MultipleInProgress(t *testing.T) {
	tasks := &Tasks{Task: []*Task{
		{Name: "one", Status: TaskInProgress},
		{Name: "two", Status: TaskInProgress},
	}}

	_, _, err := tasks.CurrentTaskStrict()
	if err == nil {
		t.Fatal("expected error for multiple in-progress tasks")
	}
}
