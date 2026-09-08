package harness

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectReportsEvidenceAndUnknownValidation(t *testing.T) {
	root := gitFixture(t)
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# Rules\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ngo 1.24.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	report, err := Inspect(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if !hasEvidence(report, "guidance", "present") || !hasEvidence(report, "validation", "not_measured") {
		t.Fatalf("unexpected evidence: %#v", report.Evidence)
	}
	if len(report.NextChecks) != 1 || report.NextChecks[0] != "go test ./..." {
		t.Fatalf("unexpected checks: %#v", report.NextChecks)
	}
}

func TestInspectIncludesStagedEvidence(t *testing.T) {
	root := gitFixture(t)
	write(t, filepath.Join(root, "README.md"), "# changed\n")
	git(t, root, "add", "README.md")
	report, err := Inspect(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if !hasEvidence(report, "staged_changes", "present") {
		t.Fatalf("expected staged evidence: %#v", report.Evidence)
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

func TestCreateTaskUsesConfiguredDirectoryWithoutID(t *testing.T) {
	root := gitFixture(t)
	task, err := CreateTask(TaskOptions{Root: root, Directory: "plans", Title: "Improve export"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(filepath.ToSlash(task.Path), "plans/improve-export/README.md") {
		t.Fatalf("unexpected path: %s", task.Path)
	}
	body, err := os.ReadFile(task.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "# Improve export") {
		t.Fatalf("unexpected task document: %s", body)
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

func hasEvidence(report Inspection, area, state string) bool {
	for _, item := range report.Evidence {
		if item.Area == area && item.State == state {
			return true
		}
	}
	return false
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
