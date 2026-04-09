package render

import (
	"fmt"
	"os"
)

// EnvFontPath is the env var for an explicit TTF/OTF font path (Chinese-capable recommended).
const EnvFontPath = "MEETING_REPORT_FONT"

// FindFont returns a usable font file path for drawing text, preferring EnvFontPath when set.
func FindFont() (string, error) {
	if p := os.Getenv(EnvFontPath); p != "" {
		st, err := os.Stat(p)
		if err != nil {
			return "", fmt.Errorf("字体路径 %q 无效: %w", p, err)
		}
		if st.IsDir() {
			return "", fmt.Errorf("字体路径 %q 是目录", p)
		}
		return p, nil
	}

	candidates := []string{
		"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
		"/Library/Fonts/Arial Unicode.ttf",
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.otf",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttf",
		`C:\Windows\Fonts\msyh.ttc`,
		`C:\Windows\Fonts\simhei.ttf`,
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("未找到可用的中文字体文件，请设置环境变量 %s 指向 .ttf/.otf 文件", EnvFontPath)
}
