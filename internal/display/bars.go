package display

import (
	"fmt"

	"pico_co2/internal/types"
	"pico_co2/internal/types/status"
)

func RenderBars(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	var lineY int16 = 0

	if r.Warning != "" {
		renderer.DrawSmallText(0, lineY, r.Warning)
	} else {
		renderer.DrawSmallText(0, lineY, fmt.Sprintf("AQI %d", r.Raw.AQI))
		renderer.DrawTwoSideBar(36, lineY, int16(r.Calculated.CO2Status-1), "CO2", 0, 4)
		co2 := fmt.Sprintf("%d", r.Raw.CO2)
		renderer.DrawSmallText(128-renderer.CalcSmallTextWidth(co2), lineY, co2)
	}

	lineY = 11
	renderer.DrawTwoSideBar(0, lineY, int16(status.CalculateComfortIndex(r.Raw.Temperature, r.Raw.Humidity)), "TEM", 3, 4)
	tem := fmt.Sprintf("%.0f", r.Raw.Temperature)
	renderer.DrawSmallText(128-renderer.CalcSmallTextWidth(tem), lineY, tem)

	lineY = 22
	renderer.DrawTwoSideBar(0, lineY, int16(status.HumidityComfortIndex(r.Raw.Humidity)), "HUM", 3, 4)
	hum := fmt.Sprintf("%.0f", r.Raw.Humidity)
	renderer.DrawSmallText(128-renderer.CalcSmallTextWidth(hum), lineY, hum)

	renderer.Display()
}
