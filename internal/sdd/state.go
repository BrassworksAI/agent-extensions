package sdd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Lane string

const (
	LaneFull Lane = "full"
	LaneVibe Lane = "vibe"
	LaneBug  Lane = "bug"
)

type PhaseStatus string

const (
	StatusInProgress PhaseStatus = "in_progress"
	StatusComplete   PhaseStatus = "complete"
	StatusBlocked    PhaseStatus = "blocked"
)

type State struct {
	Change  ChangeInfo  `toml:"change"`
	Phase   PhaseInfo   `toml:"phase"`
	Pending PendingInfo `toml:"pending"`
	Notes   NotesInfo   `toml:"notes"`
}

type ChangeInfo struct {
	Name      string    `toml:"name"`
	Lane      Lane      `toml:"lane"`
	CreatedAt time.Time `toml:"created_at"`
}

type PhaseInfo struct {
	Current string      `toml:"current"`
	Status  PhaseStatus `toml:"status"`
}

type PendingInfo struct {
	Items []string `toml:"items"`
}

type NotesInfo struct {
	Content string `toml:"content"`
}

var FullLanePhases = []string{
	"proposal", "specs", "discovery", "tasks", "plan", "implement", "reconcile", "finish",
}

var VibeLanePhases = []string{
	"context", "plan", "implement", "reconcile", "finish",
}

var BugLanePhases = []string{
	"triage", "plan", "implement", "reconcile", "finish",
}

func PhasesForLane(lane Lane) []string {
	switch lane {
	case LaneFull:
		return FullLanePhases
	case LaneVibe:
		return VibeLanePhases
	case LaneBug:
		return BugLanePhases
	default:
		return FullLanePhases
	}
}

func NewState(name string, lane Lane) *State {
	phases := PhasesForLane(lane)
	return &State{
		Change: ChangeInfo{
			Name:      name,
			Lane:      lane,
			CreatedAt: time.Now(),
		},
		Phase: PhaseInfo{
			Current: phases[0],
			Status:  StatusInProgress,
		},
		Pending: PendingInfo{Items: []string{}},
		Notes:   NotesInfo{Content: ""},
	}
}

func LoadState(changeDir string) (*State, error) {
	path := filepath.Join(changeDir, "state.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state.toml: %w", err)
	}

	var state State
	if err := toml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state.toml: %w", err)
	}

	return &state, nil
}

func (s *State) Save(changeDir string) error {
	path := filepath.Join(changeDir, "state.toml")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating state.toml: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(s); err != nil {
		return fmt.Errorf("encoding state.toml: %w", err)
	}

	return nil
}

func (s *State) PhaseIndex() int {
	phases := PhasesForLane(s.Change.Lane)
	for i, p := range phases {
		if p == s.Phase.Current {
			return i
		}
	}
	return -1
}

func (s *State) NextPhase() (string, bool) {
	phases := PhasesForLane(s.Change.Lane)
	idx := s.PhaseIndex()
	if idx < 0 || idx >= len(phases)-1 {
		return "", false
	}
	return phases[idx+1], true
}

func (s *State) PrevPhase() (string, bool) {
	phases := PhasesForLane(s.Change.Lane)
	idx := s.PhaseIndex()
	if idx <= 0 {
		return "", false
	}
	return phases[idx-1], true
}

func (s *State) SetPhase(phase string) error {
	phases := PhasesForLane(s.Change.Lane)
	for _, p := range phases {
		if p == phase {
			s.Phase.Current = phase
			s.Phase.Status = StatusInProgress
			return nil
		}
	}
	return fmt.Errorf("invalid phase %q for lane %q", phase, s.Change.Lane)
}

func (s *State) CompletePhase() {
	s.Phase.Status = StatusComplete
	s.Notes.Content = ""
}

func (s *State) AddPending(item string) {
	s.Pending.Items = append(s.Pending.Items, item)
}

func (s *State) ClearPending(index int) error {
	if index < 0 || index >= len(s.Pending.Items) {
		return fmt.Errorf("invalid pending index: %d", index)
	}
	s.Pending.Items = append(s.Pending.Items[:index], s.Pending.Items[index+1:]...)
	return nil
}

func (s *State) SetNotes(content string) {
	s.Notes.Content = content
}

func HumanizeName(name string) string {
	parts := strings.Split(strings.TrimSpace(name), "-")
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		words = append(words, strings.ToLower(part))
	}
	if len(words) == 0 {
		return name
	}
	sentence := strings.Join(words, " ")
	return strings.ToUpper(sentence[:1]) + sentence[1:]
}

func LaneLabel(lane Lane) string {
	switch lane {
	case LaneFull:
		return "Full"
	case LaneVibe:
		return "Vibe"
	case LaneBug:
		return "Bug"
	default:
		return string(lane)
	}
}
