# Best Harness contributor guide

Best Harness is a portable, evidence-first tool for coding-agent workflows. It must run on macOS, Linux, and Windows without shell-specific behavior.

- Keep the core free of organization-specific paths, requirements IDs, credentials, session transcripts, and customer data.
- Keep direct controls deterministic: task records, Git checks, and verification receipts.
- A passing receipt proves only the command that actually ran; do not add workflow scoring or suggestion reports.
- Use Go standard-library APIs and argv-based `exec.Command`; do not require a shell, `/tmp`, or a package manager for core behavior.
- For a concrete implementation, design, investigation, or delivery task, the Skill defaults to idempotent `task ensure`. Its ID and task directory remain configurable; never impose an issue tracker, title style, or directory convention.
- Preserve public CLI output contracts. Add or change behavior with focused tests that exercise the observable result.

Before committing, run `go -C skills/best-harness/scripts test ./...`, `go -C skills/best-harness/scripts vet ./...`, and `go -C skills/best-harness/scripts run . check --staged` from a Git checkout.
