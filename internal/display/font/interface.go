package font

import "tinygo.org/x/tinyfont"

type FontPrinter interface {
	// Print draws the text at the specified x, y coordinates, where x and y are
	// the top-left corner of the text.
	Print(x, y int16, text string) int16
	CalcWidth(text string) int16
	Height() int16
	Width() int16
	GetFont() tinyfont.Fonter
}
