package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type runResult struct {
	stdout string
	stderr string
	exit   int
	err    error
}

func buildCLI(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve test file path")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))

	bin := filepath.Join(t.TempDir(), "profilescli")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/mws")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, string(out))
	}
	return bin
}

func runCmd(t *testing.T, bin, dir string, args ...string) runResult {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	var outb, errb strings.Builder
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			return runResult{stdout: outb.String(), stderr: errb.String(), exit: code, err: err}
		}
	}
	return runResult{stdout: outb.String(), stderr: errb.String(), exit: code, err: nil}
}

func TestProfileLifecycle(t *testing.T) {
	bin := buildCLI(t)
	dir := t.TempDir()

	r := runCmd(t, bin, dir, "profile", "create", "-name", "dev", "-user", "alice", "-project", "app")
	if r.exit != 0 || r.stdout != "" || r.stderr != "" {
		t.Fatalf("create unexpected result: %+v", r)
	}

	r = runCmd(t, bin, dir, "profile", "get", "-name", "dev")
	if r.exit != 0 || !strings.Contains(r.stdout, "profile:  dev") || r.stderr != "" {
		t.Fatalf("get unexpected result: %+v", r)
	}

	r = runCmd(t, bin, dir, "profile", "list")
	if r.exit != 0 || !strings.Contains(r.stdout, "dev") || r.stderr != "" {
		t.Fatalf("list unexpected result: %+v", r)
	}

	r = runCmd(t, bin, dir, "profile", "delete", "-name", "dev")
	if r.exit != 0 || r.stdout != "" || r.stderr != "" {
		t.Fatalf("delete unexpected result: %+v", r)
	}
}

func TestInvalidFlagsAndMissingRequired(t *testing.T) {
	bin := buildCLI(t)
	dir := t.TempDir()

	cases := []struct {
		name string
		args []string
		want string
	}{
		{"create extra flag", []string{"profile", "create", "-name", "n", "-user", "u", "-project", "p", "-x"}, "flag provided but not defined"},
		{"create missing name", []string{"profile", "create", "-user", "u", "-project", "p"}, "name is required"},
		{"get extra flag", []string{"profile", "get", "-name", "n", "-x"}, "flag provided but not defined"},
		{"get missing name", []string{"profile", "get"}, "name is required"},
		{"list extra flag", []string{"profile", "list", "-x"}, "flag provided but not defined"},
		{"delete extra flag", []string{"profile", "delete", "-name", "n", "-x"}, "flag provided but not defined"},
		{"delete missing name", []string{"profile", "delete"}, "name is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := runCmd(t, bin, dir, tc.args...)
			if r.exit == 0 || !strings.Contains(r.stderr, tc.want) {
				t.Fatalf("unexpected result: %+v; want stderr containing %q", r, tc.want)
			}
		})
	}
}

func TestInvalidNameUserProject(t *testing.T) {
	bin := buildCLI(t)
	dir := t.TempDir()

	cases := []struct {
		name string
		args []string
		want string
	}{
		{"create invalid name", []string{"profile", "create", "-name", "bad/name", "-user", "u", "-project", "p"}, "name must not contain path separators"},
		{"create invalid user", []string{"profile", "create", "-name", "ok", "-user", " ", "-project", "p"}, "user is required"},
		{"create invalid project", []string{"profile", "create", "-name", "ok", "-user", "u", "-project", "\n"}, "project is required"},
		{"get invalid name", []string{"profile", "get", "-name", "bad/name"}, "name must not contain path separators"},
		{"delete invalid name", []string{"profile", "delete", "-name", "bad/name"}, "name must not contain path separators"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := runCmd(t, bin, dir, tc.args...)
			if r.exit == 0 || !strings.Contains(r.stderr, tc.want) {
				t.Fatalf("unexpected result: %+v; want stderr containing %q", r, tc.want)
			}
		})
	}
}

func TestFileErrors(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("permission checks are unreliable as root")
	}

	bin := buildCLI(t)
	dir := t.TempDir()

	r := runCmd(t, bin, dir, "profile", "get", "-name", "missing")
	if r.exit == 0 || !strings.Contains(r.stderr, "no such file or directory") {
		t.Fatalf("expected get missing file error, got: %+v", r)
	}

	r = runCmd(t, bin, dir, "profile", "list")
	if r.exit != 0 {
		t.Fatalf("list should pass on empty dir: %+v", r)
	}

	noRead := filepath.Join(dir, "no-read")
	if err := os.Mkdir(noRead, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(noRead, 0o755)
	r = runCmd(t, bin, noRead, "profile", "list")
	if r.err == nil && r.exit == 0 {
		t.Fatalf("expected list read permission error, got: %+v", r)
	}
	if r.err != nil && !strings.Contains(strings.ToLower(r.err.Error()), "permission denied") {
		t.Fatalf("unexpected start error: %v", r.err)
	}

	noWrite := filepath.Join(dir, "no-write")
	if err := os.Mkdir(noWrite, 0o555); err != nil {
		t.Fatal(err)
	}
	r = runCmd(t, bin, noWrite, "profile", "create", "-name", "dev", "-user", "u", "-project", "p")
	if r.exit == 0 {
		t.Fatalf("expected create write permission error, got: %+v", r)
	}
}
