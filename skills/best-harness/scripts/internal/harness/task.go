package harness

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type TaskOptions struct {
	Root      string
	Directory string
	ID        string
	Title     string
}

type Task struct {
	Title   string
	Path    string
	Created bool
}

func CanonicalTitle(id, title string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", errors.New("task title is required")
	}
	if containsPathSeparator(title) {
		return "", errors.New("task title cannot contain a path separator")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return title, nil
	}
	if containsPathSeparator(id) {
		return "", errors.New("task id cannot contain a path separator")
	}
	return id + " · " + title, nil
}

func CreateTask(options TaskOptions) (Task, error) {
	return writeTask(options, false)
}

// EnsureTask returns an existing matching task record or creates it once.
func EnsureTask(options TaskOptions) (Task, error) {
	return writeTask(options, true)
}

func writeTask(options TaskOptions, allowExisting bool) (Task, error) {
	title, err := CanonicalTitle(options.ID, options.Title)
	if err != nil {
		return Task{}, err
	}
	root, err := gitRoot(options.Root)
	if err != nil {
		return Task{}, err
	}
	directory := options.Directory
	if directory == "" {
		directory = "tasks"
	}
	if filepath.IsAbs(directory) || strings.HasPrefix(filepath.Clean(directory), "..") {
		return Task{}, errors.New("task directory must be relative to the repository root")
	}
	name := slug(options.Title)
	if options.ID != "" {
		name = slug(options.ID) + "-" + name
	}
	path := filepath.Join(root, directory, name)
	if _, err := os.Stat(path); err == nil {
		if allowExisting {
			file := filepath.Join(path, "README.md")
			if _, fileErr := os.Stat(file); fileErr == nil {
				return Task{Title: title, Path: file, Created: false}, nil
			} else if !errors.Is(fileErr, os.ErrNotExist) {
				return Task{}, fileErr
			}
			return Task{}, fmt.Errorf("existing task directory has no README.md: %s", path)
		}
		return Task{}, fmt.Errorf("task directory already exists: %s", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Task{}, err
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return Task{}, err
	}
	content := taskDocument(title, options.ID, time.Now().UTC())
	file := filepath.Join(path, "README.md")
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		return Task{}, err
	}
	return Task{Title: title, Path: file, Created: true}, nil
}

func taskDocument(title, id string, created time.Time) string {
	lines := []string{
		"# " + title,
		"",
		"- Status: draft",
		"- Created: " + created.Format("2006-01-02"),
	}
	if id != "" {
		lines = append(lines, "- Reference: "+id)
	}
	lines = append(lines,
		"",
		"## Goal",
		"",
		"Describe the intended user or system outcome.",
		"",
		"## Acceptance",
		"",
		"- [ ] Define an observable completion condition.",
		"",
		"## Evidence",
		"",
		"Record links to relevant code, tests, reviews, or delivery receipts.",
		"",
	)
	return strings.Join(lines, "\n")
}

func slug(value string) string {
	var out []rune
	dash := false
	for _, char := range strings.TrimSpace(strings.ToLower(value)) {
		switch {
		case unicode.IsLetter(char), unicode.IsDigit(char):
			out = append(out, char)
			dash = false
		case unicode.IsSpace(char), char == '-', char == '_':
			if len(out) > 0 && !dash {
				out = append(out, '-')
				dash = true
			}
		}
	}
	result := strings.Trim(string(out), "-")
	if result == "" {
		return "task"
	}
	return result
}

func containsPathSeparator(value string) bool {
	return strings.ContainsAny(value, "/\\") || strings.ContainsRune(value, 0)
}
