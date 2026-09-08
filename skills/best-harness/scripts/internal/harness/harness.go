// Package harness collects portable, evidence-bounded project observations.
package harness

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Evidence struct {
	Area   string `json:"area"`
	State  string `json:"state"`
	Detail string `json:"detail"`
}

type Inspection struct {
	Root       string     `json:"root"`
	Generated  time.Time  `json:"generated_at"`
	Staged     bool       `json:"staged"`
	Evidence   []Evidence `json:"evidence"`
	NextChecks []string   `json:"next_checks"`
}

func Inspect(start string, staged bool) (Inspection, error) {
	root, err := gitRoot(start)
	if err != nil {
		return Inspection{}, err
	}
	result := Inspection{Root: root, Generated: time.Now().UTC(), Staged: staged}
	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			result.Evidence = append(result.Evidence, Evidence{"guidance", "present", name + " is available to the agent."})
		} else if errors.Is(err, os.ErrNotExist) {
			result.Evidence = append(result.Evidence, Evidence{"guidance", "missing", name + " was not found; no claim is made about task guidance."})
		} else {
			return Inspection{}, fmt.Errorf("inspect %s: %w", name, err)
		}
	}

	changes, err := gitLines(root, "status", "--porcelain=v1")
	if err != nil {
		return Inspection{}, err
	}
	state := "clean"
	detail := "Git worktree has no reported changes."
	if len(changes) > 0 {
		state = "changed"
		detail = fmt.Sprintf("Git reports %d changed path(s).", len(changes))
	}
	result.Evidence = append(result.Evidence, Evidence{"worktree", state, detail})

	if staged {
		paths, err := gitLines(root, "diff", "--cached", "--name-only")
		if err != nil {
			return Inspection{}, err
		}
		state, detail = "empty", "No staged paths were found."
		if len(paths) > 0 {
			state = "present"
			detail = fmt.Sprintf("Git index contains %d staged path(s).", len(paths))
		}
		result.Evidence = append(result.Evidence, Evidence{"staged_changes", state, detail})
	}

	result.NextChecks = detectedChecks(root)
	result.Evidence = append(result.Evidence, Evidence{
		Area:   "validation",
		State:  "not_measured",
		Detail: "Best Harness discovers likely checks but does not claim they ran or passed.",
	})
	return result, nil
}

func gitRoot(start string) (string, error) {
	if start == "" {
		start = "."
	}
	absolute, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	output, err := run(absolute, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("best-harness requires a Git checkout: %w", err)
	}
	return filepath.Clean(strings.TrimSpace(output)), nil
}

func gitLines(root string, args ...string) ([]string, error) {
	output, err := run(root, "git", args...)
	if err != nil {
		return nil, err
	}
	return nonEmptyLines(output), nil
}

func run(dir, name string, args ...string) (string, error) {
	command := exec.Command(name, args...)
	command.Dir = dir
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return string(output), nil
}

func nonEmptyLines(value string) []string {
	var lines []string
	for _, line := range strings.Split(value, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func detectedChecks(root string) []string {
	checks := make([]string, 0, 3)
	if exists(filepath.Join(root, "go.mod")) {
		checks = append(checks, "go test ./...")
	}
	if exists(filepath.Join(root, "package.json")) {
		checks = append(checks, "npm test")
	}
	if exists(filepath.Join(root, "pyproject.toml")) {
		checks = append(checks, "python -m pytest")
	}
	if len(checks) == 0 {
		return []string{"Choose and run the repository's documented verification command."}
	}
	return checks
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
