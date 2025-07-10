package display

import (
	"fmt"

	font "github.com/Nondzu/ssd1306_font"

	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayComfortIndex(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := r.ComfortStatus()
	f.font.YPos = 0
	f.font.XPos = 0
	if r.ValidityError != "" {
		status = r.ValidityError
		f.font.XPos = 0
	}
	f.font.PrintText(status)

	if r.ValidityError == "" {
		f.font.Configure(font.Config{FontType: font.FONT_11x18})
		f.font.YPos = 16
		f.font.XPos = 0
		f.font.PrintText(fmt.Sprintf("%d", r.CO2))
	}

	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	hum := fmt.Sprintf("%.0f", r.Humidity)
	f.font.YPos = 16
	f.font.XPos = 128 - 22
	f.font.PrintText(hum)

	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	temp := fmt.Sprintf("%.0f", r.Temperature)
	f.font.YPos = 16
	f.font.XPos = int16(128 - len(temp)*11 - len(hum)*11 - 11)
	f.font.PrintText(temp)

	f.humFIFO.Enqueue(int16(r.CO2))
}
