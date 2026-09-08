---
name: best-harness
description: 检查 Git 仓库中可观察的编码 Agent 交付证据，并按需创建不绑定需求编号的任务记录。用户要求检查交付证据、查看 Agent 指导文件或暂存改动、生成 Codex 会话标题、创建轻量任务记录时使用。没有运行项目自身检查时，不得用本 Skill 宣称代码正确。
---

# Best Harness

Best Harness 把仓库中可观察的状态整理为有边界的证据：哪些文件存在、哪些改动已暂存、哪些检查只是建议、哪些结果尚未测量。

## 交付前检查

在目标 Git 仓库根目录运行：

```bash
go -C .agents/skills/best-harness/scripts run . inspect --staged --format markdown
```

报告会列出指导文件、工作区状态、暂存路径和可能的验证命令。仍须单独运行项目实际需要的检查；工具建议的命令不代表已经运行或通过。

## 命名或创建任务

只需要生成适合 Codex 的会话标题时，用 `task title`，不会创建文件：

```bash
go -C .agents/skills/best-harness/scripts run . task title --title "优化订单导出"
go -C .agents/skills/best-harness/scripts run . task title --id "GH-42" --title "优化订单导出"
```

仅在用户需要持久化本地任务记录时，用 `task new`。`--id` 是可选项，`--dir` 可指定任意仓库相对目录：

```bash
go -C .agents/skills/best-harness/scripts run . task new --title "优化订单导出"
go -C .agents/skills/best-harness/scripts run . task new --id "GH-42" --title "优化订单导出" --dir planning
```

命令只写入一份 Markdown 任务文件；不会创建 Issue、发送消息、读取会话记录，也不要求使用特定任务系统。

## 证据边界

`AGENTS.md`、任务记录、干净 Git 状态和生成报告都不能证明实际行为。没有执行、验证或结果证据时，必须明确标为“未测量”。
