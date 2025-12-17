package display

import "image/color"

type Renderer interface {
	SetPixel(x, y int16, c color.RGBA)
	Size() (width, height int16)
	Clear()
	Display() error
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
}
