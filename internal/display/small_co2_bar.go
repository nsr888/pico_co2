package display

import (
	"fmt"

	font "github.com/Nondzu/ssd1306_font"

	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayWithSmallCO2Bar(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	// Status top right corner
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := r.CO2Status
	f.font.YPos = 0
	f.font.XPos = int16(128 - (len(status) * 7))
	if r.ValidityError != "" {
		status = r.ValidityError
		f.font.XPos = 0
	}
	f.font.PrintText(status)

	// Bars for CO2 level
	if r.ValidityError == "" {
		f.font.Configure(font.Config{FontType: font.FONT_7x10})
		f.font.YPos = 0
		f.font.XPos = 0
		f.font.PrintText(printBar(r.CO2))
	}

	// Bottom big numbers
	if r.ValidityError == "" {
		f.font.Configure(font.Config{FontType: font.FONT_11x18})
		f.font.YPos = 16
		f.font.XPos = 0
		f.font.PrintText(fmt.Sprintf("%d", r.CO2))
	}

	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	temp := fmt.Sprintf("%.0f", r.Temperature)
	f.font.YPos = 16
	f.font.XPos = 64
	f.font.PrintText(temp)

	hum := fmt.Sprintf("%.0f", r.Humidity)
	f.font.YPos = 16
	f.font.XPos = 128 - 22 - 6
	f.font.PrintText(hum)
}
