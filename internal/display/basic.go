package display

import (
	"fmt"

	"pico_co2/internal/types"
)

func RenderBasic(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	width, _ := renderer.Size()
	var space int16 = 8

	if r.Warning != "" {
		renderer.DrawSmallText(0, 0, r.Warning)
	} else {
		renderer.DrawXLargeText(0, 0, fmt.Sprintf("%s", r.Calculated.CO2Status))
	}

	humStr := fmt.Sprintf("H %.0f", r.Raw.Humidity)
	humWidth := renderer.CalcSmallTextWidth(humStr)
	renderer.DrawSmallText(width-humWidth, 24, humStr)

	tempStr := fmt.Sprintf("T %.0f", r.Raw.Temperature)
	tempWidth := renderer.CalcSmallTextWidth(tempStr)
	renderer.DrawSmallText(width-tempWidth-space-humWidth, 24, tempStr)

	co2Str := fmt.Sprintf("CO2 %d", r.Raw.CO2)
	renderer.DrawSmallText(0, 24, co2Str)

	renderer.Display()
}
