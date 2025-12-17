package display

import (
	"image/color"
	"pico_co2/internal/display/font"
)

type Renderer interface {
	SetPixel(x, y int16, c color.RGBA)
	Size() (width, height int16)
	Clear()
	Display() error
	DrawPlot(data []int16, title string)
	DrawTwoSideBar(
		x, y int16,
		value int16,
		label string,
		leftCount int16,
		rightCount int16,
	) int16
	DrawSparkline(
		x, y int16,
		data []int16,
		width int16,
		height int16,
	)
	DrawSquareBar(x, y int16, value uint8)

	// NEW: Unified font management methods
	GetFont(fontType font.FontType) font.FontPrinter
	DrawText(fontType font.FontType, x, y int16, text string)
	CalcTextWidth(fontType font.FontType, text string) int16

	// Deprecated: Legacy font methods for backward compatibility
	// These will be implemented using the new font management system
	DrawXLargeText(x, y int16, text string)
	CalcXLargeTextWidth(text string) int16
	DrawLargeText(x, y int16, text string)
	DrawLargeBoldText(x, y int16, text string)
	DrawLargeSansText(x, y int16, text string)
	CalcLargeTextWidth(text string) int16
	CalcLargeBoldTextWidth(text string) int16
	CalcLargeSansTextWidth(text string) int16
	DrawSmallText(x, y int16, text string)
	CalcSmallTextWidth(text string) int16
	DrawLongText(x, y int16, text string)
}
