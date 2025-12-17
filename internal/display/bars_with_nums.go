package display

import (
	"fmt"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderBarsWithLargeNums(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var (
		y          int16 = 0
		x          int16 = 0
		barSpacing int16 = 5
		lf               = renderer.GetFont(font.FreemonoRegular12)
	)

	// Print first line with CO2 and AQI
	comfortIndex := status.CalculateComfortIndex(
		r.Raw.Temperature,
		r.Raw.Humidity,
	)
	x = renderer.DrawTwoSideBar(x, y, int16(comfortIndex), "T", 0, 4)
	co2status := int16(r.Calculated.CO2Status)
	renderer.DrawTwoSideBar(x+barSpacing, y, co2status, "C", 0, 4)

	x = 0
	y = 16
	// Print temperature
	temp := fmt.Sprintf("%.0f", r.Raw.Temperature)
	lf.Print(x, y, temp)

	// Print humidity
	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	x = int16(x + 10 + renderer.CalcLargeSansTextWidth(temp))
	lf.Print(x, y, hum)

	// Print CO2 value
	co2 := fmt.Sprintf("%d", r.Raw.CO2)
	x = 128 - renderer.CalcLargeTextWidth(co2)
	lf.Print(x, y, co2)

	renderer.Display()
}
