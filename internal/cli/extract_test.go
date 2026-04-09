package cli_test

import (
	"bytes"
	"encoding/json"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/hueryang/practice/internal/cli"
	"github.com/hueryang/practice/internal/llm"
)

func TestExtractCommand_WritesPNG(t *testing.T) {
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
	outPNG := filepath.Join(dir, "m.png")

	var stderr bytes.Buffer
	root := cli.NewRootCmd()
	root.SetOut(io.Discard)
	root.SetErr(&stderr)
	root.SetArgs([]string{"extract", path, "-o", outPNG})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute: %v stderr=%q", err, stderr.String())
	}

	st, err := os.Stat(outPNG)
	if err != nil {
		t.Fatalf("stat output png: %v", err)
	}
	if st.Size() == 0 {
		t.Fatal("empty png")
	}
	f, err := os.Open(outPNG)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cfg, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatalf("png.DecodeConfig: %v", err)
	}
	if cfg.Width != 1080 {
		t.Fatalf("width: %d want 1080", cfg.Width)
	}
}
