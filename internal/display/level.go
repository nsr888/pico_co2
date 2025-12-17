package display

import (
	"fmt"
	"pico_co2/internal/display/font"
	"pico_co2/internal/types"
)

func RenderLevel(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}

	renderer.Clear()

	width, _ := renderer.Size()

	var (
		space int16 = 8
		lf          = renderer.GetFont(font.FreemonoRegular12)
		sf          = renderer.GetFont(font.ProggySZ8)
	)

	lf.Print(0, 0, fmt.Sprintf("%s", r.Calculated.CO2Status))

	humStr := fmt.Sprintf("H %.0f", r.Raw.Humidity)
	humWidth := sf.CalcWidth(humStr)
	sf.Print(width-humWidth, 24, humStr)

	tempStr := fmt.Sprintf("T %.0f", r.Raw.Temperature)
	tempWidth := sf.CalcWidth(tempStr)
	sf.Print(width-tempWidth-space-humWidth, 24, tempStr)

	co2Str := fmt.Sprintf("CO2 %d", r.Raw.CO2)
	sf.Print(0, 24, co2Str)

	renderer.Display()
}
