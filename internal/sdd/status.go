package sdd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type StatusRenderer struct {
	State *State
	Tasks *Tasks
}

func NewStatusRenderer(state *State, tasks *Tasks) *StatusRenderer {
	return &StatusRenderer{State: state, Tasks: tasks}
}

func (r *StatusRenderer) Render() {
	r.renderHeader()
	r.renderPhaseProgress()
	r.renderTaskProgress()
	r.renderNotes()
	r.renderPending()
	r.renderNextAction()
}

func (r *StatusRenderer) renderHeader() {
	header := fmt.Sprintf("%s (%s lane)", r.State.Change.Name, r.State.Change.Lane)
	gumStyle(header, "--foreground", "212", "--bold")
	fmt.Println()
}

func (r *StatusRenderer) renderPhaseProgress() {
	phases := PhasesForLane(r.State.Change.Lane)
	currentIdx := r.State.PhaseIndex()

	var symbols []string
	var labels []string

	for i, phase := range phases {
		var symbol string
		if i < currentIdx {
			symbol = "●"
		} else if i == currentIdx {
			if r.State.Phase.Status == StatusComplete {
				symbol = "●"
			} else {
				symbol = "◐"
			}
		} else {
			symbol = "○"
		}
		symbols = append(symbols, symbol)

		label := phase
		if len(label) > 4 {
			label = label[:4]
		}
		labels = append(labels, fmt.Sprintf("%-4s", label))
	}

	phaseLine := strings.Join(symbols, " ─── ")
	labelLine := strings.Join(labels, "  ")

	gumStyle(phaseLine, "--foreground", "212")
	gumStyle(labelLine, "--foreground", "240")
	fmt.Println()
}

func (r *StatusRenderer) renderTaskProgress() {
	if r.Tasks == nil || len(r.Tasks.Task) == 0 {
		return
	}

	total, complete, _, _ := r.Tasks.Stats()
	if total == 0 {
		return
	}

	percent := float64(complete) / float64(total) * 100
	barWidth := 30
	filled := int(float64(barWidth) * float64(complete) / float64(total))
	empty := barWidth - filled

	bar := fmt.Sprintf("Tasks %s%s %.0f%% (%d/%d)",
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		percent, complete, total)

	gumStyle(bar, "--foreground", "212")

	for _, task := range r.Tasks.Task {
		if task == nil {
			continue
		}
		var symbol string
		switch task.Status {
		case TaskComplete:
			symbol = "✓"
		case TaskInProgress:
			symbol = "◐"
		default:
			symbol = "○"
		}
		name := task.Name
		if name == "" {
			name = task.Title
		}
		fmt.Printf("%s %s\n", symbol, name)
	}
	fmt.Println()
}

func (r *StatusRenderer) renderNotes() {
	if r.State.Notes.Content == "" {
		return
	}

	content := "Notes\n" + r.State.Notes.Content
	gumStyle(content,
		"--border", "rounded",
		"--padding", "0 1",
		"--border-foreground", "240")
	fmt.Println()
}

func (r *StatusRenderer) renderPending() {
	if len(r.State.Pending.Items) == 0 {
		return
	}

	gumStyle("Pending:", "--foreground", "214")
	for _, item := range r.State.Pending.Items {
		fmt.Printf("  • %s\n", item)
	}
	fmt.Println()
}

func (r *StatusRenderer) renderNextAction() {
	var next string

	if r.State.Phase.Status == StatusComplete {
		if nextPhase, ok := r.State.NextPhase(); ok {
			next = fmt.Sprintf("ae sdd phase next  (→ %s)", nextPhase)
		} else {
			next = "Change set complete!"
		}
	} else {
		name, currentTask := r.Tasks.CurrentTask()
		if currentTask != nil {
			if name == "" {
				name = "current task"
			}
			next = fmt.Sprintf("Update tasks.toml to set %s status=complete", name)
		} else {
			next = fmt.Sprintf("Continue %s phase", r.State.Phase.Current)
		}
	}

	gumStyle("Next: "+next, "--foreground", "240", "--italic")
}

func gumStyle(text string, args ...string) {
	cmdArgs := append([]string{"style"}, args...)
	cmdArgs = append(cmdArgs, text)

	cmd := exec.Command("gum", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(text)
	}
}
