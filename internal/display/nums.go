package display

import (
	"fmt"

	font "github.com/Nondzu/ssd1306_font"

	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayWithLargeCO2AndTempNums(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	// Big numbers for eCO2 and temperature
	f.font.Configure(font.Config{FontType: font.FONT_16x26})
	f.font.XPos = 0
	f.font.YPos = 0
	f.font.PrintText(fmt.Sprintf("%d", r.CO2))
	tempStr := fmt.Sprintf("%.0f", r.Temperature)
	f.font.XPos = int16(128 - (len(tempStr) * 16))
	f.font.YPos = 0
	f.font.PrintText(tempStr)

	f.font.Configure(font.Config{FontType: font.FONT_6x8})
	co2Str := "eCO2"
	f.font.XPos = 0
	f.font.YPos = 24
	f.font.PrintText(co2Str)
	tempTitleStr := "Temp"
	f.font.XPos = int16(128 - (len(tempTitleStr) * 6))
	f.font.YPos = 24
	f.font.PrintText(tempTitleStr)
	f.font.XPos = int16(128-(len(r.CO2Status)*6)) / 2
	f.font.YPos = 24
	f.font.PrintText(r.CO2Status)
}
