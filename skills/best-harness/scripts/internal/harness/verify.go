package harness

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var verificationLabel = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

type VerificationOptions struct {
	Root    string
	Label   string
	Command []string
	Output  string
}

type VerificationReceipt struct {
	Label       string `json:"label"`
	CommandHash string `json:"command_hash"`
	Status      string `json:"status"`
	ExitCode    int    `json:"exit_code"`
	StartedAt   string `json:"started_at"`
	ElapsedMS   int64  `json:"elapsed_ms"`
	Receipt     string `json:"receipt"`
}

func RunVerification(options VerificationOptions) (VerificationReceipt, int, error) {
	root, err := gitRoot(options.Root)
	if err != nil {
		return VerificationReceipt{}, 2, err
	}
	if !verificationLabel.MatchString(options.Label) {
		return VerificationReceipt{}, 2, errors.New("label must use lowercase letters, digits, underscores, or hyphens")
	}
	if len(options.Command) == 0 || strings.TrimSpace(options.Command[0]) == "" {
		return VerificationReceipt{}, 2, errors.New("verification command is required")
	}
	receiptPath, err := verificationPath(root, options.Label, options.Output)
	if err != nil {
		return VerificationReceipt{}, 2, err
	}
	started := time.Now().UTC()
	command := exec.Command(options.Command[0], options.Command[1:]...)
	command.Dir = root
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	code := 0
	status := "passed"
	if err != nil {
		status = "failed"
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			code = exitError.ExitCode()
		} else {
			code = 127
			status = "could_not_start"
		}
	}
	hash := sha256.Sum256([]byte(strings.Join(options.Command, "\x00")))
	receipt := VerificationReceipt{
		Label: options.Label, CommandHash: hex.EncodeToString(hash[:]), Status: status, ExitCode: code,
		StartedAt: started.Format(time.RFC3339), ElapsedMS: time.Since(started).Milliseconds(),
		Receipt: filepath.ToSlash(relativeOrAbsolute(root, receiptPath)),
	}
	if err := writeReceipt(receiptPath, receipt); err != nil {
		return VerificationReceipt{}, 2, err
	}
	return receipt, code, nil
}

func verificationPath(root, label, output string) (string, error) {
	if output == "" {
		return filepath.Join(root, ".best-harness", "receipts", time.Now().UTC().Format("20060102T150405.000000000Z")+"-"+label+".json"), nil
	}
	if filepath.IsAbs(output) || strings.HasPrefix(filepath.Clean(output), "..") {
		return "", errors.New("receipt output must be relative to the repository root")
	}
	return filepath.Join(root, output), nil
}

func relativeOrAbsolute(root, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return relative
}

func writeReceipt(path string, receipt VerificationReceipt) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".receipt-")
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
