package profile

import (
	"strings"
	"testing"
)

func TestValidateNewName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
	}{
		{name: "valid simple", input: "abc", wantErr: false},
		{name: "valid mixed", input: "A_B-123", wantErr: false},
		{name: "valid max length 64", input: strings.Repeat("a", 64), wantErr: false},
		{name: "invalid empty", input: "", wantErr: true, errContains: "name is required"},
		{name: "invalid length 65", input: strings.Repeat("a", 65), wantErr: true, errContains: "1-64 characters"},
		{name: "invalid with space", input: "name with space", wantErr: true, errContains: "only letters"},
		{name: "invalid exclamation", input: "name!", wantErr: true, errContains: "only letters"},
		{name: "invalid forward slash", input: "a/b", wantErr: true, errContains: "only letters"},
		{name: "invalid backslash", input: `a\\b`, wantErr: true, errContains: "only letters"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateNewName(tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("validateNewName(%q) error = nil, want error", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("validateNewName(%q) error = %v, want nil", tt.input, err)
			}
			if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Fatalf("validateNewName(%q) error = %q, want to contain %q", tt.input, err.Error(), tt.errContains)
			}
		})
	}
}

func TestValidateOldName(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "legacy dot name", in: "legacy.name", wantErr: false},
		{name: "name with space", in: "profile 1", wantErr: false},
		{name: "simple name", in: "abc", wantErr: false},
		{name: "empty", in: "", wantErr: true},
		{name: "forward slash", in: "a/b", wantErr: true},
		{name: "backslash", in: `a\b`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOldName(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateOldName(%q) error = %v, wantErr = %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty", input: "", wantErr: true},
		{name: "spaces only", input: "   ", wantErr: true},
		{name: "contains LF", input: "foo\nbar", wantErr: true},
		{name: "contains CR", input: "foo\rbar", wantErr: true},
		{name: "valid simple", input: "alice", wantErr: false},
		{name: "valid with surrounding spaces", input: " project-x ", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUser(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateUser(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateProject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty", input: "", wantErr: true},
		{name: "spaces only", input: "   ", wantErr: true},
		{name: "contains LF", input: "foo\nbar", wantErr: true},
		{name: "contains CR", input: "foo\rbar", wantErr: true},
		{name: "valid simple", input: "alice", wantErr: false},
		{name: "valid with surrounding spaces", input: " project-x ", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProject(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateProject(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
