package extract

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hueryang/practice/internal/llm"
)

// DefaultModel is a fast, cost-effective chat model on BigModel (see docs enum).
const DefaultModel = "glm-4-flash-250414"

// MeetingMinutes is structured meeting info extracted from raw text.
type MeetingMinutes struct {
	Title          string       `json:"title"`
	KeyConclusions []string     `json:"key_conclusions"`
	ActionItems    []ActionItem `json:"action_items"`
	Participants   []string     `json:"participants"`
}

// ActionItem is one todo from the meeting.
type ActionItem struct {
	Content string `json:"content"`
	Owner   string `json:"owner"`
	Due     string `json:"due"`
}

// FromMinutesText calls the LLM to extract structured meeting minutes as JSON.
func FromMinutesText(client *llm.Client, model, rawText string) (*MeetingMinutes, error) {
	if model == "" {
		model = DefaultModel
	}
	sys := `你是会议纪要分析助手。用户会提供一段会议纪要原文。
请只输出一个 JSON 对象，不要输出 Markdown 代码块或其它说明文字。
JSON 字段与含义：
- "title": 字符串，会议主题或标题；若原文未明确则概括一句。
- "key_conclusions": 字符串数组，关键结论与要点，条数适中。
- "action_items": 对象数组，每项含 "content"(待办描述)、"owner"(负责人，未知则 "")、"due"(截止时间，未知则 "")。
- "participants": 字符串数组，参会人或相关方；若原文未提及则 []。
输出语言与原文一致（中文原文则用中文）。`

	user := "以下为会议纪要原文：\n\n" + strings.TrimSpace(rawText)

	req := &llm.ChatCompletionRequest{
		Model: model,
		Messages: []llm.ChatMessage{
			{Role: "system", Content: sys},
			{Role: "user", Content: user},
		},
		Temperature: 0.3,
		Stream:      false,
		ResponseFormat: &llm.ResponseFormat{
			Type: "json_object",
		},
	}

	resp, err := client.ChatCompletion(req)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	payload, err := normalizeAssistantJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("解析模型返回的 JSON 文本: %w", err)
	}
	var out MeetingMinutes
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil, fmt.Errorf("反序列化会议纪要 JSON: %w", err)
	}
	return &out, nil
}

var fenceJSON = regexp.MustCompile("(?s)```(?:json)?\\s*([\\s\\S]*?)```")

func normalizeAssistantJSON(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("模型返回内容为空")
	}
	if m := fenceJSON.FindStringSubmatch(s); len(m) == 2 {
		s = strings.TrimSpace(m[1])
	}
	return []byte(s), nil
}
