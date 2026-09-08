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

## 记录可复用纠正

只有已经验证的用户纠正、修复或评审发现才记录为演进观察。观察必须带当前 Skill 目标、仓库内证据文件和实际检查；不保存聊天全文或原始日志。

```bash
go -C .agents/skills/best-harness/scripts run . evolve observe \
  --task export-fix \
  --target .agents/skills/best-harness/SKILL.md \
  --lesson-key verify-before-delivery \
  --signal user-correction \
  --summary "交付前必须记录实际验证回执。" \
  --evidence .best-harness/receipts/<receipt>.json \
  --check "go test ./..."
```

`user-correction` 一条已验证观察即可进入 `review_pending`；其他信号需要两个不同任务。候选和脱敏观察保存在 `.best-harness/evolution/`，不改业务代码、提交或推送。

候选进入 `review_pending` 后，先阅读观察和目标 Skill，再由 Agent 归纳一条短的可复用规则，调用：

```bash
go -C .agents/skills/best-harness/scripts run . evolve apply \
  --key <candidate-key> \
  --lesson "运行这类任务前先保存实际验证回执。"
```

`apply` 会重新核对目标和所有证据哈希；任一变化都会拒绝写入。它只向目标 Skill 追加规则，不修改业务代码、提交或推送。
