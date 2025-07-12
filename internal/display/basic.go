package display

import (
	"fmt"
	"strings"

	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayBasic(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	errMsg := strings.Split(r.ValidityError, ": ")

	font7 := font.NewFont7(f.display)

	f.printLines(errMsg[:1], font7)

	humStr := fmt.Sprintf("H %.0f", r.Humidity)
	font7.Print(int16(128-font7.CalcWidth(humStr)), 24, humStr)

	tempStr := fmt.Sprintf("T %.0f", r.Temperature)
	font7.Print(int16(128-font7.CalcWidth(humStr)-font7.CalcWidth(tempStr)-8), 24, tempStr)
}
