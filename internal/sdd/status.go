package sdd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	fmt.Println()
	header := fmt.Sprintf("%s (%s lane)", HumanizeName(r.State.Change.Name), LaneLabel(r.State.Change.Lane))
	styleText(header, lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true))
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

	styleText(phaseLine, lipgloss.NewStyle().Foreground(lipgloss.Color("212")))
	styleText(labelLine, lipgloss.NewStyle().Foreground(lipgloss.Color("240")))
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

	styleText(bar, lipgloss.NewStyle().Foreground(lipgloss.Color("212")))
	fmt.Println()

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

	styleText("Notes:", lipgloss.NewStyle().Foreground(lipgloss.Color("81")))
	for _, line := range strings.Split(r.State.Notes.Content, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Printf("  • %s\n", line)
	}
	fmt.Println()
}

func (r *StatusRenderer) renderPending() {
	if len(r.State.Pending.Items) == 0 {
		return
	}

	styleText("Pending:", lipgloss.NewStyle().Foreground(lipgloss.Color("214")))
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
		if r.Tasks == nil {
			next = fmt.Sprintf("ae sdd phase complete  (current: %s)", r.State.Phase.Current)
			styleText("Next: "+next, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true))
			return
		}

		name, currentTask, err := r.Tasks.CurrentTaskStrict()
		if err != nil {
			next = err.Error()
		} else if currentTask != nil {
			if name == "" {
				name = "current task"
			}
			next = fmt.Sprintf("ae sdd task complete --next  (current: %s)", name)
		} else if nextName, nextTask := r.Tasks.NextPendingTask(); nextTask != nil {
			next = fmt.Sprintf("ae sdd task start  (next: %s)", nextName)
		} else if r.State.Phase.Current == "implement" {
			next = "ae sdd phase complete --next"
		} else {
			next = fmt.Sprintf("ae sdd phase complete  (current: %s)", r.State.Phase.Current)
		}
	}

	styleText("Next: "+next, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true))
}

func styleText(text string, style lipgloss.Style) {
	fmt.Println(style.Render(text))
}
