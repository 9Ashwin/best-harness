# Architecture

Best Harness is a self-contained execution bundle:

- `skills/best-harness/SKILL.md`: default task lifecycle and command routing;
- `skills/best-harness/scripts/`: Go CLI, checks, task records, verification receipts, and tests;
- `.codex-plugin/plugin.json`: Codex plugin metadata.

The CLI has three direct controls: idempotent task records, deterministic Git checks, and explicit verification commands with receipts. Runtime receipts live under `.best-harness/` in the target repository and are intentionally excluded from the Skill source tree.

No session collector, scoring model, workflow analyzer, project-specific convention, or external service is part of the core.
