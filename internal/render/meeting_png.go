package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/hueryang/practice/internal/extract"
)

const (
	canvasWidth  = 1080
	margin       = 48.0
	titleSize    = 30.0
	sectionSize  = 20.0
	bodySize     = 16.0
	lineFactor   = 1.45
	sectionGap   = 28.0
	paragraphGap = 14.0
	maxCanvasH   = 20000
)

// WriteMeetingReportPNG renders structured minutes to a PNG file.
func WriteMeetingReportPNG(outPath string, m *extract.MeetingMinutes) error {
	fontPath, err := FindFont()
	if err != nil {
		return err
	}

	dc := gg.NewContext(canvasWidth, maxCanvasH)
	dc.SetColor(color.RGBA{R: 250, G: 252, B: 255, A: 255})
	dc.Clear()

	x0 := margin
	maxW := float64(canvasWidth) - 2*margin
	y := margin

	measureWrap := func(size float64, txt string) []string {
		_ = dc.LoadFontFace(fontPath, size)
		txt = strings.TrimSpace(txt)
		return wrapLines(txt, maxW, func(s string) float64 {
			w, _ := dc.MeasureString(s)
			return w
		})
	}

	drawCentered := func(size float64, txt string) {
		_ = dc.LoadFontFace(fontPath, size)
		dc.SetRGB(0.12, 0.14, 0.18)
		lines := measureWrap(size, txt)
		lineH := size * lineFactor
		for _, line := range lines {
			w, _ := dc.MeasureString(line)
			dc.DrawString(line, x0+(maxW-w)/2, y+size)
			y += lineH
		}
	}

	drawLeft := func(size float64, txt string) {
		_ = dc.LoadFontFace(fontPath, size)
		dc.SetRGB(0.12, 0.14, 0.18)
		lines := measureWrap(size, txt)
		lineH := size * lineFactor
		for _, line := range lines {
			dc.DrawString(line, x0, y+size)
			y += lineH
		}
	}

	title := strings.TrimSpace(m.Title)
	if title == "" {
		title = "会议纪要"
	}
	drawCentered(titleSize, title)
	y += sectionGap

	_ = dc.LoadFontFace(fontPath, sectionSize)
	dc.SetRGB(0.2, 0.45, 0.85)
	dc.DrawString("关键结论", x0, y+sectionSize)
	y += sectionSize*lineFactor + paragraphGap

	_ = dc.LoadFontFace(fontPath, bodySize)
	dc.SetRGB(0.12, 0.14, 0.18)
	if len(m.KeyConclusions) == 0 {
		drawLeft(bodySize, "（无）")
	} else {
		for i, c := range m.KeyConclusions {
			c = strings.TrimSpace(c)
			if c == "" {
				continue
			}
			block := "• " + c
			lines := measureWrap(bodySize, block)
			lineH := bodySize * lineFactor
			for _, line := range lines {
				_ = dc.LoadFontFace(fontPath, bodySize)
				dc.SetRGB(0.12, 0.14, 0.18)
				dc.DrawString(line, x0, y+bodySize)
				y += lineH
			}
			if i < len(m.KeyConclusions)-1 {
				y += paragraphGap / 2
			}
		}
	}

	y += sectionGap
	_ = dc.LoadFontFace(fontPath, sectionSize)
	dc.SetRGB(0.2, 0.45, 0.85)
	dc.DrawString("待办事项", x0, y+sectionSize)
	y += sectionSize*lineFactor + paragraphGap

	_ = dc.LoadFontFace(fontPath, bodySize)
	dc.SetRGB(0.12, 0.14, 0.18)
	if len(m.ActionItems) == 0 {
		drawLeft(bodySize, "（无）")
	} else {
		n := 0
		for _, it := range m.ActionItems {
			line := strings.TrimSpace(it.Content)
			if line == "" {
				continue
			}
			n++
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
			block := fmt.Sprintf("%d. %s", n, line)
			lines := measureWrap(bodySize, block)
			lineH := bodySize * lineFactor
			for _, ln := range lines {
				_ = dc.LoadFontFace(fontPath, bodySize)
				dc.SetRGB(0.12, 0.14, 0.18)
				dc.DrawString(ln, x0, y+bodySize)
				y += lineH
			}
			y += paragraphGap / 2
		}
		if n == 0 {
			drawLeft(bodySize, "（无）")
		}
	}

	if len(m.Participants) > 0 {
		parts := strings.Join(trimStrings(m.Participants), "、")
		if parts != "" {
			y += sectionGap / 2
			_ = dc.LoadFontFace(fontPath, bodySize)
			dc.SetRGB(0.35, 0.37, 0.42)
			block := "参与人：" + parts
			lines := measureWrap(bodySize, block)
			lineH := bodySize * lineFactor
			for _, line := range lines {
				dc.DrawString(line, x0, y+bodySize)
				y += lineH
			}
		}
	}

	y += margin
	h := int(y + margin)
	if h < 400 {
		h = 400
	}
	if h > maxCanvasH {
		h = maxCanvasH
	}

	r := image.Rect(0, 0, canvasWidth, h)
	cropped := image.NewRGBA(r)
	draw.Draw(cropped, r, dc.Image(), image.Point{}, draw.Src)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("创建输出目录: %w", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("创建输出文件 %q: %w", outPath, err)
	}
	defer f.Close()
	if err := png.Encode(f, cropped); err != nil {
		return fmt.Errorf("写入 PNG %q: %w", outPath, err)
	}
	return nil
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
