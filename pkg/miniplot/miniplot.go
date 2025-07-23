package miniplot

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
)

type MiniPlot struct {
	display       drivers.Displayer
	font          tinyfont.Fonter
	fontWidth     int16 // Width of a single character in the font
	fontHeight    int16 // Height of a single character in the font
	DisplayWidth  int16
	DisplayHeight int16
	Color         color.RGBA
}

type Drawer interface {
	Line(display drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, color color.RGBA)
	Rectangle(display drivers.Displayer, x int16, y int16, w int16, h int16, color color.RGBA) error
	FilledRectangle(display drivers.Displayer, x int16, y int16, w int16, h int16, color color.RGBA) error
	Circle(display drivers.Displayer, x0 int16, y0 int16, r int16, color color.RGBA)
	FilledCircle(display drivers.Displayer, x0 int16, y0 int16, r int16, color color.RGBA)
	Triangle(display drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA)
	FilledTriangle(display drivers.Displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA)
}

type Plotter interface {
	// DrawLineChart draws a line chart on the display from the provided data.
	// Plot automatically adjusts height to fit the display.
	// If length of data is less than the display width, it will fill the rest
	// with zeros. Drawing will start from right to left.
	// If the data is longer than the display width, it will be truncated.
	DrawLineChart([]int) error
}

func NewMiniPlot(
	display drivers.Displayer,
	font tinyfont.Fonter,
	displayWidth int16,
	displayHeight int16,
) *MiniPlot {
	return &MiniPlot{
		display:       display,
		font:          font,
		fontWidth:     int16(tinyfont.GetGlyph(font, '0').Info().Width),
		fontHeight:    int16(tinyfont.GetGlyph(font, '0').Info().Height),
		DisplayWidth:  displayWidth,
		DisplayHeight: displayHeight,
		Color:         color.RGBA{R: 1, G: 1, B: 1, A: 255},
	}
}

func (mp *MiniPlot) drawText(x, y int16, text string) {
	tinyfont.WriteLine(mp.display, mp.font, x, y, text, mp.Color)
}

// drawAxis draws the axis lines and labels for the plot.
// It should draw values -100, -50 and 0 on the Y-axis,
// and two labels on the X-axis: 0 and the maximum value of the data.
func (mp *MiniPlot) drawAxis(maxValue int) {
}

// drawGrid draws a grid on the plot.
// It should draw horizontal lines with space for each 25% of the Y-axis,
func (mp *MiniPlot) drawGrid() {
}

// drawData draws the data points on the plot.
func (mp *MiniPlot) drawData(data []int) {
	// Calculate the scaling factor to fit the data into the display height.
	maxValue := 0
	for _, value := range data {
		if value > maxValue {
			maxValue = value
		}
	}

	if maxValue == 0 {
		return // No data to draw
	}

	scaleFactor := float64(mp.DisplayHeight) / float64(maxValue)

	for i, value := range data {
		x := int16(i)
		y := int16(mp.DisplayHeight - int16(float64(value)*scaleFactor))
		tinydraw.FilledCircle(mp.display, x, y, 2, mp.Color)
	}
}
