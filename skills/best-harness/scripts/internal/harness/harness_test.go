package harness

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckReportsGitDiffErrors(t *testing.T) {
	root := gitFixture(t)
	write(t, filepath.Join(root, "README.md"), "line with space \n")
	report, err := Check(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Errors) != 1 || report.Errors[0] != "git diff --check failed" {
		t.Fatalf("unexpected report: %#v", report)
	}
}

func TestRunVerificationWritesSanitizedReceipt(t *testing.T) {
	root := gitFixture(t)
	receipt, code, err := RunVerification(VerificationOptions{Root: root, Label: "git-status", Command: []string{"git", "status", "--short"}})
	if err != nil || code != 0 || receipt.Status != "passed" {
		t.Fatalf("receipt=%#v code=%d err=%v", receipt, code, err)
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(receipt.Receipt)))
	if err != nil {
		t.Fatal(err)
	}
	var stored VerificationReceipt
	if err := json.Unmarshal(data, &stored); err != nil {
		t.Fatal(err)
	}
	if stored.CommandHash == "" || strings.Contains(string(data), "git status") {
		t.Fatalf("receipt must store the hash, not raw command: %s", data)
	}
}

func TestCanonicalTitleAllowsOptionalReference(t *testing.T) {
	withoutID, err := CanonicalTitle("", "Fix export")
	if err != nil || withoutID != "Fix export" {
		t.Fatalf("got %q, %v", withoutID, err)
	}
	withID, err := CanonicalTitle("GH-42", "Fix export")
	if err != nil || withID != "GH-42 · Fix export" {
		t.Fatalf("got %q, %v", withID, err)
	}
}

func TestEnsureTaskCreatesThenReusesDefaultTask(t *testing.T) {
	root := gitFixture(t)
	options := TaskOptions{Root: root, Title: "Improve export"}
	first, err := EnsureTask(options)
	if err != nil || !first.Created {
		t.Fatalf("first ensure: %#v, %v", first, err)
	}
	second, err := EnsureTask(options)
	if err != nil || second.Created || second.Path != first.Path {
		t.Fatalf("second ensure: %#v, %v", second, err)
	}
}

func TestCreateTaskPrefixesOptionalID(t *testing.T) {
	root := gitFixture(t)
	task, err := CreateTask(TaskOptions{Root: root, ID: "REQ-7", Title: "调整订单导出"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(filepath.ToSlash(task.Path), "tasks/req-7-调整订单导出/README.md") {
		t.Fatalf("unexpected path: %s", task.Path)
	}
}

func TestObserveEvolutionUsesEvidenceGate(t *testing.T) {
	root := gitFixture(t)
	write(t, filepath.Join(root, ".agents", "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: Demo.\n---\n\n# Demo\n")
	write(t, filepath.Join(root, "evidence.md"), "verified evidence\n")
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "demo"), 0755); err != nil {
		t.Fatal(err)
	}
	observation, candidate, err := ObserveEvolution(ObserveOptions{
		Root: root, Task: "task-one", Target: ".agents/skills/demo/SKILL.md", LessonKey: "check-first",
		Signal: "user-correction", Summary: "A verified correction requires the shared check.", Evidence: "evidence.md", Check: "go test ./...",
	})
	if err != nil || observation.ID == "" || candidate.State != "review_pending" {
		t.Fatalf("observation=%#v candidate=%#v err=%v", observation, candidate, err)
	}
	status, err := EvolutionStatus(root)
	if err != nil || len(status) != 1 || status[0].Key != candidate.Key {
		t.Fatalf("status=%#v err=%v", status, err)
	}
	applied, err := ApplyEvolution(root, candidate.Key, "运行这类任务前先保存实际验证回执。")
	if err != nil || applied.State != "applied" {
		t.Fatalf("applied=%#v err=%v", applied, err)
	}
	target, err := os.ReadFile(filepath.Join(root, ".agents", "skills", "demo", "SKILL.md"))
	if err != nil || !strings.Contains(string(target), "## Harness lessons") {
		t.Fatalf("target=%s err=%v", target, err)
	}
}

func gitFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.email", "fixture@example.invalid")
	git(t, root, "config", "user.name", "Fixture")
	write(t, filepath.Join(root, "README.md"), "# fixture\n")
	git(t, root, "add", "README.md")
	git(t, root, "commit", "-m", "fixture")
	return root
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}
