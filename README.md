<p align="center"><strong>Best Harness</strong></p>

<p align="center">直接复用编码 Agent 的任务、检查与验证闭环。</p>

[English](README.en.md) · [架构](docs/architecture.md) · [贡献](CONTRIBUTING.md)

Best Harness 是一个可安装的 Agent Skill。它不分析会话、不生成工作流评分，也不根据项目文件给建议；它把可复用 Harness 直接带进任何 Git 项目：默认任务目录、确定性 Git 检查和实际验证回执。

## 三个直接动作

| 动作 | 命令 | 结果 |
| --- | --- | --- |
| 建立任务边界 | `task ensure` | 创建或复用 `tasks/<名称>/README.md`，并输出 Codex 会话标题 |
| 交付前检查 | `check --staged` | 运行 `git diff --check` 并检查未合并索引项 |
| 记录验证 | `verify run -- <命令>` | 执行项目检查，写入脱敏退出码、耗时和命令哈希回执 |

任务 ID 可选：有 ID 时标题为 `ID · 标题`、目录为 `tasks/<id>-<标题>/`；无 ID 时直接使用标题和 `tasks/<标题>/`。用户提供的标题、ID、目录或命名格式优先。

## 使用 npx 安装

```bash
npx -y skills@latest add 9Ashwin/best-harness \
  --skill best-harness \
  --agent codex \
  --yes
```

安装器会将完整 Skill 和 Go 脚本复制到 `.agents/skills/best-harness/`。

## 在项目中使用

对明确任务，Agent 默认先创建或复用任务记录：

```bash
go -C .agents/skills/best-harness/scripts run . task ensure \
  --id "GH-42" --title "提升导出可靠性"
```

交付前检查实际暂存内容：

```bash
go -C .agents/skills/best-harness/scripts run . check --staged
```

运行并记录项目自己的验证：

```bash
go -C .agents/skills/best-harness/scripts run . verify run \
  --label unit-tests -- go test ./...
```

验证回执默认保存到 `.best-harness/receipts/`。它只记录标签、命令哈希、状态、退出码和耗时；不保存原始命令或输出。

## 边界

- 不要求需求 ID、Issue、特定任务系统或固定目录。
- 不读取 Agent 会话、业务数据、凭据或外部服务。
- 不自动修改业务代码、提交、推送或发布。
- Git 检查和验证回执只证明实际运行过的条件，不证明业务结果。

## 开发

```bash
go -C skills/best-harness/scripts test ./...
go -C skills/best-harness/scripts vet ./...
go -C skills/best-harness/scripts run . check --staged
```

GitHub Actions 会在 Ubuntu、macOS 和 Windows 上运行测试与 vet。

## 致谢

项目从本地可执行 Harness 的任务、检查与证据边界中提取通用部分；不包含任何工作区规则、业务资料或个人路径。

## 许可证

[MIT](LICENSE)
