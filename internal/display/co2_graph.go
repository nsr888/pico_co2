package display

import (
	"pico_co2/internal/types"
)

func RenderCO2Graph(renderer Renderer, r *types.Readings) {
	if renderer == nil {
		return
	}
	renderer.Clear()

	rawData := r.History.CO2.Contiguous()

	renderer.DrawPlot(rawData, "CO2")

	renderer.Display()
}
