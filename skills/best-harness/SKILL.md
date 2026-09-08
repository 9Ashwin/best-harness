---
name: best-harness
description: Inspect a Git repository's observable coding-agent workflow evidence and create an optional task directory with a portable title. Use when the user asks to assess delivery evidence, inspect agent guidance or staged changes, establish a task title, or create a lightweight task record. Do not use it to claim code is correct without running the project's own checks.
---

# Best Harness

Best Harness turns observable repository state into bounded evidence. It distinguishes files that exist, changes that are present, checks that are only suggested, and checks that actually ran.

## Inspect before delivery

Run from the target repository:

```bash
go run github.com/9Ashwin/best-harness/cmd/best-harness@latest inspect --staged --format markdown
```

Use the report to identify guidance files, worktree state, staged paths, and likely verification commands. Run the relevant project checks separately; a suggested command is not a passing result.

## Name or create a task

Use `task title` to produce a Codex-friendly title without creating files:

```bash
best-harness task title --title "Improve order export"
best-harness task title --id "GH-42" --title "Improve order export"
```

Use `task new` only when the user wants a persistent local task record. `--id` is optional, and `--dir` can be any repository-relative directory:

```bash
best-harness task new --title "Improve order export"
best-harness task new --id "GH-42" --title "Improve order export" --dir planning
```

The command writes one Markdown task file. It never creates issues, sends messages, reads session transcripts, or assumes a specific tracker.

## Evidence boundary

Do not treat `AGENTS.md`, a task record, a clean Git state, or a generated report as proof of behavior. Keep absent execution, verification, and outcome evidence marked as not measured.
