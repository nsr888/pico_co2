package display

import (
	"fmt"

	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderBarsWithLargeNums(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var (
		y          int16
		x          int16
		barSpacing int16 = 5
	)

	x = 0
	y = 0
	// Print first line with CO2 and AQI
	x = renderer.DrawTwoSideBar(x, y, int16(status.CalculateComfortIndex(r.Raw.Temperature, r.Raw.Humidity)), "T", 0, 4)
	renderer.DrawTwoSideBar(x+barSpacing, y, int16(r.Calculated.CO2Status), "C", 0, 4)

	x = 0
	y = 16
	// Print temperature
	temp := fmt.Sprintf("%.0f", r.Raw.Temperature)
	renderer.DrawLargeSansText(x, y, temp)

	// Print humidity
	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	x = int16(x + 10 + renderer.CalcLargeSansTextWidth(temp))
	renderer.DrawLargeSansText(x, y, hum)

	// Print CO2 value
	co2 := fmt.Sprintf("%d", r.Raw.CO2)
	x = 128 - renderer.CalcLargeTextWidth(co2)
	renderer.DrawLargeSansText(x, y, co2)

	renderer.Display()
}
