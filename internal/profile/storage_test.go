package profile

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadAndCreatesYAMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(oldWd); chdirErr != nil {
			t.Fatalf("restore working directory: %v", chdirErr)
		}
	})
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("change working directory to temp dir %q: %v", tmpDir, err)
	}

	expected := newProfile("p1", "u1", "pr1")
	if err := Save(expected); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := Load("p1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got.Name != "p1" {
		t.Fatalf("Name = %q, want %q", got.Name, "p1")
	}
	if got.Data.User != "u1" {
		t.Fatalf("Data.User = %q, want %q", got.Data.User, "u1")
	}
	if got.Data.Project != "pr1" {
		t.Fatalf("Data.Project = %q, want %q", got.Data.Project, "pr1")
	}

	if _, err := os.Stat(filepath.Join("p1.yaml")); err != nil {
		t.Fatalf("expected file %q to be created: %v", "p1.yaml", err)
	}
}

func TestExistsLifecycleAndRemoveMissing(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(oldWd); chdirErr != nil {
			t.Fatalf("restore working directory: %v", chdirErr)
		}
	})
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("change working directory to temp dir %q: %v", tmpDir, err)
	}

	ok, err := exists("p1")
	if err != nil {
		t.Fatalf("exists before save: %v", err)
	}
	if ok {
		t.Fatalf("exists(\"p1\") before save = true, want false")
	}

	if err = Save(newProfile("p1", "user1", "project1")); err != nil {
		t.Fatalf("save profile: %v", err)
	}

	ok, err = exists("p1")
	if err != nil {
		t.Fatalf("exists after save: %v", err)
	}
	if !ok {
		t.Fatalf("exists(\"p1\") after save = false, want true")
	}

	if err = Remove("p1"); err != nil {
		t.Fatalf("remove existing profile: %v", err)
	}

	ok, err = exists("p1")
	if err != nil {
		t.Fatalf("exists after remove: %v", err)
	}
	if ok {
		t.Fatalf("exists(\"p1\") after remove = true, want false")
	}

	err = Remove("missing")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("remove missing profile error = %v, want os.ErrNotExist", err)
	}

	if _, statErr := os.Stat(filepath.Join(tmpDir, "p1.yaml")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("expected removed file p1.yaml to not exist, got: %v", statErr)
	}
}

func TestGetProfileName(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		want      string
		wantErrIs error
	}{
		{
			name:     "valid simple yaml",
			fileName: "abc.yaml",
			want:     "abc",
		},
		{
			name:     "valid dotted yaml",
			fileName: "a.b.yaml",
			want:     "a.b",
		},
		{
			name:      "invalid yml extension",
			fileName:  "abc.yml",
			wantErrIs: ErrNotYAML,
		},
		{
			name:      "invalid no extension",
			fileName:  "abc",
			wantErrIs: ErrNotYAML,
		},
		{
			name:      "invalid temporary extension",
			fileName:  ".yaml.tmp",
			wantErrIs: ErrNotYAML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getProfileName(tt.fileName)
			if tt.wantErrIs != nil {
				if !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("getProfileName(%q) error = %v, want errors.Is(err, %v)", tt.fileName, err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("getProfileName(%q) unexpected error: %v", tt.fileName, err)
			}
			if got != tt.want {
				t.Fatalf("getProfileName(%q) = %q, want %q", tt.fileName, got, tt.want)
			}
		})
	}
}

func TestValidateFileStructure(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []struct {
		name          string
		content       string
		expected      FileStructure
		expectProfile ProfileData
		checkData     bool
	}{
		{
			name:     "valid",
			content:  "user: alice\nproject: demo\n",
			expected: Valid,
			expectProfile: ProfileData{
				User:    "alice",
				Project: "demo",
			},
			checkData: true,
		},
		{
			name:     "extended",
			content:  "user: bob\nproject: prod\nextra: value\n",
			expected: ValidOrExtended,
			expectProfile: ProfileData{
				User:    "bob",
				Project: "prod",
			},
			checkData: true,
		},
		{
			name:     "invalid-structure",
			content:  "user: [oops\nproject: demo\n",
			expected: All,
		},
		{
			name:     "invalid-values-empty-user",
			content:  "user: \nproject: demo\n",
			expected: All,
		},
		{
			name:     "invalid-values-empty-project",
			content:  "user: alice\nproject: \n",
			expected: All,
		},
		{
			name:     "invalid-values-empty-both",
			content:  "user: \nproject: \n",
			expected: All,
		},
		{
			name:     "invalid-values-user-multiline",
			content:  "user: |\n  alice\n  bob\nproject: demo\n",
			expected: All,
		},
		{
			name:     "invalid-values-project-multiline",
			content:  "user: alice\nproject: |\n  demo\n  other\n",
			expected: All,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tc.name+".yaml")
			if err := os.WriteFile(path, []byte(tc.content), 0o644); err != nil {
				t.Fatalf("write test file: %v", err)
			}

			p := newDefaultProfile(tc.name)
			struc, err := validateFileStructure(path, &p)
			if err != nil {
				t.Fatalf("validateFileStructure() error = %v", err)
			}
			if struc != tc.expected {
				t.Fatalf("validateFileStructure() = %v, want %v", struc, tc.expected)
			}

			if tc.checkData && p.Data != tc.expectProfile {
				t.Fatalf("profile data = %+v, want %+v", p.Data, tc.expectProfile)
			}
		})
	}
}

func TestSearchAll_FilterModesAndIgnoreNonYAML(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"ok.yaml":       "user: alice\nproject: app\n",
		"extended.yaml": "user: bob\nproject: api\nextra: value\n",
		"bad.yaml":      "user: \nproject: broken\n",
		"readme.txt":    "this file must be ignored by SearchAll\n",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	t.Run("Valid returns only ok", func(t *testing.T) {
		got, err := SearchAll(tmpDir, Valid)
		if err != nil {
			t.Fatalf("SearchAll(Valid): %v", err)
		}
		assertProfileNames(t, got, []string{"ok"}, true)
	})

	t.Run("ValidOrExtended returns ok and extended", func(t *testing.T) {
		got, err := SearchAll(tmpDir, ValidOrExtended)
		if err != nil {
			t.Fatalf("SearchAll(ValidOrExtended): %v", err)
		}
		assertProfileNames(t, got, []string{"extended", "ok"}, true)
	})

	t.Run("All returns all YAML including invalid", func(t *testing.T) {
		got, err := SearchAll(tmpDir, All)
		if err != nil {
			t.Fatalf("SearchAll(All): %v", err)
		}
		assertProfileNames(t, got, []string{"bad", "extended", "ok"}, true)
	})
}
