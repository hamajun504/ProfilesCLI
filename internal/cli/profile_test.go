package cli

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type errReader struct{}

func (e errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestAskToOverwriteFromReaderAnswers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "y with newline", input: "y\n", want: true},
		{name: "yes with newline", input: "yes\n", want: true},
		{name: "n with newline", input: "n\n", want: false},
		{name: "empty with newline", input: "\n", want: false},
		{name: "random text", input: "something\n", want: false},
		{name: "EOF without newline y", input: "y", want: true},
		{name: "EOF without newline yes", input: "yes", want: true},
		{name: "EOF without newline random", input: "random", want: false},
		{name: "EOF without newline empty", input: "", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := askToOverwriteFromReader("demo", strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected result: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAskToOverwriteFromReaderReturnsErrorOnlyForIOFailure(t *testing.T) {
	t.Parallel()

	_, err := askToOverwriteFromReader("demo", errReader{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

var _ io.Reader = errReader{}
