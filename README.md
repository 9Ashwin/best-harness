<p align="center"><strong>Best Harness</strong></p>

<p align="center">先看证据，再改进编码 Agent 的交付闭环。</p>

[English](README.en.md) · [架构](docs/architecture.md) · [贡献](CONTRIBUTING.md)

Best Harness 是一个轻量、可移植的编码 Agent 工作流工具。它只报告 Git 仓库中可观察的事实：指导文件、工作区变更、暂存路径和可能的验证命令；没有证据的部分明确标为未测量，不会伪装成质量评分或正确性结论。

对于明确的实现、设计、排查和交付任务，它默认创建或复用轻量 Markdown 任务记录，并生成适合 Codex 的标题。需求 ID 是可选项，因此可用于 GitHub Issue、外部需求系统或完全没有编号的个人项目。

## 快速开始

在项目源码目录执行：

```bash
go -C skills/best-harness/scripts run . inspect --staged --format markdown
go -C skills/best-harness/scripts run . task title --title "提升导出可靠性"
go -C skills/best-harness/scripts run . task ensure --title "提升导出可靠性"
```

带可选外部编号、并指定任务目录：

```bash
go -C skills/best-harness/scripts run . task ensure \
  --id "GH-42" \
  --title "提升导出可靠性" \
  --dir planning
```

有 `--id` 时会创建 `planning/gh-42-提升导出可靠性/README.md`；没有编号时默认创建 `tasks/提升导出可靠性/README.md`。同名目录已存在时，`task ensure` 会复用其中的 `README.md` 而不覆盖内容。

## 使用 npx 安装 Skill

使用 `npx skills` 将完整 Skill 包（包括 Go 脚本）安装到当前项目的 Codex：

```bash
npx -y skills@latest add 9Ashwin/best-harness \
  --skill best-harness \
  --agent codex \
  --yes
```

安装器会创建 `.agents/skills/best-harness/`。从安装目录直接运行 CLI：

```bash
go -C .agents/skills/best-harness/scripts run . inspect --staged --format markdown
```

加上 `--global` 可安装到所有 Codex 项目。安装后新开一个 Agent 会话，让 Skill 清单重新加载。

## 它能观察什么

| 范围 | 可观察证据 | 不代表什么 |
| --- | --- | --- |
| 指导 | `AGENTS.md`、`CLAUDE.md` 是否存在 | 不代表 Agent 已遵守 |
| Git 状态 | 工作区和可选暂存路径数量 | 不代表改动质量 |
| 验证 | 从 `go.mod`、`package.json`、`pyproject.toml` 推断的命令 | 不代表命令已经运行或通过 |
| 任务记录 | 可选的本地 Markdown 文件 | 不会创建 Issue，也不要求编号 |

## Codex 插件

仓库包含 Codex 插件清单和 [`$best-harness` Skill](skills/best-harness/SKILL.md)。它没有会话采集器、外部服务或特定 Agent 宿主依赖。

## 开发

```bash
go -C skills/best-harness/scripts test ./...
go -C skills/best-harness/scripts vet ./...
go -C skills/best-harness/scripts run . inspect --format markdown
```

GitHub Actions 会在 Ubuntu、macOS 和 Windows 上运行测试与 vet。

## 致谢

项目的 README 结构和“证据优先”原则参考了 [QoderAI/better-harness](https://github.com/QoderAI/better-harness)。本项目是独立、精简的 Go 实现，不包含其工作流资产或任何个人项目数据。

## 许可证

[MIT](LICENSE)
