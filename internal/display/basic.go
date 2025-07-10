package display

import (
	"fmt"
	"strings"

	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayBasicTempAndHumidity(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	lines := strings.Split(r.ValidityError, ": ")

	f.printLines(lines[:1])
	humStr := fmt.Sprintf("H %.0f", r.Humidity)
	f.font.XPos = int16(128 - (len(humStr) * 7))
	f.font.YPos = 24
	f.font.PrintText(humStr)
	tempStr := fmt.Sprintf("T %.0f", r.Temperature)
	f.font.XPos = int16(128 - ((len(humStr) * 7) + (len(tempStr) * 7) + 8)) // 8 for padding
	f.font.YPos = 24
	f.font.PrintText(tempStr)
}
