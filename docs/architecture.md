# Architecture

Best Harness is a single portable Skill bundle:

- `skills/best-harness/SKILL.md`: agent-facing guidance;
- `skills/best-harness/scripts/`: the self-contained Go CLI and tests;
- `.codex-plugin/plugin.json`: Codex plugin metadata.

Installing the Skill copies its CLI source alongside the instructions. The installed command uses no shell and only the Go standard library plus Git.

The CLI collects local repository metadata: Git status, staged paths, guidance-file presence, and likely verification commands. It does not read coding-agent transcripts, network services, credentials, or organization-specific files.

A report separates observed evidence from suggested checks. Project owners retain responsibility for choosing, running, and interpreting their own validation commands.
