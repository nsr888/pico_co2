package display

import (
	"fmt"
	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderCO2BarWithNums(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var (
		y int16
		x int16
	)

	x = 0
	y = 0
	hi := status.GetHeatIndex(r.Raw.Temperature, r.Raw.Humidity)
	renderer.DrawTwoSideBar(x, y, int16(hi), "HEAT  ", 0, 4)

	x = 0
	y = 11
	co2status := int16(status.CO2Index(r.Raw.CO2))
	renderer.DrawTwoSideBar(x, y, co2status, "CO2   ", 0, 4)

	x = 0
	y = 22
	co2Str := fmt.Sprintf("       %d", r.Raw.CO2)
	renderer.DrawSmallText(x, y, co2Str)

	temp := fmt.Sprintf("%.0f", r.Raw.Temperature)
	x = 128 - renderer.CalcLargeTextWidth(temp)
	y = 0
	renderer.DrawLargeText(x, y, temp)

	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	x = 128 - renderer.CalcLargeTextWidth(hum)
	y = 16
	renderer.DrawLargeText(x, y, hum)

	renderer.Display()
}
