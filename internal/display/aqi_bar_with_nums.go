package display

import (
	"fmt"

	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderAqiBarWithNums(renderer Renderer, r *types.Readings) {
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
	renderer.DrawTwoSideBar(x, y, int16(r.Calculated.HeatIndex), "HEAT  ", 0, 4)

	x = 0
	y = 11
	humidityComfortIndex := status.HumidityComfortIndex(r.Raw.Humidity)
	renderer.DrawTwoSideBar(x, y, humidityComfortIndex, "HUMID ", 0, 4)

	x = 0
	y = 22
	if r.Warning != "" {
		renderer.DrawSmallText(x, y, r.Warning)
	} else {
		renderer.DrawTwoSideBar(x, y, int16(r.Raw.AQI-1), "AQI   ", 0, 4)
	}

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
