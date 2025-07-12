package display

import "pico_co2/internal/display/font"

func (f *FontDisplay) DisplayError(msg string) {
	// Split longer messages into multiple lines
	lines := []string{msg}

	fp := font.NewFont7(f.display)

	maxCharPerLine := f.width / fp.Width()
	if len(msg) > int(maxCharPerLine) {
		lines = splitStringToLines(msg, int(maxCharPerLine))
	}

	f.printLines(lines, fp)
}

func (f *FontDisplay) printLines(lines []string, fp font.FontPrinter) {
	if f == nil || len(lines) == 0 {
		return
	}

	maxLines := f.height / (fp.Height() + 1)

	for i, line := range lines {
		if i >= int(maxLines) {
			break
		}

		fp.Print(0, int16(i)*(fp.Height()+1), line)
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
