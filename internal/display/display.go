package display

import (
	"log"

	"machine"

	font "github.com/Nondzu/ssd1306_font"
	"tinygo.org/x/drivers/ssd1306"

	"pico_co2/pkg/fifo"
	"pico_co2/pkg/sparkline"
)

// Display Configuration
const (
	displayWidth   int16 = 128
	displayHeight  int16 = 32
	displayAddress       = ssd1306.Address_128_32
	graphWidth     int16 = 60
	graphHeight    int16 = 14
	barWidth       int16 = 1
)

func NewFontDisplay(bus *machine.I2C) (*FontDisplay, error) {
	display := ssd1306.NewI2C(bus)
	display.Configure(ssd1306.Config{
		Width:   displayWidth,
		Height:  displayHeight,
		Address: displayAddress,
	})
	log.Printf("Display configured: Width=%d, Height=%d, Address=%#x\n", displayWidth, displayHeight, displayAddress)

	fontLib := font.NewDisplay(display)

	var graphMeasurementsCount int16 = graphWidth / barWidth

	return &FontDisplay{
		font:         &fontLib,
		clearDisplay: display.ClearDisplay,
		display:      &display,
		humFIFO:      fifo.NewFIFO16(int(graphMeasurementsCount)),
		sparkline:    sparkline.NewSparkline(int(graphHeight)),
	}, nil
}

type FontDisplay struct {
	display      *ssd1306.Device
	font         *font.Display
	clearDisplay func()
	humFIFO      *fifo.FIFO16
	sparkline    *sparkline.Sparkline
}
