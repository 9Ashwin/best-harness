package harness

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CheckReport struct {
	Root         string   `json:"root"`
	Mode         string   `json:"mode"`
	CheckedAt    string   `json:"checked_at"`
	Changed      int      `json:"changed_paths"`
	Errors       []string `json:"errors"`
	Measurements []string `json:"measurements"`
}

func Check(start string, staged bool) (CheckReport, error) {
	root, err := gitRoot(start)
	if err != nil {
		return CheckReport{}, err
	}
	mode := "worktree"
	basis := "HEAD"
	if staged {
		mode, basis = "staged", "--cached"
	}
	report := CheckReport{Root: root, Mode: mode, CheckedAt: time.Now().UTC().Format(time.RFC3339), Errors: []string{}, Measurements: []string{}}
	changed, err := gitLines(root, "diff", basis, "--name-only")
	if err != nil {
		return CheckReport{}, err
	}
	report.Changed = len(changed)
	report.Measurements = append(report.Measurements, fmt.Sprintf("%d changed path(s) in %s scope", report.Changed, mode))
	if _, err := gitOutput(root, "diff", basis, "--check"); err != nil {
		report.Errors = append(report.Errors, "git diff --check failed")
	}
	unmerged, err := gitLines(root, "ls-files", "-u")
	if err != nil {
		return CheckReport{}, err
	}
	if len(unmerged) > 0 {
		report.Errors = append(report.Errors, "unmerged index entries present")
	}
	return report, nil
}

func gitRoot(start string) (string, error) {
	absolute, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	output, err := gitOutput(absolute, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("best-harness requires a Git checkout: %w", err)
	}
	return filepath.Clean(strings.TrimSpace(output)), nil
}

func gitLines(root string, args ...string) ([]string, error) {
	output, err := gitOutput(root, args...)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func gitOutput(root string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return string(output), nil
}
