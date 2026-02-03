package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type UI struct {
	Theme Theme
}

func New() *UI {
	return &UI{Theme: GetTheme()}
}

func (u *UI) Choose(header string, options []string) (string, error) {
	args := []string{"choose", "--header", header}
	args = append(args, "--cursor.foreground", u.Theme.Primary)
	args = append(args, "--header.foreground", u.Theme.Secondary)
	args = append(args, options...)

	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (u *UI) ChooseMulti(header string, options []string) ([]string, error) {
	if len(options) == 0 {
		return []string{}, nil
	}

	isRetry := false
	for {
		// Clear previous warning on retry
		if isRetry {
			u.ClearLines(1)
		}

		args := []string{"choose", "--no-limit", "--header", header}
		args = append(args, "--cursor.foreground", u.Theme.Primary)
		args = append(args, "--selected.foreground", u.Theme.Success)
		args = append(args, "--header.foreground", u.Theme.Secondary)
		args = append(args, options...)

		cmd := exec.Command("gum", args...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		result := strings.TrimSpace(string(out))
		if result == "" {
			u.Warn("Please select at least one (Space to toggle, Enter to confirm)")
			isRetry = true
			continue
		}

		return strings.Split(result, "\n"), nil
	}
}

func (u *UI) Confirm(prompt string) (bool, error) {
	args := []string{"confirm", prompt}
	args = append(args, "--affirmative", "Yes")
	args = append(args, "--negative", "No")
	args = append(args, "--prompt.foreground", u.Theme.Primary)
	args = append(args, "--prompt.margin", "0")
	args = append(args, "--selected.margin", "0")
	args = append(args, "--unselected.margin", "0")

	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

func (u *UI) Spin(title string, fn func() error) error {
	// TODO: Implement gum spin subprocess
	// For now, just run the function directly with a simple indicator
	// For now, just run the function directly with a simple indicator
	fmt.Printf("%s... ", title)
	err := fn()
	if err != nil {
		u.Error("failed")
		return err
	}
	u.Success("done")
	return nil
}

func (u *UI) Success(msg string) {
	fmt.Printf("\033[38;2;%sm✓\033[0m %s\n", hexToRGB(u.Theme.Success), msg)
}

func (u *UI) Error(msg string) {
	fmt.Printf("\033[38;2;%sm✗\033[0m %s\n", hexToRGB(u.Theme.Error), msg)
}

func (u *UI) Info(msg string) {
	fmt.Printf("\033[38;2;%sm•\033[0m %s\n", hexToRGB(u.Theme.Primary), msg)
}

func (u *UI) Warn(msg string) {
	fmt.Printf("\033[38;2;%sm!\033[0m %s\n", hexToRGB(u.Theme.Warning), msg)
}

func (u *UI) Header(msg string) {
	fmt.Printf("\n\033[38;2;%sm%s\033[0m\n", hexToRGB(u.Theme.Secondary), msg)
}

func (u *UI) ClearLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[1A\033[2K")
	}
}

func (u *UI) Summary(content string) {
	args := []string{"style",
		"--border", "rounded",
		"--padding", "0 1",
		"--margin", "0 0 0 0",
		"--border-foreground", u.Theme.Muted,
		content,
	}

	cmd := exec.Command("gum", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
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
	fmt.Printf("\033[38;2;%sm%s\033[0m\n", hexToRGB(u.Theme.Primary), title)
	fmt.Printf("\033[38;2;%sm%s\033[0m\n\n", hexToRGB(u.Theme.Muted), subtitle)
}

func (u *UI) Section(title string) {
	// Minimal section - just the title with muted color
	fmt.Printf("\033[38;2;%sm%s\033[0m\n", hexToRGB(u.Theme.Muted), title)
}

func hexToRGB(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d;%d;%d", r, g, b)
}
