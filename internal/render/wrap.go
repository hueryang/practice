package render

import "strings"

// wrapLines splits s into lines that fit within maxWidth when measured by measureWidth.
func wrapLines(s string, maxWidth float64, measureWidth func(string) float64) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var lines []string
	for _, para := range strings.Split(s, "\n") {
		if para == "" {
			lines = append(lines, "")
			continue
		}
		var b strings.Builder
		for _, r := range para {
			trial := b.String() + string(r)
			if measureWidth(trial) > maxWidth && b.Len() > 0 {
				lines = append(lines, b.String())
				b.Reset()
				b.WriteRune(r)
				continue
			}
			b.WriteRune(r)
		}
		if b.Len() > 0 {
			lines = append(lines, b.String())
		}
	}
	return lines
}
