package cli

import (
	"bytes"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteWithArgs_LogsErrorOnMissingFile(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prev) })
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError})))

	code := ExecuteWithArgs([]string{"read", filepath.Join(t.TempDir(), "missing.txt")})
	if code != 1 {
		t.Fatalf("exit code: got %d want 1", code)
	}
	out := buf.String()
	if !strings.Contains(out, "命令执行失败") {
		t.Fatalf("expected log to contain 命令执行失败, got: %q", out)
	}
	if !strings.Contains(out, "error") {
		t.Fatalf("expected slog text output, got: %q", out)
	}
}
