<p align="center"><strong>Best Harness</strong></p>

<p align="center">Evidence first. Better coding-agent delivery loops.</p>

[简体中文](README.zh-CN.md) · [Architecture](docs/architecture.md) · [Contributing](CONTRIBUTING.md)

Best Harness is a small, portable toolkit for inspecting the evidence around a coding-agent task. It reports what a Git repository actually exposes—guidance files, worktree changes, staged paths, and likely verification commands—without turning missing evidence into a score or a claim of correctness.

It also provides an optional task helper for a Codex-friendly title and a lightweight Markdown task record. IDs are optional, so it works with GitHub issues, external requirement systems, or no tracker at all.

## Quick start

From a source checkout:

```bash
go run ./cmd/best-harness inspect --staged --format markdown
go run ./cmd/best-harness task title --title "Improve export reliability"
go run ./cmd/best-harness task new --title "Improve export reliability"
```

With an optional external reference and a custom task directory:

```bash
go run ./cmd/best-harness task new \
  --id "GH-42" \
  --title "Improve export reliability" \
  --dir planning
```

The task helper creates `planning/gh-42-improve-export-reliability/README.md`. Without `--id`, it creates `tasks/improve-export-reliability/README.md`.

## What it observes

| Area | Evidence | Boundary |
| --- | --- | --- |
| Guidance | Presence of `AGENTS.md` and `CLAUDE.md` | Does not prove the agent used them |
| Git state | Worktree and optional staged-path counts | Does not assess change quality |
| Verification | Likely commands inferred from `go.mod`, `package.json`, or `pyproject.toml` | Does not run or claim those checks passed |
| Task record | An optional local Markdown file | Does not create an issue or require an ID |

## Codex plugin

The repository includes a Codex plugin manifest and the `$best-harness` Skill under [`skills/best-harness`](skills/best-harness/SKILL.md). It intentionally has no host adapter, session collector, or external service dependency.

## Development

```bash
go test ./...
go vet ./...
go run ./cmd/best-harness inspect --format markdown
```

The GitHub Actions matrix runs tests and vet on Ubuntu, macOS, and Windows.

## Inspiration

The project structure and evidence-first stance were informed by [QoderAI/better-harness](https://github.com/QoderAI/better-harness), while this implementation is a separate, smaller Go project with no copied workflow assets or project data.

## License

[MIT](LICENSE)
