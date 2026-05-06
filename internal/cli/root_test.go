package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_NoArgs_ShowsHelpWithoutError(t *testing.T) {
	err := Run([]string{})
	if err != nil {
		t.Fatalf("Run([]string{}) returned error: %v", err)
	}
}

func TestRun_UnknownCommand_ReturnsUnknownCommandError(t *testing.T) {
	err := Run([]string{"unknown"})
	if err == nil {
		t.Fatal("Run([]string{\"unknown\"}) returned nil error")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("expected error to contain %q, got %q", "unknown command", err.Error())
	}
}

func TestRun_ProfileHelp_WithoutError(t *testing.T) {
	err := Run([]string{"profile", "help"})
	if err != nil {
		t.Fatalf("Run([]string{\"profile\", \"help\"}) returned error: %v", err)
	}
}

func TestRun_ProfileCommands_RoutedCorrectly(t *testing.T) {
	tmpDir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() failed: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Chdir(%q) failed: %v", tmpDir, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	if err := Run([]string{"profile", "create", "--name=test", "--user=alice", "--project=alpha", "--force"}); err != nil {
		t.Fatalf("profile create failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "test.yaml")); err != nil {
		t.Fatalf("expected created file test.yaml, stat error: %v", err)
	}

	if err := Run([]string{"profile", "get", "--name=test"}); err != nil {
		t.Fatalf("profile get failed: %v", err)
	}

	if err := Run([]string{"profile", "list"}); err != nil {
		t.Fatalf("profile list failed: %v", err)
	}

	if err := Run([]string{"profile", "delete", "--name=test"}); err != nil {
		t.Fatalf("profile delete failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "test.yaml")); !os.IsNotExist(err) {
		t.Fatalf("expected test.yaml to be deleted, stat error: %v", err)
	}
}
