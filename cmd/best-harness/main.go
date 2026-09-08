package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/best-harness/internal/harness"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "best-harness:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}
	switch args[0] {
	case "inspect":
		return inspect(args[1:])
	case "task":
		return task(args[1:])
	case "help", "--help", "-h":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func inspect(args []string) error {
	flags := flag.NewFlagSet("inspect", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	root := flags.String("root", ".", "Git checkout to inspect")
	staged := flags.Bool("staged", false, "include staged path evidence")
	format := flags.String("format", "json", "output format: json or markdown")
	if err := flags.Parse(args); err != nil {
		return err
	}
	report, err := harness.Inspect(*root, *staged)
	if err != nil {
		return err
	}
	switch *format {
	case "json":
		output, err := harness.JSON(report)
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	case "markdown":
		fmt.Print(harness.Markdown(report))
	default:
		return fmt.Errorf("unknown format %q", *format)
	}
	return nil
}

func task(args []string) error {
	if len(args) == 0 {
		return errors.New("task requires new or title")
	}
	switch args[0] {
	case "title":
		flags := flag.NewFlagSet("task title", flag.ContinueOnError)
		flags.SetOutput(os.Stderr)
		id := flags.String("id", "", "optional issue or external reference")
		title := flags.String("title", "", "human-readable task title")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		result, err := harness.CanonicalTitle(*id, *title)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	case "new":
		flags := flag.NewFlagSet("task new", flag.ContinueOnError)
		flags.SetOutput(os.Stderr)
		root := flags.String("root", ".", "Git checkout that owns the task")
		directory := flags.String("dir", "tasks", "relative task directory")
		id := flags.String("id", "", "optional issue or external reference")
		title := flags.String("title", "", "human-readable task title")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		result, err := harness.CreateTask(harness.TaskOptions{Root: *root, Directory: *directory, ID: *id, Title: *title})
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(*root, result.Path)
		if err != nil {
			relative = result.Path
		}
		fmt.Printf("created %s\n%s\n", result.Title, filepath.ToSlash(relative))
		return nil
	default:
		return fmt.Errorf("unknown task command %q", args[0])
	}
}

func printUsage() {
	fmt.Fprint(os.Stdout, `Best Harness: evidence-first checks for coding-agent workflows.

Usage:
  best-harness inspect [--root <path>] [--staged] [--format json|markdown]
  best-harness task title --title <text> [--id <reference>]
  best-harness task new --title <text> [--id <reference>] [--root <path>] [--dir <directory>]
`)
}
