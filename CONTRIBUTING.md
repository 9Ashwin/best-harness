# Contributing

Start with the smallest surface that matches the change.

- CLI behavior and data shape: `skills/best-harness/scripts/internal/harness/` with focused Go tests.
- Agent guidance: `skills/best-harness/SKILL.md`.
- Public project behavior: `README.md` and `docs/`.

Keep additions portable and executable. Do not add organization-specific task IDs, filesystem paths, session readers, credentials, workflow scoring, or suggestion reports.

Run `go -C skills/best-harness/scripts test ./...` and `go -C skills/best-harness/scripts vet ./...` before opening a pull request.
