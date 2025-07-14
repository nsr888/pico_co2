package display

import (
	"fmt"
	"strings"

	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayWithLargeCO2Bar(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()
	font7 := font.NewFont7(f.display)
	font11 := font.NewFont11(f.display)

	// Display CO2 status string at the top right corner
	status := r.CO2Status
	var (
		XPos       = int16(128 - (len(status) * 7))
		YPos int16 = 0
	)
	font7.Print(XPos, YPos, status)

	// Bars for CO2 level
	XPos = 0
	YPos = 0
	font11.Print(XPos, YPos, printBar(r.CO2))

	// Small font
	co2Str := fmt.Sprintf("CO2 %d", r.CO2)
	if r.ValidityError != "" {
		co2Str = fmt.Sprintf("CO2 %d*", r.CO2)
	}
	XPos = 0
	YPos = 24
	font7.Print(XPos, YPos, co2Str)

	humStr := fmt.Sprintf("H %.0f", r.Humidity)
	XPos = int16(128 - (len(humStr) * 7))
	YPos = 24
	font7.Print(XPos, YPos, humStr)
	tempStr := fmt.Sprintf("T %.0f", r.Temperature)
	XPos = int16(128 - ((len(humStr) * 7) + (len(tempStr) * 7) + 8)) // 8 for padding
	YPos = 24
	font7.Print(XPos, YPos, tempStr)
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
