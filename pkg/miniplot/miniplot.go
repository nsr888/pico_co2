package miniplot

import (
	"errors"
	"fmt"
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
	StartX        int16 // X position to start drawing the plot
	StartY        int16 // Y position to start drawing the plot
	AutoScale     bool  // Whether to automatically scale the Y-axis
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
	c color.RGBA,
) (*MiniPlot, error) {
	if displayWidth <= 0 || displayHeight <= 0 {
		return nil, errors.New("display dimensions must be greater than zero")
	}

	if display == nil {
		return nil, errors.New("display cannot be nil")
	}

	if font == nil {
		return nil, errors.New("font cannot be nil")
	}

	return &MiniPlot{
		display:       display,
		font:          font,
		fontWidth:     int16(tinyfont.GetGlyph(font, '0').Info().Width),
		fontHeight:    int16(tinyfont.GetGlyph(font, '0').Info().Height),
		DisplayWidth:  displayWidth,
		DisplayHeight: displayHeight,
		Color:         c,
		StartX:        20,
		StartY:        22,
	}, nil
}

func (mp *MiniPlot) DrawLineChart(
	data []int16,
	title string,
) error {
	if len(data) == 0 {
		return nil
	}

	// Cut slice to fit display width, only keep the last N samples
	if len(data) > int(mp.DisplayWidth-mp.StartX-2) {
		data = data[len(data)-int(mp.DisplayWidth-mp.StartX-2):]
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
	mp.drawAxis(maxVal, minVal, title)

	// Draw data
	mp.drawData(data, minVal, maxVal)

	return nil
}

func (mp *MiniPlot) drawText(x, y int16, text string) {
	tinyfont.WriteLine(mp.display, mp.font, x, y, text, mp.Color)
}

func (mp *MiniPlot) drawAxis(
	maxValue, minValue int16,
	title string,
) {
	startX := mp.StartX // Start X position for the axis
	startY := mp.StartY // Start Y position for the axis

	// Draw Y-axis line
	tinydraw.Line(mp.display, startX, startY, startX, 0, mp.Color)

	// Draw X-axis line
	tinydraw.Line(mp.display, startX, startY, mp.DisplayWidth-1, startY, mp.Color)

	// Draw Y-axis labels
	rangeVal := maxValue - minValue
	if rangeVal == 0 {
		rangeVal = 1
	}

	// Title in left bottom corner
	mp.drawText(1, startY+mp.fontHeight, title)

	// Y-axis labels
	label := mp.formatValue(minValue)
	mp.drawText(1, startY, label)

	label = mp.formatValue(maxValue)
	mp.drawText(1, mp.fontHeight, label)

	// X-axis labels
	text := "0h"
	textWidth := int16(len(text)) * mp.fontWidth
	halfTextWidth := textWidth / 2
	xPos := mp.DisplayWidth - textWidth
	mp.drawText(xPos, startY+mp.fontHeight, text)

	text = "-1h"
	textWidth = int16(len(text)) * mp.fontWidth
	halfTextWidth = textWidth / 2
	xPos = mp.DisplayWidth - int16(60) - halfTextWidth
	mp.drawText(xPos, startY+mp.fontHeight, text)
}

// drawGrid draws a grid on the plot.
// It should draw horizontal lines with space for each 25% of the Y-axis,
func (mp *MiniPlot) drawGrid() {
}

// drawData draws an auto-scrolling line chart that grows from the right edge.
func (mp *MiniPlot) drawData(samples []int16, minV, maxV int16) {
	var (
		baseY  = int16(22)           // the horizontal axis pixel row
		rightX = mp.DisplayWidth - 2 // right-most drawable column
		topY   = int16(1)            // graph top (screen row 1)
	)

	if len(samples) == 0 {
		return
	}

	rangeVal := float64(maxV - minV)
	if rangeVal == 0 {
		rangeVal = 1
	}
	pixelsPerUnit := float64(baseY-topY) / rangeVal

	pixelY := func(v int16) int16 {
		return int16(float64(maxV-v) * pixelsPerUnit)
	}

	n := len(samples)
	for i := 1; i < n; i++ {
		x1 := rightX - int16(i-1)
		x2 := rightX - int16(i)
		y1 := pixelY(samples[n-i])
		y2 := pixelY(samples[n-i-1])

		tinydraw.Line(mp.display, x1, y1, x2, y2, mp.Color)
	}

	mp.display.Display()
}

func (mp *MiniPlot) formatValue(value int16) string {
	if value >= 1000 {
		return fmt.Sprintf("%dk", value/1000)
	}

	return fmt.Sprintf("%d", value)
}
