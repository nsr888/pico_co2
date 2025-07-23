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

func (mp *MiniPlot) DrawLineChart(data []int16) error {
	if len(data) == 0 {
		return nil
	}

	// Find min and max values for scaling
	minVal := data[0]
	maxVal := data[0]
	for _, v := range data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Clear display area
	tinydraw.FilledRectangle(mp.display, 0, 0, mp.DisplayWidth, mp.DisplayHeight, color.RGBA{0, 0, 0, 255})

	// Draw axes
	mp.drawAxis(maxVal, minVal)
	
	// Draw data
	mp.drawData(data, minVal, maxVal)
	
	return nil
}

func (mp *MiniPlot) drawText(x, y int16, text string) {
	tinyfont.WriteLine(mp.display, mp.font, x, y, text, mp.Color)
}

func (mp *MiniPlot) drawAxis(maxValue, minValue int16) {
	// Draw Y-axis line
	tinydraw.Line(mp.display, 20, 0, 20, mp.DisplayHeight-1, mp.Color)
	
	// Draw X-axis line
	tinydraw.Line(mp.display, 20, mp.DisplayHeight-1, mp.DisplayWidth-1, mp.DisplayHeight-1, mp.Color)

	// Draw Y-axis labels
	rangeVal := maxValue - minValue
	if rangeVal == 0 {
		rangeVal = 1
	}
	
	label := mp.formatValue(minValue)
	mp.drawText(0, mp.DisplayHeight-1-mp.fontHeight, label)
	
	label = mp.formatValue((maxValue+minValue)/2)
	mp.drawText(0, mp.DisplayHeight/2, label)
	
	label = mp.formatValue(maxValue)
	mp.drawText(0, 0, label)
}

// drawGrid draws a grid on the plot.
// It should draw horizontal lines with space for each 25% of the Y-axis,
func (mp *MiniPlot) drawGrid() {
}

func (mp *MiniPlot) drawData(data []int16, minValue, maxValue int16) {
	if len(data) == 0 {
		return
	}

	rangeVal := float64(maxValue - minValue)
	if rangeVal == 0 {
		rangeVal = 1
	}

	// Calculate scaling factors
	xScale := float64(mp.DisplayWidth-21) / float64(len(data)-1)
	yScale := float64(mp.DisplayHeight-2) / rangeVal

	// Draw line chart
	prevX := int16(21)
	prevY := int16(float64(maxValue-data[0]) * yScale)
	
	for i, value := range data {
		if i == 0 {
			continue
		}
		
		x := int16(21 + float64(i)*xScale)
		y := int16(float64(maxValue-value) * yScale)
		
		tinydraw.Line(mp.display, prevX, prevY, x, y, mp.Color)
		prevX, prevY = x, y
	}
}

func (mp *MiniPlot) formatValue(value int16) string {
	// Format value to fit in small space
	if value >= 1000 {
		return string(rune('0' + (value/1000)%10)) + "k"
	}
	return string(rune('0' + value%1000))
}
