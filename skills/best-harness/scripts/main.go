package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/best-harness/internal/harness"
)

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "best-harness:", err)
		if code == 0 {
			code = 2
		}
	}
	os.Exit(code)
}

func run(args []string) (int, error) {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return 0, nil
	}
	switch args[0] {
	case "check":
		return check(args[1:])
	case "verify":
		return verify(args[1:])
	case "task":
		return task(args[1:])
	default:
		return 2, fmt.Errorf("unknown command %q", args[0])
	}
}

func check(args []string) (int, error) {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	root := flags.String("root", ".", "Git checkout to check")
	staged := flags.Bool("staged", false, "check only the staged change")
	if err := flags.Parse(args); err != nil {
		return 2, err
	}
	report, err := harness.Check(*root, *staged)
	if err != nil {
		return 2, err
	}
	printJSON(report)
	if len(report.Errors) != 0 {
		return 1, nil
	}
	return 0, nil
}

func verify(args []string) (int, error) {
	if len(args) == 0 || args[0] != "run" {
		return 2, errors.New("verify requires run")
	}
	flags := flag.NewFlagSet("verify run", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	root := flags.String("root", ".", "Git checkout that owns the verification")
	label := flags.String("label", "", "stable check label")
	output := flags.String("output", "", "receipt path relative to the repository root")
	if err := flags.Parse(args[1:]); err != nil {
		return 2, err
	}
	command := flags.Args()
	if *label == "" || len(command) == 0 {
		return 2, errors.New("verify run requires --label and -- command")
	}
	receipt, code, err := harness.RunVerification(harness.VerificationOptions{
		Root: *root, Label: *label, Command: command, Output: *output,
	})
	if err != nil {
		return 2, err
	}
	printJSON(receipt)
	return code, nil
}

func task(args []string) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("task requires ensure, new, or title")
	}
	switch args[0] {
	case "title":
		flags := flag.NewFlagSet("task title", flag.ContinueOnError)
		flags.SetOutput(os.Stderr)
		id := flags.String("id", "", "optional issue or external reference")
		title := flags.String("title", "", "human-readable task title")
		if err := flags.Parse(args[1:]); err != nil {
			return 2, err
		}
		result, err := harness.CanonicalTitle(*id, *title)
		if err != nil {
			return 2, err
		}
		fmt.Println(result)
		return 0, nil
	case "new":
		return createTask(args[1:], false)
	case "ensure":
		return createTask(args[1:], true)
	default:
		return 2, fmt.Errorf("unknown task command %q", args[0])
	}
}

func createTask(args []string, ensure bool) (int, error) {
	name := "task new"
	if ensure {
		name = "task ensure"
	}
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	root := flags.String("root", ".", "Git checkout that owns the task")
	directory := flags.String("dir", "tasks", "relative task directory")
	id := flags.String("id", "", "optional issue or external reference")
	title := flags.String("title", "", "human-readable task title")
	if err := flags.Parse(args); err != nil {
		return 2, err
	}
	options := harness.TaskOptions{Root: *root, Directory: *directory, ID: *id, Title: *title}
	var (
		result harness.Task
		err    error
	)
	if ensure {
		result, err = harness.EnsureTask(options)
	} else {
		result, err = harness.CreateTask(options)
	}
	if err != nil {
		return 2, err
	}
	relative, err := filepath.Rel(*root, result.Path)
	if err != nil {
		relative = result.Path
	}
	state := "existing"
	if result.Created {
		state = "created"
	}
	fmt.Printf("%s %s\n%s\n", state, result.Title, filepath.ToSlash(relative))
	return 0, nil
}

func printJSON(value any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
}

func printUsage() {
	fmt.Fprint(os.Stdout, `Best Harness: direct task, check, and verification controls.

Usage:
  best-harness check [--root <path>] [--staged]
  best-harness verify run --label <name> [--root <path>] [--output <path>] -- <command> [args...]
  best-harness task title --title <text> [--id <reference>]
  best-harness task ensure --title <text> [--id <reference>] [--root <path>] [--dir <directory>]
  best-harness task new --title <text> [--id <reference>] [--root <path>] [--dir <directory>]
`)
}
