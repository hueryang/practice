package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/hueryang/practice/internal/cli"
)

func TestReadCommand_PrintsFileContents(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "minutes.txt")
	want := "line1\nline2 中文\n"
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := cli.NewRootCmd()
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs([]string{"read", path})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v\nstderr: %s", err, stderr.String())
	}
	if got := stdout.String(); got != want {
		t.Fatalf("stdout: got %q want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestReadCommand_RequiresExactlyOneArg(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := cli.NewRootCmd()
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs([]string{"read"})

	if err := root.Execute(); err == nil {
		t.Fatal("expected error when file path is missing")
	}
}

func TestReadCommand_MissingFile(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := cli.NewRootCmd()
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs([]string{"read", filepath.Join(t.TempDir(), "missing.txt")})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
}
