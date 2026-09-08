---
name: best-harness
description: 为编码 Agent 提供可直接执行的任务目录、Git 检查和验证回执。用户请求实现、修改、设计、排查或交付明确任务时使用，默认创建或复用任务记录；交付前执行暂存检查与项目验证命令。不要把它当作工作流评分、建议报告或会话记录分析工具。
---

# Best Harness

Best Harness 是执行 Harness，不做评分或建议报告。它直接创建任务边界、运行确定性检查，并保存实际验证回执。

## 默认任务记录

对明确的实现、修改、设计、排查或交付任务，先自动执行一次 `task ensure`。普通问答、只读解释、单条命令，以及用户明确说“不建任务记录”的情况不创建文件。

- 默认会话标题：`<ID> · <标题>`；无 ID 时为 `<标题>`。
- 默认目录：`tasks/<id>-<标题>/README.md`；无 ID 时为 `tasks/<标题>/README.md`。
- 用户提供的标题、ID、`--dir` 或命名格式优先。

```bash
go -C .agents/skills/best-harness/scripts run . task ensure --title "优化订单导出"
go -C .agents/skills/best-harness/scripts run . task ensure --id "GH-42" --title "优化订单导出"
```

同名目录已存在时，`ensure` 复用 `README.md`，不覆盖正文。

## 交付检查

对准备交付的 Git 变更执行实际检查：

```bash
go -C .agents/skills/best-harness/scripts run . check --staged
```

它执行 `git diff --check` 并拒绝未合并索引项。通过只证明这些 Git 条件，不代表业务正确。

## 运行并记录验证

显式传入项目自己的验证命令。工具会运行它，将退出码、耗时和命令哈希写入 `.best-harness/receipts/`；不保存命令输出或原始参数。

```bash
go -C .agents/skills/best-harness/scripts run . verify run \
  --label unit-tests -- go test ./...
```

不要把建议命令、配置文件存在或任务记录当作通过证据；只有该命令实际运行并返回成功才是验证回执。
