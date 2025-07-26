package display

import (
	"fmt"

	"pico_co2/internal/types"
)

func RenderLargeBar(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}
	renderer.Clear()

	width, _ := renderer.Size()

	var (
		XPos int16 = 0
		YPos int16 = 0
	)
	renderer.DrawSmallText(XPos, YPos, r.Calculated.CO2Status.String())

	XPos = 0
	YPos = 12
	renderer.DrawSquareBar(XPos, YPos, uint8(r.Calculated.CO2Status))

	co2Str := fmt.Sprintf("CO2 %d", r.Raw.CO2)
	if r.Warning != "" {
		co2Str = r.Warning
	}
	XPos = 0
	YPos = 24
	renderer.DrawSmallText(XPos, YPos, co2Str)

	humStr := fmt.Sprintf("H %.0f", r.Raw.Humidity)
	humWidth := renderer.CalcSmallTextWidth(humStr)
	XPos = int16(width - humWidth)
	YPos = 24
	renderer.DrawSmallText(XPos, YPos, humStr)

	tempStr := fmt.Sprintf("T %.0f", r.Raw.Temperature)
	tempWidth := renderer.CalcSmallTextWidth(tempStr)
	XPos = int16(width - (humWidth) - (tempWidth) - 8) // 8 for padding
	YPos = 24
	renderer.DrawSmallText(XPos, YPos, tempStr)

	renderer.Display()
}
