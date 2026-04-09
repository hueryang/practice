package render

import (
	"strings"
	"testing"
)

func TestWrapLines(t *testing.T) {
	t.Parallel()
	measure := func(s string) float64 { return float64(len([]rune(s))) * 10 }
	lines := wrapLines("abcdefghijklmnop", 50, measure)
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %q", lines)
	}
	if !strings.Contains(strings.Join(lines, ""), "abcdefghijklmnop") {
		t.Fatalf("lost content: %q", lines)
	}
}
