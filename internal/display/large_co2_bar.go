package display

import (
	"fmt"
	"strings"

	font "github.com/Nondzu/ssd1306_font"

	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayWithLargeCO2Bar(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	// Display CO2 status string at the top right corner
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	status := r.CO2Status
	f.font.XPos = int16(128 - (len(status) * 7))
	f.font.YPos = 0
	f.font.PrintText(status)

	// Bars for CO2 level
	f.font.Configure(font.Config{FontType: font.FONT_11x18})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(printBar(r.CO2))

	// Small font
	f.font.Configure(font.Config{FontType: font.FONT_7x10})
	co2Str := fmt.Sprintf("CO2 %d", r.CO2)
	if r.ValidityError != "" {
		co2Str = fmt.Sprintf("CO2 %d*", r.CO2)
	}
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(co2Str)

	humStr := fmt.Sprintf("H %.0f", r.Humidity)
	f.font.XPos = int16(128 - (len(humStr) * 7))
	f.font.YPos = 24
	f.font.PrintText(humStr)
	tempStr := fmt.Sprintf("T %.0f", r.Temperature)
	f.font.XPos = int16(128 - ((len(humStr) * 7) + (len(tempStr) * 7) + 8)) // 8 for padding
	f.font.YPos = 24
	f.font.PrintText(tempStr)
}

func printBar(count uint16) string {
	cnt := int(count)

	cnt = cnt - 400 // Adjust count to start from 400

	if cnt <= 0 {
		return ""
	}

	numBars := cnt / 100 // Corrected value 1 - 11

	return strings.Repeat("|", numBars)
}
