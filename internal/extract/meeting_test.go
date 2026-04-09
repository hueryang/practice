package extract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hueryang/practice/internal/llm"
)

func TestNormalizeAssistantJSON_RawJSON(t *testing.T) {
	t.Parallel()
	raw := `{"title":"x","key_conclusions":[],"action_items":[],"participants":[]}`
	b, err := normalizeAssistantJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	var m MeetingMinutes
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.Title != "x" {
		t.Fatalf("title: %q", m.Title)
	}
}

func TestNormalizeAssistantJSON_Fenced(t *testing.T) {
	t.Parallel()
	raw := "```json\n{\"title\":\" fenced \",\"key_conclusions\":[],\"action_items\":[],\"participants\":[]}\n```"
	b, err := normalizeAssistantJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	var m MeetingMinutes
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(m.Title) != "fenced" {
		t.Fatalf("title: %q", m.Title)
	}
}

func TestFromMinutesText_MockServer(t *testing.T) {
	minutesJSON := `{"title":"周会","key_conclusions":["按期发布"],"action_items":[{"content":"补充文档","owner":"李四","due":"周五"}],"participants":["王五"]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/paas/v4/chat/completions" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("method: %s", r.Method)
		}
		body := map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"content": minutesJSON,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(srv.Close)

	t.Setenv(llm.EnvAPIKey, "test-key")
	t.Setenv(llm.EnvBaseURL, srv.URL+"/api/paas/v4")

	client, err := llm.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client.HTTPClient = srv.Client()

	got, err := FromMinutesText(client, "glm-4-flash-250414", "会议记录原文……")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "周会" {
		t.Fatalf("title: %q", got.Title)
	}
	if len(got.KeyConclusions) != 1 || got.KeyConclusions[0] != "按期发布" {
		t.Fatalf("conclusions: %+v", got.KeyConclusions)
	}
	if len(got.ActionItems) != 1 || got.ActionItems[0].Content != "补充文档" {
		t.Fatalf("actions: %+v", got.ActionItems)
	}
}

func TestFormatHumanReadable(t *testing.T) {
	t.Parallel()
	s := FormatHumanReadable(&MeetingMinutes{
		Title:          "项目评审",
		KeyConclusions: []string{"通过方案 A"},
		ActionItems: []ActionItem{
			{Content: "出图", Owner: "赵六", Due: ""},
		},
		Participants: []string{"钱七"},
	})
	if !strings.Contains(s, "项目评审") || !strings.Contains(s, "通过方案") {
		t.Fatalf("output: %s", s)
	}
}
