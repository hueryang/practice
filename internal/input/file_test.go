package input_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hueryang/practice/internal/input"
)

func TestReadMeetingFile(t *testing.T) {
	t.Parallel()

	t.Run("success utf8 and newlines", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "minutes.txt")
		want := "主题：周会\n结论：按期发布\n"
		if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := input.ReadMeetingFile(path)
		if err != nil {
			t.Fatalf("ReadMeetingFile: %v", err)
		}
		if string(got) != want {
			t.Fatalf("content: got %q want %q", got, want)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "empty.txt")
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := input.ReadMeetingFile(path)
		if err != nil {
			t.Fatalf("ReadMeetingFile: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("len: got %d want 0", len(got))
		}
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()
		_, err := input.ReadMeetingFile(filepath.Join(t.TempDir(), "nope.txt"))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
