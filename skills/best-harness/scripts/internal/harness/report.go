package harness

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSON(report Inspection) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

func Markdown(report Inspection) string {
	lines := []string{
		"# Best Harness inspection",
		"",
		"- Repository: `" + report.Root + "`",
		"- Generated: `" + report.Generated.Format("2006-01-02T15:04:05Z") + "`",
		"",
		"## Evidence",
		"",
	}
	for _, item := range report.Evidence {
		lines = append(lines, fmt.Sprintf("- **%s** — `%s`: %s", item.Area, item.State, item.Detail))
	}
	lines = append(lines, "", "## Suggested verification", "")
	for _, check := range report.NextChecks {
		lines = append(lines, "- `"+check+"`")
	}
	lines = append(lines, "", "## Boundary", "", "This report observes repository state. It does not prove that an agent followed guidance, that suggested checks ran, or that a workflow improved.", "")
	return strings.Join(lines, "\n")
}
