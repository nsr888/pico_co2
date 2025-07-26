package display

import (
	"image/color"

	"pico_co2/pkg/fifo"
	"pico_co2/pkg/sparkline"
)

const (
	graphWidth  int16 = 128
	graphHeight int16 = 14
	barWidth    int16 = 1
)

type FontDisplay struct {
	renderer     Renderer
	clearDisplay func()
	humFIFO      *fifo.FIFO16
	sparkline    *sparkline.Sparkline
	color        color.RGBA
	width        int16
	height       int16
}

func NewFontDisplay(renderer Renderer) (*FontDisplay, error) {
	var (
		graphMeasurementsCount int16 = graphWidth / barWidth
		white                        = color.RGBA{255, 255, 255, 255}
		black                        = color.RGBA{0, 0, 0, 255}
	)
	displayWidth, displayHeight := renderer.Size()

	return &FontDisplay{
		clearDisplay: func() {
			for y := int16(0); y < displayHeight; y++ {
				for x := int16(0); x < displayWidth; x++ {
					renderer.SetPixel(x, y, black)
				}
			}
		},
		renderer:  renderer,
		humFIFO:   fifo.NewFIFO16(int(graphMeasurementsCount)),
		sparkline: sparkline.NewSparkline(int(graphHeight)),
		color:     white,
		width:     displayWidth,
		height:    displayHeight,
	}, nil
}

// Renderer returns the underlying renderer
func (f *FontDisplay) Renderer() Renderer {
	return f.renderer
}
