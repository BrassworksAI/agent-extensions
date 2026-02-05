package sdd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskInProgress TaskStatus = "in_progress"
	TaskComplete   TaskStatus = "complete"
)

type Task struct {
	Name         string     `toml:"name"`
	Title        string     `toml:"title"`
	Description  string     `toml:"description"`
	Status       TaskStatus `toml:"status"`
	Requirements []string   `toml:"requirements"`
}

type Tasks struct {
	Task []*Task `toml:"task"`
}

func NewTasks() *Tasks {
	return &Tasks{
		Task: []*Task{},
	}
}

func LoadTasks(changeDir string) (*Tasks, error) {
	path := filepath.Join(changeDir, "tasks.toml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewTasks(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tasks.toml: %w", err)
	}

	var tasks Tasks
	tasksErr := toml.Unmarshal(data, &tasks)
	if tasksErr == nil {
		if tasks.Task == nil {
			tasks.Task = []*Task{}
		}
		if len(tasks.Task) > 0 {
			normalizeTasks(tasks.Task)
			return &tasks, nil
		}
	}

	var legacy legacyTasks
	legacyErr := toml.Unmarshal(data, &legacy)
	if legacyErr == nil && len(legacy.Task) > 0 {
		names := make([]string, 0, len(legacy.Task))
		for name := range legacy.Task {
			names = append(names, name)
		}
		sort.Strings(names)

		converted := NewTasks()
		for _, name := range names {
			task := legacy.Task[name]
			if task == nil {
				continue
			}
			normalizeTask(task, name)
			converted.Task = append(converted.Task, task)
		}

		return converted, nil
	}

	if tasksErr == nil {
		normalizeTasks(tasks.Task)
		return &tasks, nil
	}

	if legacyErr != nil {
		return nil, fmt.Errorf("parsing tasks.toml: %w", legacyErr)
	}

	return nil, fmt.Errorf("parsing tasks.toml: %w", tasksErr)
}

func (t *Tasks) Save(changeDir string) error {
	path := filepath.Join(changeDir, "tasks.toml")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating tasks.toml: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(t); err != nil {
		return fmt.Errorf("encoding tasks.toml: %w", err)
	}

	return nil
}

func (t *Tasks) Add(shortName string, task *Task) error {
	if task == nil {
		return fmt.Errorf("task %q is nil", shortName)
	}
	if shortName == "" {
		return fmt.Errorf("task name must be non-empty")
	}
	if _, exists := t.Get(shortName); exists {
		return fmt.Errorf("task %q already exists", shortName)
	}
	if task.Name == "" {
		task.Name = shortName
	}
	if task.Status == "" {
		task.Status = TaskPending
	}
	t.Task = append(t.Task, task)
	return nil
}

func (t *Tasks) Get(shortName string) (*Task, bool) {
	for _, task := range t.Task {
		if task == nil {
			continue
		}
		if task.Name == shortName {
			return task, true
		}
	}
	return nil, false
}

func (t *Tasks) Start(shortName string) error {
	task, exists := t.Get(shortName)
	if !exists {
		return fmt.Errorf("task %q not found", shortName)
	}
	task.Status = TaskInProgress
	return nil
}

func (t *Tasks) Complete(shortName string) error {
	task, exists := t.Get(shortName)
	if !exists {
		return fmt.Errorf("task %q not found", shortName)
	}
	task.Status = TaskComplete
	return nil
}

func (t *Tasks) List() []string {
	names := make([]string, 0, len(t.Task))
	for _, task := range t.Task {
		if task == nil {
			continue
		}
		name := task.Name
		if name == "" {
			name = task.Title
		}
		names = append(names, name)
	}
	return names
}

func (t *Tasks) Stats() (total, complete, inProgress, pending int) {
	for _, task := range t.Task {
		if task == nil {
			continue
		}
		total++
		switch task.Status {
		case TaskComplete:
			complete++
		case TaskInProgress:
			inProgress++
		case TaskPending:
			pending++
		}
	}
	return
}

func (t *Tasks) CurrentTask() (string, *Task) {
	for _, task := range t.Task {
		if task == nil {
			continue
		}
		if task.Status == TaskInProgress {
			return task.Name, task
		}
	}
	return "", nil
}

func (t *Tasks) AllComplete() bool {
	for _, task := range t.Task {
		if task == nil {
			continue
		}
		if task.Status != TaskComplete {
			return false
		}
	}
	return len(t.Task) > 0
}

type legacyTasks struct {
	Task map[string]*Task `toml:"task"`
}

func normalizeTasks(tasks []*Task) {
	for i, task := range tasks {
		if task == nil {
			continue
		}
		normalizeTask(task, fmt.Sprintf("task-%d", i+1))
	}
}

func normalizeTask(task *Task, fallbackName string) {
	if task.Name == "" {
		task.Name = fallbackName
	}
	if task.Status == "" {
		task.Status = TaskPending
	}
}
