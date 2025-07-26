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

	width, _ := renderer.Size()

	renderer.DrawSmallText(0, 0, fmt.Sprintf("CO2: %s", r.Calculated.CO2Status))

	renderer.DrawXLargeText(0, 8, fmt.Sprintf("%d", r.Raw.CO2))

	tempStr := fmt.Sprintf("T %.0f", r.Raw.Temperature)
	XPos := int16(width - (renderer.CalcSmallTextWidth(tempStr)))
	YPos := int16(0)
	renderer.DrawSmallText(XPos, YPos, tempStr)

	humStr := fmt.Sprintf("H %.0f", r.Raw.Humidity)
	XPos = int16(width - (renderer.CalcSmallTextWidth(humStr)))
	YPos = int16(11)
	renderer.DrawSmallText(XPos, YPos, humStr)

	aqiStr := fmt.Sprintf("AQI %d", r.Raw.AQI)
	XPos = int16(width - (renderer.CalcSmallTextWidth(aqiStr)))
	YPos = int16(22)
	renderer.DrawSmallText(XPos, YPos, aqiStr)

	renderer.Display()
}
