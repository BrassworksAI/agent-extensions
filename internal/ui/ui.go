package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	Theme Theme
}

func New() *UI {
	return &UI{Theme: GetTheme()}
}

func (u *UI) Choose(header string, options []string) (string, error) {
	if len(options) == 0 {
		return "", errors.New("no options available")
	}

	selected := options[0]
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(header).
				Options(optionsToHuh(options)...).
				Value(&selected),
		),
	).Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(selected), nil
}

func (u *UI) ChooseMulti(header string, options []string) ([]string, error) {
	if len(options) == 0 {
		return []string{}, nil
	}

	selected := make([]string, 0, len(options))
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(header).
				Options(optionsToHuh(options)...).
				Validate(func(value []string) error {
					if len(value) == 0 {
						return errors.New("please select at least one")
					}
					return nil
				}).
				Value(&selected),
		),
	).Run(); err != nil {
		return nil, err
	}

	return selected, nil
}

func (u *UI) Confirm(prompt string) (bool, error) {
	confirmed := false
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Affirmative("Yes").
				Negative("No").
				Value(&confirmed),
		),
	).Run(); err != nil {
		return false, err
	}

	return confirmed, nil
}

func (u *UI) Spin(title string, fn func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- fn()
		close(done)
	}()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Primary))

	model := spinModel{
		spinner: s,
		title:   title,
		done:    done,
	}

	finalModel, err := tea.NewProgram(model, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return err
	}

	result, ok := finalModel.(spinModel)
	if !ok {
		return errors.New("unexpected spinner model type")
	}

	if result.err != nil {
		u.Error("failed")
		return result.err
	}

	u.Success("done")
	return nil
}

func (u *UI) Success(msg string) {
	symbol := lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Success)).Render("✓")
	fmt.Printf("%s %s\n", symbol, msg)
}

func (u *UI) Error(msg string) {
	symbol := lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Error)).Render("✗")
	fmt.Printf("%s %s\n", symbol, msg)
}

func (u *UI) Info(msg string) {
	symbol := lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Primary)).Render("•")
	fmt.Printf("%s %s\n", symbol, msg)
}

func (u *UI) Warn(msg string) {
	symbol := lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Warning)).Render("!")
	fmt.Printf("%s %s\n", symbol, msg)
}

func (u *UI) Header(msg string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Secondary))
	fmt.Printf("\n%s\n", style.Render(msg))
}

func (u *UI) Summary(content string) {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(u.Theme.Muted)).
		Padding(0, 1)

	fmt.Println(style.Render(content))
}

func (u *UI) Title() {
	title := `
 █████╗  ███████╗
██╔══██╗ ██╔════╝
███████║ █████╗  
██╔══██║ ██╔══╝  
██║  ██║ ███████╗
╚═╝  ╚═╝ ╚══════╝`
	subtitle := "Supercharge your AI agents"
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Primary)).Render(title))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(u.Theme.Muted)).Render(subtitle))
	fmt.Println()
}

func optionsToHuh(options []string) []huh.Option[string] {
	out := make([]huh.Option[string], 0, len(options))
	for _, option := range options {
		out = append(out, huh.NewOption(option, option))
	}
	return out
}

type spinDoneMsg struct {
	err error
}

type spinModel struct {
	spinner spinner.Model
	title   string
	done    <-chan error
	err     error
}

func (m spinModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForSpinDone(m.done))
}

func (m spinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case spinDoneMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	return m, nil
}

func (m spinModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.title)
}

func waitForSpinDone(done <-chan error) tea.Cmd {
	return func() tea.Msg {
		return spinDoneMsg{err: <-done}
	}
}
