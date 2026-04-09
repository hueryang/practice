package cli_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hueryang/practice/internal/cli"
	"github.com/hueryang/practice/internal/llm"
)

func TestExtractCommand_PrintsStructuredOutput(t *testing.T) {
	minutesJSON := `{"title":"CLI测试会","key_conclusions":["结论A"],"action_items":[{"content":"跟进事项","owner":"","due":""}],"participants":[]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/paas/v4/chat/completions" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": minutesJSON}},
			},
		})
	}))
	t.Cleanup(srv.Close)

	t.Setenv(llm.EnvAPIKey, "k")
	t.Setenv(llm.EnvBaseURL, srv.URL+"/api/paas/v4")

	dir := t.TempDir()
	path := filepath.Join(dir, "m.txt")
	if err := os.WriteFile(path, []byte("任意纪要正文"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout strings.Builder
	var stderr strings.Builder
	root := cli.NewRootCmd()
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs([]string{"extract", path})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v stderr=%q", err, stderr.String())
	}
	out := stdout.String()
	if !containsAll(out, "CLI测试会", "关键结论", "待办事项", "结论A", "跟进事项") {
		t.Fatalf("unexpected stdout:\n%s", out)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}
