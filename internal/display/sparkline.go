package display

import (
	"fmt"

	"pico_co2/internal/types"
)

func RenderSparkline(renderer Renderer, r *types.Readings) {
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

	renderer.DrawTwoSideBar(x, y, int16(r.Calculated.CO2Status), "CO2", 0, 4)

	humStr := fmt.Sprintf("H %.0f", r.Raw.Humidity)
	humWidth := renderer.CalcSmallTextWidth(humStr)
	width, _ := renderer.Size()
	renderer.DrawSmallText(width-humWidth, 0, humStr)

	tempStr := fmt.Sprintf("T %.0f", r.Raw.Temperature)
	tempWidth := renderer.CalcSmallTextWidth(tempStr)
	renderer.DrawSmallText(width-tempWidth-humWidth-8, 0, tempStr)

	x = 0
	y = 11

	var graphHeight int16 = 21
	var graphWidth int16 = 128

	data := r.History.CO2.CopyTo()
	renderer.DrawSparkline(x, y, data, graphWidth, graphHeight)

	renderer.Display()
}
