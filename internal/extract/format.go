package extract

import (
	"fmt"
	"strings"
)

// FormatHumanReadable prints extracted minutes for CLI stdout (key info first).
func FormatHumanReadable(m *MeetingMinutes) string {
	var b strings.Builder
	fmt.Fprintf(&b, "会议主题：%s\n", strings.TrimSpace(m.Title))
	if len(m.Participants) > 0 {
		fmt.Fprintf(&b, "参与人：%s\n", strings.Join(trimStrings(m.Participants), "、"))
	} else {
		fmt.Fprintf(&b, "参与人：（未提取到）\n")
	}
	fmt.Fprintf(&b, "\n关键结论：\n")
	if len(m.KeyConclusions) == 0 {
		fmt.Fprintf(&b, "（无）\n")
	} else {
		for _, line := range m.KeyConclusions {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(line))
		}
	}
	fmt.Fprintf(&b, "\n待办事项：\n")
	if len(m.ActionItems) == 0 {
		fmt.Fprintf(&b, "（无）\n")
	} else {
		for i, it := range m.ActionItems {
			line := strings.TrimSpace(it.Content)
			extras := []string{}
			if o := strings.TrimSpace(it.Owner); o != "" {
				extras = append(extras, "负责人："+o)
			}
			if d := strings.TrimSpace(it.Due); d != "" {
				extras = append(extras, "截止："+d)
			}
			if len(extras) > 0 {
				line += "（" + strings.Join(extras, "，") + "）"
			}
			fmt.Fprintf(&b, "%d. %s\n", i+1, line)
		}
	}
	return b.String()
}

func trimStrings(ss []string) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
