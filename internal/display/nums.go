package display

import (
	"fmt"

	"pico_co2/internal/types"
)

func RenderNums(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}
	renderer.Clear()

	var (
		y int16
		x int16
	)

	renderer.DrawSmallText(x, y, fmt.Sprintf("CO2: %s", r.Calculated.CO2Status))

	x = 0
	y = 8
	renderer.DrawXLargeText(x, y, fmt.Sprintf("%d", r.Raw.CO2))

	x = 90
	y = 0
	renderer.DrawSmallText(x, y, "T")

	temp := fmt.Sprintf("%.0f", r.Raw.Temperature)
	x = 128 - renderer.CalcLargeTextWidth(temp)
	y = 0
	renderer.DrawLargeText(x, y, temp)

	x = 90
	y = 16
	renderer.DrawSmallText(x, y, "H")

	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	x = 128 - renderer.CalcLargeTextWidth(hum)
	y = 16
	renderer.DrawLargeText(x, y, hum)

	renderer.Display()
}
