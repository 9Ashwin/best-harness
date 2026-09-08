package harness

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const evolutionDirectory = ".best-harness/evolution"

var evolutionKey = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,63}$`)

type EvolutionObservation struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	Task         string `json:"task"`
	Target       string `json:"target"`
	TargetHash   string `json:"target_hash"`
	LessonKey    string `json:"lesson_key"`
	Signal       string `json:"signal"`
	Summary      string `json:"summary"`
	Evidence     string `json:"evidence"`
	EvidenceHash string `json:"evidence_hash"`
	Check        string `json:"check"`
}

type EvolutionCandidate struct {
	Key            string   `json:"key"`
	Target         string   `json:"target"`
	TargetHash     string   `json:"target_hash"`
	LessonKey      string   `json:"lesson_key"`
	State          string   `json:"state"`
	ObservationIDs []string `json:"observation_ids"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type ObserveOptions struct {
	Root, Task, Target, LessonKey, Signal, Summary, Evidence, Check string
}

func ObserveEvolution(options ObserveOptions) (EvolutionObservation, EvolutionCandidate, error) {
	root, err := gitRoot(options.Root)
	if err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	if !evolutionKey.MatchString(options.LessonKey) {
		return EvolutionObservation{}, EvolutionCandidate{}, errors.New("lesson key must use lowercase letters, digits, or hyphens")
	}
	if !validSignal(options.Signal) {
		return EvolutionObservation{}, EvolutionCandidate{}, errors.New("signal must be user-correction, verified-fix, or review-finding")
	}
	if strings.TrimSpace(options.Task) == "" || strings.TrimSpace(options.Summary) == "" || strings.TrimSpace(options.Check) == "" {
		return EvolutionObservation{}, EvolutionCandidate{}, errors.New("task, summary, and check are required")
	}
	target, targetHash, err := evolutionTarget(root, options.Target)
	if err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	evidence, evidenceHash, err := evolutionEvidence(root, options.Evidence)
	if err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	id, err := randomID()
	if err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	observation := EvolutionObservation{ID: id, CreatedAt: now, Task: strings.TrimSpace(options.Task), Target: target, TargetHash: targetHash, LessonKey: options.LessonKey, Signal: options.Signal, Summary: strings.TrimSpace(options.Summary), Evidence: evidence, EvidenceHash: evidenceHash, Check: strings.TrimSpace(options.Check)}
	if err := writeJSON(filepath.Join(root, evolutionDirectory, "observations", id+".json"), observation); err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	observations, err := matchingObservations(root, target, targetHash, options.LessonKey)
	if err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	candidate := EvolutionCandidate{Key: candidateKey(target, targetHash, options.LessonKey), Target: target, TargetHash: targetHash, LessonKey: options.LessonKey, State: "collecting", CreatedAt: now, UpdatedAt: now}
	for _, item := range observations {
		candidate.ObservationIDs = append(candidate.ObservationIDs, item.ID)
	}
	if eligible(observations) {
		candidate.State = "review_pending"
	}
	if err := writeJSON(filepath.Join(root, evolutionDirectory, "candidates", candidate.Key+".json"), candidate); err != nil {
		return EvolutionObservation{}, EvolutionCandidate{}, err
	}
	return observation, candidate, nil
}

func EvolutionStatus(root string) ([]EvolutionCandidate, error) {
	root, err := gitRoot(root)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(root, evolutionDirectory, "candidates"))
	if errors.Is(err, os.ErrNotExist) {
		return []EvolutionCandidate{}, nil
	}
	if err != nil {
		return nil, err
	}
	result := make([]EvolutionCandidate, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var candidate EvolutionCandidate
		if err := readJSON(filepath.Join(root, evolutionDirectory, "candidates", entry.Name()), &candidate); err != nil {
			return nil, err
		}
		result = append(result, candidate)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result, nil
}

func ApplyEvolution(root, key, lesson string) (EvolutionCandidate, error) {
	root, err := gitRoot(root)
	if err != nil {
		return EvolutionCandidate{}, err
	}
	if !evolutionKey.MatchString(key) {
		return EvolutionCandidate{}, errors.New("candidate key is invalid")
	}
	lesson = strings.TrimSpace(lesson)
	if lesson == "" || len(lesson) > 480 || strings.ContainsAny(lesson, "\r\n") {
		return EvolutionCandidate{}, errors.New("lesson must be one line with at most 480 characters")
	}
	path := filepath.Join(root, evolutionDirectory, "candidates", key+".json")
	var candidate EvolutionCandidate
	if err := readJSON(path, &candidate); err != nil {
		return EvolutionCandidate{}, err
	}
	if candidate.State != "review_pending" {
		return EvolutionCandidate{}, fmt.Errorf("candidate is not ready to apply: %s", candidate.State)
	}
	_, hash, err := evolutionTarget(root, candidate.Target)
	if err != nil || hash != candidate.TargetHash {
		return EvolutionCandidate{}, errors.New("candidate target changed; record a new observation")
	}
	for _, id := range candidate.ObservationIDs {
		var observation EvolutionObservation
		if err := readJSON(filepath.Join(root, evolutionDirectory, "observations", id+".json"), &observation); err != nil {
			return EvolutionCandidate{}, err
		}
		_, hash, err := evolutionEvidence(root, observation.Evidence)
		if err != nil || hash != observation.EvidenceHash {
			return EvolutionCandidate{}, errors.New("candidate evidence changed; record a new observation")
		}
	}
	targetPath := filepath.Join(root, filepath.FromSlash(candidate.Target))
	content, err := os.ReadFile(targetPath)
	if err != nil {
		return EvolutionCandidate{}, err
	}
	updated := string(content)
	if !strings.Contains(updated, "## Harness lessons\n") {
		updated = strings.TrimRight(updated, "\n") + "\n\n## Harness lessons\n"
	}
	updated += "\n- " + lesson + "\n"
	if err := atomicWrite(targetPath, []byte(updated), 0644); err != nil {
		return EvolutionCandidate{}, err
	}
	candidate.State = "applied"
	candidate.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if err := writeJSON(path, candidate); err != nil {
		return EvolutionCandidate{}, err
	}
	return candidate, nil
}

func validSignal(value string) bool {
	return value == "user-correction" || value == "verified-fix" || value == "review-finding"
}
func eligible(items []EvolutionObservation) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if item.Signal == "user-correction" {
			return true
		}
	}
	tasks := map[string]bool{}
	for _, item := range items {
		tasks[item.Task] = true
	}
	return len(tasks) >= 2
}

func evolutionTarget(root, target string) (string, string, error) {
	if filepath.IsAbs(target) {
		return "", "", errors.New("evolution target must be repository-relative")
	}
	clean := filepath.ToSlash(filepath.Clean(target))
	if !strings.HasPrefix(clean, ".agents/skills/") || !strings.HasSuffix(clean, ".md") || strings.Contains(clean, "..") {
		return "", "", errors.New("evolution target must be a Skill Markdown file under .agents/skills")
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(clean)))
	if err != nil {
		return "", "", err
	}
	if !strings.HasPrefix(string(data), "---\n") {
		return "", "", errors.New("evolution target must have Skill frontmatter")
	}
	return clean, digest(data), nil
}

func evolutionEvidence(root, evidence string) (string, string, error) {
	if filepath.IsAbs(evidence) {
		return "", "", errors.New("evidence must be repository-relative")
	}
	clean := filepath.ToSlash(filepath.Clean(evidence))
	if clean == "." || strings.Contains(clean, "..") {
		return "", "", errors.New("evidence must remain inside the repository")
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(clean)))
	if err != nil {
		return "", "", err
	}
	return clean, digest(data), nil
}

func matchingObservations(root, target, targetHash, lesson string) ([]EvolutionObservation, error) {
	entries, err := os.ReadDir(filepath.Join(root, evolutionDirectory, "observations"))
	if errors.Is(err, os.ErrNotExist) {
		return []EvolutionObservation{}, nil
	}
	if err != nil {
		return nil, err
	}
	var result []EvolutionObservation
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var item EvolutionObservation
		if err := readJSON(filepath.Join(root, evolutionDirectory, "observations", entry.Name()), &item); err != nil {
			return nil, err
		}
		if item.Target == target && item.TargetHash == targetHash && item.LessonKey == lesson {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func candidateKey(target, targetHash, lesson string) string {
	return digest([]byte(target + "\x00" + targetHash + "\x00" + lesson))[:20]
}
func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func randomID() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return time.Now().UTC().Format("20060102T150405.000000000Z") + "-" + hex.EncodeToString(b), nil
}
func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".state-")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Chmod(name, 0600); err != nil {
		return err
	}
	return os.Rename(name, path)
}
func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".write-")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Chmod(name, mode); err != nil {
		return err
	}
	return os.Rename(name, path)
}
