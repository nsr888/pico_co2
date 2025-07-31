package display

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/proggy"

	"pico_co2/internal/display/bar"
	"pico_co2/pkg/miniplot"
	"pico_co2/pkg/sparkline"
)

// SSD1306Adapter wraps an ssd1306 display
type SSD1306Adapter struct {
	dev   drivers.Displayer
	white color.RGBA
	black color.RGBA
}

func NewSSD1306Adapter(dev drivers.Displayer) *SSD1306Adapter {
	return &SSD1306Adapter{
		dev:   dev,
		white: color.RGBA{255, 255, 255, 255},
		black: color.RGBA{0, 0, 0, 255},
	}
}

func (v *SSD1306Adapter) SetPixel(x, y int16, c color.RGBA) {
	v.dev.SetPixel(x, y, c)
}

func (v *SSD1306Adapter) Size() (int16, int16) {
	return v.dev.Size()
}

func (v *SSD1306Adapter) Clear() {
	width, height := v.dev.Size()
	for x := int16(0); x < width; x++ {
		for y := int16(0); y < height; y++ {
			v.dev.SetPixel(x, y, v.black)
		}
	}
}

func (v *SSD1306Adapter) Display() error {
	v.dev.Display()
	return nil
}

func (v *SSD1306Adapter) DrawXLargeText(x, y int16, text string) {
	font := &freemono.Bold18pt7b
	y += int16(font.GetGlyph('0').Info().Height)
	tinyfont.WriteLine(v, font, x, y, text, v.white)
}

func (v *SSD1306Adapter) CalcXLargeTextWidth(text string) int16 {
	font := &freemono.Bold18pt7b
	if len(text) == 0 {
		return 0
	}

	_, width := tinyfont.LineWidth(font, text)

	return int16(width)
}

func (v *SSD1306Adapter) DrawLargeText(x, y int16, text string) {
	font := &freemono.Regular12pt7b
	y += int16(font.GetGlyph('0').Info().Height)
	tinyfont.WriteLine(v, font, x, y, text, v.white)
}

func (v *SSD1306Adapter) CalcLargeTextWidth(text string) int16 {
	font := &freemono.Regular12pt7b
	if len(text) == 0 {
		return 0
	}

	_, width := tinyfont.LineWidth(font, text)

	return int16(width)
}

func (v *SSD1306Adapter) DrawSmallText(x, y int16, text string) {
	font := &proggy.TinySZ8pt7b
	y += int16(font.GetGlyph('0').Info().Height)
	tinyfont.WriteLine(v, font, x, y, text, v.white)
}

func (v *SSD1306Adapter) CalcSmallTextWidth(text string) int16 {
	font := &proggy.TinySZ8pt7b
	if len(text) == 0 {
		return 0
	}

	_, width := tinyfont.LineWidth(font, text)

	return int16(width)
}

func (v *SSD1306Adapter) DrawLongText(x, y int16, text string) {
	font := &proggy.TinySZ8pt7b
	width, height := v.Size()

	// Wrap text to fit the display width, accounting for the initial x offset.
	lines := wrapText(text, font, width-x)

	fontHeight := int16(font.GetGlyph('A').Info().Height)
	lineHeight := fontHeight + 1 // Font height + 1px spacing

	// Calculate how many lines can be drawn in the available vertical space.
	availableHeight := height - y
	if availableHeight < 0 {
		availableHeight = 0
	}
	maxLines := int(availableHeight / lineHeight)

	for i, line := range lines {
		if i >= maxLines {
			break
		}
		// Y-coordinate for the baseline of the current line.
		drawY := y + fontHeight + (int16(i) * lineHeight)
		tinyfont.WriteLine(v, font, x, drawY, line, v.white)
	}
}

func (v *SSD1306Adapter) DrawPlot(data []int16, title string) {
	font := &proggy.TinySZ8pt7b
	width, height := v.Size()
	miniPlot, err := miniplot.NewMiniPlot(
		v,
		font,
		width,
		height,
		v.white,
	)
	if err != nil {
		println("Error creating CO2 graph:", err.Error())
		return
	}

	miniPlot.DrawLineChart(data, "CO2")
}

func (v *SSD1306Adapter) DrawTwoSideBar(
	x, y int16,
	value int16,
	label string,
	leftCount int16,
	rightCount int16,
) int16 {
	var (
		radius     int16 = 3
		barSpacing int16 = 5
	)

	font := &proggy.TinySZ8pt7b
	hbar := bar.NewTwoSideBar(
		v,
		radius,
		barSpacing,
		leftCount,
		rightCount,
		font,
		v.white,
	)

	return hbar.Draw(x, y, value, label)
}

func (v *SSD1306Adapter) DrawSparkline(
	x, y int16,
	data []int16,
	width int16,
	height int16,
) {
	if len(data) == 0 {
		return
	}

	sp := sparkline.NewSparkline(
		int(height),
	)

	if len(data) > int(width) {
		data = data[len(data)-int(width):]
	}

	result := sp.Process(data)

	y = y + height

	for i, barHeight := range result {
		topY := y - barHeight

		tinydraw.FilledRectangle(v, x+int16(i), topY, barWidth, barHeight, v.white)
	}
}

func (v *SSD1306Adapter) DrawSquareBar(x, y int16, value uint8) {
	if value == 0 {
		return
	}

	squareWidth := int16((128 - 6*3) / 4)
	squareHeight := int16(9)
	totalCount := 4
	filledCount := int16(value)

	for range totalCount {
		if filledCount > 0 {
			filledCount--
			tinydraw.FilledRectangle(v, x, y, squareWidth, squareHeight, v.white)
		} else {
			tinydraw.Rectangle(v, x, y, squareWidth, squareHeight, v.white)
		}
		x = x + squareWidth + 6
	}
}
