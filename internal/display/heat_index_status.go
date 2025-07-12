package display

import (
	"fmt"

	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayReadingsWithHIAndStatus(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	radiusFilled := int16(3)

	font7 := font.NewFont7(f.display)
	font11 := font.NewFont11(f.display)
	if r.ValidityError != "" {
		font7.Print(0, 0, r.ValidityError)
	} else {
		// CO2Status
		status := "CO2"
		font7.Print(0, 0, status)

		f.DrawBar(30, 4, r.CO2Index(), radiusFilled)

		// CO2 value
		co2Value := fmt.Sprintf("%d", r.CO2)
		var (
			YPos int16 = 0
			XPos       = int16(128 - (len(co2Value) * 11))
		)
		font11.Print(XPos, YPos, co2Value)
	}

	// Heat Index status
	hiStatus := "HI"
	font7.Print(0, 10, hiStatus)

	f.DrawBar(30, 14, r.HeatIndexRating(), radiusFilled)

	tempHum := fmt.Sprintf("%.0f %.0f", r.Temperature, r.Humidity)
	var (
		YPos int16 = 16
		XPos int16 = int16(128 - (len(tempHum) * 11))
	)
	font11.Print(XPos, YPos, tempHum)

	status := r.ComfortStatus()
	font7.Print(0, 20, status)
}
