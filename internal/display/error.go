package display

import (
	font "github.com/Nondzu/ssd1306_font"
)

func (f *FontDisplay) DisplayError(msg string) {
	// Split longer messages into multiple lines
	lines := []string{msg}
	maxCharPerLine := 128 / 7
	if len(msg) > maxCharPerLine {
		lines = splitStringToLines(msg, maxCharPerLine)
	}

	f.printLines(lines)
}

func (f *FontDisplay) printLines(lines []string) {
	if f == nil || len(lines) == 0 {
		return
	}

	f.clearDisplay()

	f.font.Configure(font.Config{FontType: font.FONT_7x10})

	for i, line := range lines {
		if i >= 3 { // Max 3 lines on a 32px display
			break
		}

		f.font.XPos = 0
		f.font.YPos = int16(i * 11)
		f.font.PrintText(line)
	}
}

// Split string into multiple lines if it exceeds the display width
func splitStringToLines(s string, maxCharPerLine int) []string {
	lines := make([]string, 0)
	for i := 0; i < len(s); i += maxCharPerLine {
		end := i + maxCharPerLine
		if end > len(s) {
			end = len(s)
		}
		lines = append(lines, s[i:end])
	}

	return lines
}
