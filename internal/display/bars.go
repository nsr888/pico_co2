package display

import (
	"fmt"

	"pico_co2/internal/display/bar"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func (f *FontDisplay) DisplayBars(r *types.Readings) {
	if f == nil {
		return
	}
	f.clearDisplay()

	var (
		lineY        int16 = 0
		radiusFilled int16 = 3
		barSpacing   int16 = 5
	)

	font7 := font.NewFont7(f.display)
	if font7 == nil {
		return
	}
	hbarTop := bar.NewTwoSideBar(
		f.display,
		radiusFilled,
		barSpacing,
		0,
		4,
		font.NewFont7(f.display),
	)

	if r.ValidityError != "" {
		font7.Print(0, lineY, r.ValidityError)
	} else {
		x := hbarTop.Draw(0, lineY, r.CO2Index(), "C")
		hbarTop.Draw(x+9, lineY, int16(r.AQI-1), "A")
	}

	hbarBottom := bar.NewTwoSideBar(
		f.display,
		radiusFilled,
		barSpacing,
		3,
		4,
		font.NewFont7(f.display),
	)

	lineY = 11
	hbarBottom.Draw(0, lineY, r.CalculateComfortIndex(), "TEM")
	tem := fmt.Sprintf("%.0f", r.Temperature)
	font7.Print(128-font7.CalcWidth(tem), lineY, tem)

	lineY = 22
	hbarBottom.Draw(0, lineY, r.HumidityComfortIndex(), "HUM")
	hum := fmt.Sprintf("%.0f", r.Humidity)
	font7.Print(128-font7.CalcWidth(hum), lineY, hum)
}
