# Architecture

Best Harness has three small, independent surfaces:

- `cmd/best-harness`: a portable CLI with no shell dependency;
- `internal/harness`: Git evidence collection, rendering, and optional task-record creation;
- `skills/best-harness`: an agent-facing wrapper that routes to the CLI without inventing evidence.

The CLI collects only local repository metadata: Git status, staged paths, guidance-file presence, and likely verification commands inferred from standard project files. It does not read coding-agent transcripts, network services, credentials, or organization-specific files.

A report separates observed evidence from suggested checks. Project owners retain responsibility for choosing, running, and interpreting their own validation commands.
