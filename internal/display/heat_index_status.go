package display

import (
	"fmt"

	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderHeatIndexStatus(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}
	renderer.Clear()
	var (
		y int16
		x int16
	)

	width, _ := renderer.Size()

	x = 0
	y = 0
	renderer.DrawTwoSideBar(x, y, int16(status.CO2Index(r.Raw.CO2)), "CO2 ", 0, 4)

	co2Value := fmt.Sprintf("%d", r.Raw.CO2)
	renderer.DrawLargeText(int16(width-renderer.CalcLargeTextWidth(co2Value)), y, co2Value)

	// Heat Index status
	x = 0
	y = 11
	hi := status.GetHeatIndex(r.Raw.Temperature, r.Raw.Humidity)
	renderer.DrawTwoSideBar(x, y, int16(hi), "HI  ", 0, 4)

	y = 22
	humStr := fmt.Sprintf("%.0f", r.Raw.Humidity)
	humWidth := renderer.CalcSmallTextWidth(humStr)
	renderer.DrawSmallText(int16(width-humWidth), y, humStr)
	tempStr := fmt.Sprintf("%.0f", r.Raw.Temperature)
	tempWidth := renderer.CalcSmallTextWidth(tempStr)
	renderer.DrawSmallText(int16(width-humWidth-tempWidth-5), y, tempStr)

	status := status.ComfortStatus(
		r.Raw.CO2,
		r.Raw.AQI,
		r.Raw.Humidity,
		r.Raw.Temperature,
	)
	renderer.DrawSmallText(x, y, status)

	renderer.Display()
}
