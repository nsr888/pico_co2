package display

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/proggy"

	"pico_co2/internal/display/bar"
	"pico_co2/pkg/miniplot"
	"pico_co2/pkg/sparkline"
)

type VirtualDisplay struct {
	buffer []color.RGBA
	width  int16
	height int16
	white  color.RGBA
	black  color.RGBA
}

func NewVirtualDisplay(w, h int16) *VirtualDisplay {
	size := w * h
	return &VirtualDisplay{
		buffer: make([]color.RGBA, size),
		width:  w,
		height: h,
		white:  color.RGBA{255, 255, 255, 255},
		black:  color.RGBA{0, 0, 0, 255},
	}
}

func (v *VirtualDisplay) Clear() {
	for i := range v.buffer {
		v.buffer[i] = v.black
	}
}

func (v *VirtualDisplay) SetPixel(x, y int16, c color.RGBA) {
	if x >= 0 && x < v.width && y >= 0 && y < v.height {
		v.buffer[y*v.width+x] = c
	}
}

func (v *VirtualDisplay) Size() (int16, int16) {
	return v.width, v.height
}

func (v *VirtualDisplay) Display() error {
	return nil
}

func (v *VirtualDisplay) SavePNG(filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, int(v.width), int(v.height)))
	for y := 0; y < int(v.height); y++ {
		for x := 0; x < int(v.width); x++ {
			idx := y*int(v.width) + x
			img.Set(x, y, v.buffer[idx])
		}
	}

	// Resize the image to 5x its original size
	resizedImg := resize.Resize(
		uint(v.width*2),
		uint(v.height*2),
		img,
		resize.NearestNeighbor,
	)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, resizedImg)
}

func (v *VirtualDisplay) DrawXLargeText(x, y int16, text string) {
	font := &freemono.Bold18pt7b
	y += int16(font.GetGlyph('0').Info().Height)
	tinyfont.WriteLine(v, font, x, y, text, v.white)
}

func (v *VirtualDisplay) CalcXLargeTextWidth(text string) int16 {
	font := &freemono.Bold18pt7b
	if len(text) == 0 {
		return 0
	}

	_, width := tinyfont.LineWidth(font, text)

	return int16(width)
}

func (v *VirtualDisplay) DrawLargeText(x, y int16, text string) {
	font := &freemono.Regular12pt7b
	y += int16(font.GetGlyph('0').Info().Height)
	tinyfont.WriteLine(v, font, x, y, text, v.white)
}

func (v *VirtualDisplay) CalcLargeTextWidth(text string) int16 {
	font := &freemono.Regular12pt7b
	if len(text) == 0 {
		return 0
	}

	_, width := tinyfont.LineWidth(font, text)

	return int16(width)
}

func (v *VirtualDisplay) DrawSmallText(x, y int16, text string) {
	font := &proggy.TinySZ8pt7b
	y += int16(font.GetGlyph('0').Info().Height)
	tinyfont.WriteLine(v, font, x, y, text, v.white)
}

func (v *VirtualDisplay) CalcSmallTextWidth(text string) int16 {
	font := &proggy.TinySZ8pt7b
	if len(text) == 0 {
		return 0
	}

	_, width := tinyfont.LineWidth(font, text)

	return int16(width)
}

// wrapText splits a string into lines at word boundaries to fit a maximum width.
func wrapText(text string, font *tinyfont.Font, maxWidth int16) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	currentLine := words[0]

	for i := 1; i < len(words); i++ {
		word := words[i]
		testLine := currentLine + " " + word
		_, w := tinyfont.LineWidth(font, testLine)

		if int16(w) > maxWidth {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}
	lines = append(lines, currentLine)
	return lines
}

func (v *VirtualDisplay) DrawLongText(x, y int16, text string) {
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

func (v *VirtualDisplay) DrawPlot(data []int16, title string) {
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

func (v *VirtualDisplay) DrawTwoSideBar(
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

func (v *VirtualDisplay) DrawSparkline(
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

func (v *VirtualDisplay) DrawSquareBar(x, y int16, value uint8) {
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
