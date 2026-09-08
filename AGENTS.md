# Best Harness contributor guide

Best Harness is a portable, evidence-first tool for coding-agent workflows. It must run on macOS, Linux, and Windows without shell-specific behavior.

- Keep the core free of organization-specific paths, requirements IDs, credentials, session transcripts, and customer data.
- Treat configured files as evidence that a mechanism exists, never evidence that an agent used it or that a change succeeded.
- Keep missing evidence explicit in CLI output and documentation.
- Use Go standard-library APIs and argv-based `exec.Command`; do not require a shell, `/tmp`, or a package manager for core behavior.
- `task` is opt-in. Its ID and task directory are configurable; never impose an issue tracker, title style, or directory convention.
- Preserve public CLI output contracts. Add or change behavior with focused tests that exercise the observable result.

Before committing, run `go -C skills/best-harness/scripts test ./...`, `go -C skills/best-harness/scripts vet ./...`, and `go -C skills/best-harness/scripts run . inspect --format markdown` from a Git checkout.
