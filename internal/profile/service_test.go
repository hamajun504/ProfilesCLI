package profile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func chdirToTempDir(t *testing.T) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir to temp dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
}

func TestServiceLifecycle_ValidAttributes(t *testing.T) {
	cases := []struct {
		name    string
		profile string
		user1   string
		proj1   string
		user2   string
		proj2   string
	}{
		{name: "simple", profile: "p1", user1: "u1", proj1: "pr1", user2: "u2", proj2: "pr2"},
		{name: "with_hyphen", profile: "team-1", user1: "alice", proj1: "backend", user2: "bob", proj2: "frontend"},
		{name: "with_underscore", profile: "dev_ops", user1: "svc-account", proj1: "infra_2026", user2: "svc-account-2", proj2: "infra_next"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chdirToTempDir(t)

			if err := Create(tc.profile, tc.user1, tc.proj1); err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			exists, err := Exists(tc.profile)
			if err != nil {
				t.Fatalf("Exists() error = %v", err)
			}
			if !exists {
				t.Fatalf("Exists() = false, want true")
			}

			p, err := Get(tc.profile)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if p.Data.User != tc.user1 || p.Data.Project != tc.proj1 {
				t.Fatalf("Get() = user/project %q/%q, want %q/%q", p.Data.User, p.Data.Project, tc.user1, tc.proj1)
			}

			if err := Update(tc.profile, tc.user2, tc.proj2); err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			p, err = Get(tc.profile)
			if err != nil {
				t.Fatalf("Get() after update error = %v", err)
			}
			if p.Data.User != tc.user2 || p.Data.Project != tc.proj2 {
				t.Fatalf("Get() after update = user/project %q/%q, want %q/%q", p.Data.User, p.Data.Project, tc.user2, tc.proj2)
			}

			if err := Delete(tc.profile); err != nil {
				t.Fatalf("Delete() error = %v", err)
			}

			exists, err = Exists(tc.profile)
			if err != nil {
				t.Fatalf("Exists() after delete error = %v", err)
			}
			if exists {
				t.Fatalf("Exists() after delete = true, want false")
			}

			if err := Delete(tc.profile); err != nil {
				t.Fatalf("second Delete() error = %v", err)
			}
		})
	}
}

func TestServiceLifecycle_InvalidAttributes(t *testing.T) {
	cases := []struct {
		name      string
		profile   string
		user      string
		project   string
		existsErr bool
		getErr    bool
		deleteErr bool
		updateErr bool
		createErr bool
	}{
		{name: "invalid_name", profile: "bad/name", user: "u1", project: "pr1", existsErr: true, getErr: true, deleteErr: true, updateErr: true, createErr: true},
		{name: "empty_user", profile: "p1", user: "", project: "pr1", existsErr: false, getErr: true, deleteErr: false, updateErr: true, createErr: true},
		{name: "project_with_newline", profile: "p2", user: "u2", project: "pr2\nnext", existsErr: false, getErr: true, deleteErr: false, updateErr: true, createErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chdirToTempDir(t)

			err := Create(tc.profile, tc.user, tc.project)
			if (err != nil) != tc.createErr {
				t.Fatalf("Create() error = %v, createErr=%v", err, tc.createErr)
			}

			exists, err := Exists(tc.profile)
			if (err != nil) != tc.existsErr {
				t.Fatalf("Exists() error = %v, existsErr=%v", err, tc.existsErr)
			}
			if err == nil && exists {
				t.Fatalf("Exists() = true, want false")
			}

			_, err = Get(tc.profile)
			if (err != nil) != tc.getErr {
				t.Fatalf("Get() error = %v, getErr=%v", err, tc.getErr)
			}

			err = Update(tc.profile, tc.user, tc.project)
			if (err != nil) != tc.updateErr {
				t.Fatalf("Update() error = %v, updateErr=%v", err, tc.updateErr)
			}

			err = Delete(tc.profile)
			if (err != nil) != tc.deleteErr {
				t.Fatalf("Delete() error = %v, deleteErr=%v", err, tc.deleteErr)
			}
			if !tc.deleteErr {
				if err := Delete(tc.profile); err != nil && !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("second Delete() error = %v", err)
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Run("returns already exists error on duplicate profile", func(t *testing.T) {
		tempDir := t.TempDir()

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("get working directory: %v", err)
		}
		t.Cleanup(func() {
			if chdirErr := os.Chdir(cwd); chdirErr != nil {
				t.Fatalf("restore working directory: %v", chdirErr)
			}
		})

		if err := os.Chdir(filepath.Clean(tempDir)); err != nil {
			t.Fatalf("change to temp dir: %v", err)
		}

		if err := Create("p1", "user1", "project1"); err != nil {
			t.Fatalf("create initial profile: %v", err)
		}

		err = Create("p1", "user2", "project2")
		if err == nil {
			t.Fatal("expected duplicate profile error, got nil")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("expected error to contain %q, got %q", "already exists", err.Error())
		}
	})
}

func TestSearchAllFiltersByMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	files := map[string]string{
		"ok.yaml":            "user: alice\nproject: demo\n",
		"extended.yaml":      "user: alice\nproject: demo\nextra: value\n",
		"invalid.yaml":       "user: \nproject: demo\n",
		"invalid_field.yaml": "user: alice\nproject: demo\ninvalid_field: \"line1\\nline2\"\n",
		"note.txt":           "not a yaml profile\n",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	t.Run("valid only", func(t *testing.T) {
		profiles, err := SearchAll(dir, Valid)
		if err != nil {
			t.Fatalf("SearchAll Valid returned error: %v", err)
		}
		assertProfileNames(t, profiles, []string{"ok"}, false)
	})

	t.Run("valid or extended", func(t *testing.T) {
		profiles, err := SearchAll(dir, ValidOrExtended)
		if err != nil {
			t.Fatalf("SearchAll ValidOrExtended returned error: %v", err)
		}
		assertProfileNames(t, profiles, []string{"extended", "invalid_field", "ok"}, false)
	})

	t.Run("all yaml profiles", func(t *testing.T) {
		profiles, err := SearchAll(dir, All)
		if err != nil {
			t.Fatalf("SearchAll All returned error: %v", err)
		}
		assertProfileNames(t, profiles, []string{"extended", "invalid", "invalid_field", "ok"}, false)
	})
}

func TestListReturnsSortedByName(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	files := map[string]string{
		"zeta.yaml":  "user: user-z\nproject: project-z\n",
		"alpha.yaml": "user: user-a\nproject: project-a\n",
		"beta.yaml":  "user: user-b\nproject: project-b\n",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(currentDir); chdirErr != nil {
			t.Fatalf("restore cwd: %v", chdirErr)
		}
	})

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir temp dir: %v", err)
	}

	profiles, err := List(Valid)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	assertProfileNames(t, profiles, []string{"alpha", "beta", "zeta"}, false)
}
