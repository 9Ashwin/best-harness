<p align="center"><strong>Best Harness</strong></p>

<p align="center">Reusable task, check, and verification controls for coding agents.</p>

[简体中文](README.md) · [Architecture](docs/architecture.md) · [Contributing](CONTRIBUTING.md)

Best Harness is an installable Agent Skill. It does not analyze sessions, score workflows, or turn project files into suggestions. It brings a reusable execution Harness into any Git project: default task records, deterministic Git checks, and verification receipts.

## Three direct actions

| Action | Command | Result |
| --- | --- | --- |
| Establish a task boundary | `task ensure` | Creates or reuses `tasks/<name>/README.md` and prints a Codex-friendly title |
| Check before delivery | `check --staged` | Runs `git diff --check` and checks for unmerged index entries |
| Record verification | `verify run -- <command>` | Runs a project check and stores a sanitized receipt with exit code, duration, and command hash |
| Capture a verified correction | `evolve observe` | Records a reusable Skill-improvement candidate without storing the full session |

Task IDs are optional. With an ID, the title is `ID · title` and the directory is `tasks/<id>-<title>/`; without one, the title and `tasks/<title>/` are used directly. A user-specified title, ID, directory, or naming format wins.

## Install with npx

```bash
npx -y skills@latest add 9Ashwin/best-harness \
  --skill best-harness \
  --agent codex \
  --yes
```

The installer copies the complete Skill and its Go scripts to `.agents/skills/best-harness/`.

## Use in a project

For a concrete task, the agent first creates or reuses a task record:

```bash
go -C .agents/skills/best-harness/scripts run . task ensure \
  --id "GH-42" --title "Improve export reliability"
```

Check the actual staged change before delivery:

```bash
go -C .agents/skills/best-harness/scripts run . check --staged
```

Run and record the project's own verification command:

```bash
go -C .agents/skills/best-harness/scripts run . verify run \
  --label unit-tests -- go test ./...
```

Receipts are stored under `.best-harness/receipts/`. They retain a label, command hash, status, exit code, and duration, never raw command arguments or output.

For a verified repeated correction, `evolve observe` records a candidate under `.best-harness/evolution/`. It accepts only a current Skill Markdown target and repository evidence, and never modifies business code or commits automatically. A `review_pending` candidate is applied with `evolve apply --key <candidate> --lesson "<rule>"`, which rechecks the target and evidence hashes.

## Boundaries

- No requirement ID, issue tracker, task system, or fixed directory is required.
- It does not read agent sessions, business data, credentials, or external services.
- It does not modify business code, commit, push, or publish automatically.
- A Git check or verification receipt proves only the condition that actually ran.

## Development

```bash
go -C skills/best-harness/scripts test ./...
go -C skills/best-harness/scripts vet ./...
go -C skills/best-harness/scripts run . check --staged
```

GitHub Actions runs tests and vet on Ubuntu, macOS, and Windows.

## License

[MIT](LICENSE)
